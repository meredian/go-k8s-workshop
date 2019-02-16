package handlers

import "net/http"

func (s *Server) ReadyzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
