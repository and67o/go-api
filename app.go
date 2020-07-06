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

var apiErrors []string
var answer Answer

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
	apiErrors = []string{}

	redisOrders, err := a.Redis.GetOrders()
	handleError(err)

	if len(redisOrders) == 0 {
		orders, err := a.DB.GetOrders()
		handleError(err)

		_, err = a.Redis.Set(_redis.USERS, orders)
		handleError(err)

		answer.setAnswerResult(orders)
	} else {
		answer.setAnswerResult(redisOrders)
	}

	answer.setAnswerErrors()

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
	apiErrors = []string{}

	redisUsers, err := a.Redis.GetUsers()
	handleError(err)

	if len(redisUsers) == 0 {
		users, err := a.DB.GetUsers()
		handleError(err)

		_, err = a.Redis.Set(_redis.USERS, users)
		handleError(err)

		answer.setAnswerResult(users)
	} else {
		answer.setAnswerResult(redisUsers)
	}

	answer.setAnswerErrors()

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
	apiErrors = []string{}
	answer=Answer{}

	UserID, err := strconv.Atoi(vars["userId"])
	handleError(err)

	redisUsers, err := a.Redis.GetUsers()
	handleError(err)

	if len(redisUsers) == 0 {
		user, err := a.DB.GetUser(UserID)
		handleError(err)
		answer.setAnswerResult(user)
	} else {
		findUserByTgId(UserID, redisUsers)
	}

	answer.setAnswerErrors()

	respondWithJSON(w, http.StatusOK, answer)
}

func (a *Answer) setAnswerErrors()  {
	a.Errors = apiErrors
}

func (a *Answer) setAnswerResult(res interface{})  {
	a.Result = res
}

func findUserByTgId(tgId int, redisUsers []db.User) {
	for _, user := range redisUsers {
		if tgId == user.TgId {
			answer.setAnswerResult(user)
			return
		}
	}
	apiErrors=append(apiErrors, "no user")
	return
}

func respondWithJSON(w http.ResponseWriter, code int, res Answer) {
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
