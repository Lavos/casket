package examples

import (
	"bytes"
	"io"
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"

	"github.com/Lavos/casket"
)

type Server struct {
	ContentStorer casket.ContentStorer
	Filer casket.Filer

	mux *web.Mux
	socket string
}

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func JSONCopy (data interface{}, w io.Writer) {
	var b bytes.Buffer
	j := json.NewEncoder(&b)
	j.Encode(data)
	b.WriteTo(w)
}

func NewServer(c casket.ContentStorer, f casket.Filer, socket string) *Server {
	return &Server{c, f, web.New(), socket}
}

func (s *Server) Run() {
	s.mux.Use(middleware.EnvInit)
	s.mux.Options("/*", s.Options)
	s.mux.Get("/r/:sha", s.GetRaw)
	s.mux.Get("/f/*", s.GetFileMeta)

	listener := bind.Socket(s.socket)
	log.Println("Starting Goji on", listener.Addr())

	graceful.HandleSignals()
	bind.Ready()

	err := graceful.Serve(listener, s.mux)

	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}

func (s *Server) Error(status_code int, error_code, error_message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status_code)
	JSONCopy(ErrorResponse{error_code, error_message}, w)
}

func (s *Server) Options(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "accept, content-type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH")
}

func (s *Server) GetFileMeta(c web.C, w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/f/")
	file, err := s.Filer.Get(name)

	log.Printf("file: %#v", file)

	if err != nil {
		s.Error(404, "file_notfound", "File was not found.", w)
		return
	}

	JSONCopy(file, w)
}

func (s *Server) GetRaw(c web.C, w http.ResponseWriter, r *http.Request) {
	sha_string := c.URLParams["sha"]
	log.Printf("sha_string: %s", sha_string)

	sha1sum := casket.NewSHA1SumFromString(sha_string)
	b, err := s.ContentStorer.Get(sha1sum)

	if err != nil {
		s.Error(404, "rev_notfound", "Revision was not found.", w)
		return
	}

	w.Write(b)
}
