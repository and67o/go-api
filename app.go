package main

import (
	"encoding/json"
	"errors"
	"go-api/db"
	_redis "go-api/redis"
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
	Redis  *_redis.Redis
}

var apiErrors = []string{}

func (a *App) Initialize() {
	a.DB = db.Db
	a.Router = mux.NewRouter()
	a.Redis, _ = _redis.NewClient()
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
	answer := Answer{}
	apiErrors = []string{}

	redisOrders, err := a.Redis.GetOrders()
	handleError(err)

	if len(redisOrders) == 0 {
		orders, err := a.DB.GetOrders()
		handleError(err)

		_, err = a.Redis.Set(_redis.USERS, orders)
		handleError(err)

		answer = Answer{orders, apiErrors}
	} else {
		answer = Answer{redisOrders, apiErrors}
	}

	respondWithJSON(w, http.StatusOK, answer)
}

func (a *App) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiErrors = []string{}
	orderID, _ := strconv.Atoi(vars["orderId"])
	order := a.DB.GetOrder(orderID)
	respondWithJSON(w, http.StatusOK, Answer{order, apiErrors})
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	answer := Answer{}
	apiErrors = []string{}

	redisUsers, err := a.Redis.GetUsers()
	handleError(err)

	if len(redisUsers) == 0 {
		users, err := a.DB.GetUsers()
		handleError(err)

		_, err = a.Redis.Set(_redis.USERS, users)
		handleError(err)

		answer = Answer{users, apiErrors}
	} else {
		answer = Answer{redisUsers, apiErrors}
	}
	respondWithJSON(w, http.StatusOK, answer)
}

func handleError(err error) {
	if err != nil {
		apiErrors = append(apiErrors, err.Error())
	}
}

func (a *App) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, ok := vars["orderId"]
	if ok {
		orderID, _ := strconv.Atoi(orderID)
		order, err := a.DB.DeleteOrder(orderID)
		if err != nil {
			apiErrors = append(apiErrors, err.Error())
		}
		respondWithJSON(w, http.StatusOK, Answer{order, apiErrors})
	}
	apiErrors = append(apiErrors, "no Id")
	respondWithJSON(w, http.StatusOK, Answer{Errors: apiErrors})
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	answer := Answer{}
	apiErrors = []string{}

	UserID, err := strconv.Atoi(vars["userId"])
	handleError(err)

	redisUsers, err := a.Redis.GetUsers()
	handleError(err)

	if len(redisUsers) == 0 {
		user, err := a.DB.GetUser(UserID)
		handleError(err)
		answer = Answer{user, apiErrors}
	} else {
		user, err := findUserByTgId(UserID, redisUsers)
		handleError(err)
		answer = Answer{user, apiErrors}
	}
	respondWithJSON(w, http.StatusOK, answer)
}

func findUserByTgId(tgId int, redisUsers []db.User) (user db.User, err error) {
	for _, user := range redisUsers {
		if tgId == user.TgId {
			return user, nil
		}
	}
	return db.User{}, errors.New("no user")
}

func respondWithJSON(w http.ResponseWriter, code int, res Answer) {
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
