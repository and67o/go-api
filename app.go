package main

import (
	"database/sql"
	"encoding/json"
	"go-api/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() {
	db := db.Db
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/orders", a.getOrders).Methods("GET")
}

func (a *App) getOrders(w http.ResponseWriter, r *http.Request) {
	products := `{"name":"John", "age":30, "city":"New York"}`
	respondWithJSON(w, http.StatusOK, products)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
