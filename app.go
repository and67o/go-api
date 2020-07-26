package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	DB     db.Operations
	Redis  _redis.Operations
}

var apiErrors []string
var answer Answer

func (a *App) Initialize() {
	a.DB = db.Db
	a.Router = mux.NewRouter()
	a.Redis = _redis.Redis
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/orders", a.getOrders).Methods("GET")
	a.Router.HandleFunc("/order/{orderId:[0-9]+}", a.getOrder).Methods("GET")
	a.Router.HandleFunc("/order/{orderId:[0-9]+}", a.deleteOrder).Methods("DELETE")
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.addUser).Methods("POST")
	a.Router.HandleFunc("/user/{userId:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{userId:[0-9]+}", a.deleteUser).Methods("DELETE")
}

func (a *App) getOrders(w http.ResponseWriter, r *http.Request) {
	apiErrors = []string{}
	answer = Answer{}

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

func sendError(err error, w http.ResponseWriter, status int) {
	handleError(err)
	answer.setAnswerErrors()
	respondWithJSON(w, status, answer)
	return
}

func (a *App) addUser(w http.ResponseWriter, r *http.Request) {
	user := db.User{}
	apiErrors = []string{}
	answer = Answer{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		sendError(errors.New("Invalid request payload"), w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newUser, err := a.DB.AddUser(user.Name, user.TgId)
	if err != nil {
		sendError(err, w, http.StatusNotFound)
		return
	}

	redisUsers, _ := a.Redis.GetUsers()
	if len(redisUsers) == 0 {
		users := redisUsers
		a.Redis.DeleteUsers()
		users = append(users, newUser)
		_, err = a.Redis.Set(_redis.USERS, users)
	}

	answer.setAnswerResult(newUser)

	respondWithJSON(w, http.StatusCreated, answer)
	return
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiErrors = []string{}
	answer = Answer{}

	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		sendError(errors.New("Invalid User ID"), w, http.StatusBadRequest)
		return
	}

	err = a.DB.DeleteUser(userId)

	if err != nil {
		sendError(err, w, http.StatusInternalServerError)
		return
	}

	answer.setAnswerResult(fmt.Sprintf("user %d delete ", userId))

	respondWithJSON(w, http.StatusOK, answer)
	return
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiErrors = []string{}
	answer = Answer{}

	UserID, err := strconv.Atoi(vars["userId"])

	if err != nil {
		sendError(err, w, http.StatusBadRequest)
		return
	}

	redisUsers, _ := a.Redis.GetUsers()
	if len(redisUsers) == 0 {
		user, err := a.DB.GetUser(UserID)
		if err != nil {
			sendError(err, w, http.StatusNotFound)
			return
		}
		answer.setAnswerResult(user)
	} else {
		user, err := findUserByTgId(int64(UserID), redisUsers)
		if err != nil {
			sendError(err, w, http.StatusNotFound)
			return
		}
		answer.setAnswerResult(user)
	}

	respondWithJSON(w, http.StatusOK, answer)
	return
}

func (a *Answer) setAnswerErrors() {
	a.Errors = apiErrors
}

func (a *Answer) setAnswerResult(res interface{}) {
	a.Result = res
}

func findUserByTgId(UserID int64, redisUsers []db.User) (user db.User, err error) {
	for _, user := range redisUsers {
		if UserID == user.Id {
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
