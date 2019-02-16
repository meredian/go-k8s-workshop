package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
)

type Server struct {
	Session *gocql.Session
}

type Action struct {
	MessageID string    `json:"message_id,-,omitempty"`
	UserID    string    `json:"user_id,-,omitempty"`
	Status    string    `json:"status,-,omitempty"`
	Timestamp time.Time `json:"timestamp,-,omitempty"`
}

func httpErr(w http.ResponseWriter, code int, err error) {
	log.Printf("Error while processing request: %v\n", err)
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`Error: %v`, err)))
}

func (s *Server) SaveActionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	defer r.Body.Close()

	action := &Action{}
	err = json.Unmarshal(body, action)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = s.Session.Query(
		`INSERT INTO tracking (messageID, userID, status, timestamp) VALUES (?, ?, ?, ?);`,
		action.MessageID, action.UserID, action.Status, time.Now(),
	).Exec()
	if err != nil {
		httpErr(w, 500, err)
		return
	}

	w.WriteHeader(201)
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(`{"result":"ok"}`))
}

func (s *Server) GetActionStatusHandler(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	messageID := queryString.Get("message_id")
	userID := queryString.Get("user_id")

	start := time.Now()

	fmt.Printf("message_id: %v, user_id: %v\n", messageID, userID)

	var status string
	err := s.Session.Query(
		`SELECT status FROM tracking where messageID = ? and userID = ?`,
		messageID, userID,
	).Scan(&status)

	if err == gocql.ErrNotFound {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		httpErr(w, 500, err)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`{"result":"%s"}`, status)))

	Counter.WithLabelValues(messageID).Inc()
	Timings.WithLabelValues(messageID).Observe(time.Since(start).Seconds())
}
