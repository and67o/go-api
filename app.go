package main

import (
	"encoding/json"
	"go-api/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Answer struct {
	Result interface{}
	Errors []string
}

type App struct {
	Router *mux.Router
	DB     db.DBOperations
}

var errors = []string{}

func (a *App) Initialize() {
	a.DB = db.Db
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/orders", a.getOrders).Methods("GET")
	a.Router.HandleFunc("/order/{orderId}", a.getOrder).Methods("GET")
	a.Router.HandleFunc("/order/{orderId}", a.getOrder).Methods("DELETE")
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user/{userId}", a.getUser).Methods("GET")
}

func (a *App) getOrders(w http.ResponseWriter, r *http.Request) {
	orders := a.DB.GetOrders()
	respondWithJSON(w, http.StatusOK, Answer{orders, errors})
}

func (a *App) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, _ := strconv.Atoi(vars["orderId"])
	order := a.DB.GetOrder(orderID)
	respondWithJSON(w, http.StatusOK, Answer{order, errors})
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	redisKey := "users"
	user, err := redis.client.Get()
	users := a.DB.GetUsers()
	respondWithJSON(w, http.StatusOK, Answer{users, errors})
}

func (a *App) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, ok := vars["orderId"]
	if ok {
		orderID, _ := strconv.Atoi(orderID)
		order, err := a.DB.DeleteOrder(orderID)
		if err == nil {
			errors = append(errors, err.Error())
		}
		respondWithJSON(w, http.StatusOK, Answer{order, errors})
	}
	errors = append(errors, "no Id")
	respondWithJSON(w, http.StatusOK, Answer{Errors: errors})
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserID, _ := strconv.Atoi(vars["userId"])
	user := a.DB.GetOrder(UserID)
	respondWithJSON(w, http.StatusOK, Answer{user, errors})
}

func respondWithJSON(w http.ResponseWriter, code int, res Answer) {
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
