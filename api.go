package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/account/{id}", jwtAuth(makeHttpHnadlerFunc(s.handleGetAccount)))
	router.HandleFunc("/transfer", makeHttpHnadlerFunc(s.handleTransfer))

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
	defer r.Body.Close()
	account := NewAccount(req.FirstName, req.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := signJWT(account)
	if err != nil {
		return err
	}
	fmt.Println(tokenString)
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
	return writeJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferreq := new(TransferReq)
	if err := json.NewDecoder(r.Body).Decode(transferreq); err != nil {
		return err
	}
	defer r.Body.Close()
	return writeJson(w, http.StatusOK, transferreq)
}

func signJWT(account *Account) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	claims := &jwt.MapClaims{
		"accountNumber": account.AccNumber,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func jwtAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("token")

		_, err := validateJWT(tokenString)

		if err != nil {
			writeJson(w, http.StatusForbidden, ApiErr{Error: "Invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
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
