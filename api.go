package main

import "net/http"

type ApiServer struct {
	listenAddr string
}


type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHttpHnadlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err:=f(w, r); err != nil {
			// handle the error
		}
	}
}

func NewApiServer(listenAddr string) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
	}
}

func (s *ApiServer) Run() {
	
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
