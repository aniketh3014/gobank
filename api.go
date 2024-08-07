package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	store      Storage
	listenAddr string
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHnadlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHnadlerFunc(s.handleGetAccount))
	log.Println("Api server is running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not supported")
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account := NewAccount(req.FirstName, req.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, account)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getId(r)

		if err != nil {
			return err
		}

		acc, err := s.store.GetAccoountById(id)
		if err != nil {
			return err
		}

		return writeJson(w, http.StatusOK, acc)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method: %s not supported", r.Method)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// json writer
func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiErr struct {
	Error string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHttpHnadlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJson(w, http.StatusBadRequest, ApiErr{Error: err.Error()})
		}
	}
}

func getId(r *http.Request) (int, error) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)

	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idstr)
	}
	return id, nil
}
