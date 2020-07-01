package main

import (
	"encoding/json"
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

var errors = []string{}

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
	redisUsers, err := a.Redis.Get(_redis.USERS, []db.User{})

	answer := Answer{}
	handleError(err)

	if redisUsers == "" {
		users, err := a.DB.GetUsers()
		handleError(err)

		_, err = a.Redis.Set(_redis.USERS, users)
		handleError(err)
		answer = Answer{users, errors}
	} else {
		answer = Answer{redisUsers, errors}
	}
	respondWithJSON(w, http.StatusOK, answer)
}

func handleError(err error) {
	if err != nil {
		errors = append(errors, err.Error())
	}
}

func (a *App) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, ok := vars["orderId"]
	if ok {
		orderID, _ := strconv.Atoi(orderID)
		order, err := a.DB.DeleteOrder(orderID)
		if err != nil {
			errors = append(errors, err.Error())
		}
		respondWithJSON(w, http.StatusOK, Answer{order, errors})
	}
	errors = append(errors, "no Id")
	respondWithJSON(w, http.StatusOK, Answer{Errors: errors})
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	redisKey := "users"
	UserID, _ := strconv.Atoi(vars["userId"])
	redisUsers, err := a.Redis.Client.Get(a.Redis.Client.Context(), redisKey).Result()
	if err != nil {
		panic(err)
	}
	if redisUsers == "" {
		user := a.DB.GetUser(UserID)
		respondWithJSON(w, http.StatusOK, Answer{user, errors})
	} else {
		var users []_redis.User
		err := json.Unmarshal([]byte(redisUsers), &users)
		if err != nil {
			panic(err)
		}
		user := findUserByTgId(UserID, users)
		if user.TgId == 0 {
			errors = append(errors, "No user")
			respondWithJSON(w, http.StatusOK, Answer{Errors: errors})
		} else {
			respondWithJSON(w, http.StatusOK, Answer{user, errors})
		}
	}
}

func findUserByTgId(tgId int, redisUsers []_redis.User) _redis.User {
	for _, user := range redisUsers {
		if tgId == user.TgId {
			return user
		}
	}
	return _redis.User{}
}

func respondWithJSON(w http.ResponseWriter, code int, res Answer) {
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
