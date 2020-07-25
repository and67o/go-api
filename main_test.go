package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

type AnswerTest struct {
	Result interface{}
	Errors []string
}

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize()

	checkForExistTable()
	code := m.Run()

	clearUsers()
	os.Exit(code)
}

func checkForExistTable() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearUsers() {
	if _, err := a.DB.DeleteUsers(); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1"); err != nil {
		log.Fatal(err)
	}
	a.Redis.DeleteUsers()
}

func TestEmptyTable(t *testing.T) {
	clearUsers()

	var answer Answer

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	err := json.Unmarshal(response.Body.Bytes(), &answer)
	if err != nil {
		t.Errorf(err.Error())
	}

	if answer.Result != nil {
		t.Errorf("Expected an empty array. Got %s", answer.Result)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
	id int NOT NULL AUTO_INCREMENT,
	tg_id int DEFAULT NULL,
	name varchar(36) NOT NULL,
	password varchar(200) NOT NULL,
	updated_at datetime DEFAULT NULL,
	created_at datetime NOT NULL,
	PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4`

func TestGetNonExistentUser(t *testing.T) {
	clearUsers()

	req, _ := http.NewRequest("GET", "/user/999", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m Answer
	json.Unmarshal(response.Body.Bytes(), &m)
	if m.Errors[0] != "no user" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m.Result)
	}
}

func TestCreateUser(t *testing.T) {
	clearUsers()
	var answer Answer

	payload := []byte(`{"name":"test user","tg_id":30}`)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	json.Unmarshal(response.Body.Bytes(), &answer)

	user := answer.Result

	data := user.(map[string]interface{})

	if data["name"] != "test user" {
		t.Errorf("Expected user name to be 'test user'. Got '%v'", data["name"])
	}

	if data["tg_id"] != 30.0 {
		t.Errorf("Expected user tgId to be '30'. Got '%v'", data["tg_id"])
	}

	if data["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", data["id"])
	}
}

func TestGetUser(t *testing.T) {
	a.Redis.DeleteUsers()
	a.DB.AddUser("oleg", 123.0)

	req, _ := http.NewRequest("GET", "/user/2", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	json.Unmarshal(response.Body.Bytes(), &answer)

	user := answer.Result

	data := user.(map[string]interface{})

	if data["name"] != "oleg" {
		t.Errorf("Expected user name to be 'oleg'. Got '%v'", data["name"])
	}

	if data["tg_id"] != 123.0 {
		t.Errorf("Expected user tgId to be '123'. Got '%v'", data["tg_id"])
	}

	if data["id"] != 2.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", data["id"])
	}
}