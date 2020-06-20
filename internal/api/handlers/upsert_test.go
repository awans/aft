package handlers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpsertUpdate(t *testing.T) {
	appDB := db.NewTest()
	eventbus := bus.New()
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.upsert", strings.NewReader(
		`{"update":{
			"firstName":"Chase"
		},
		"create":{
			"firstName":"Bob"
		},
		"where": {
			"firstName": "Andrew"
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := UpsertHandler{db: appDB, bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(bytes, &data)

	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Chase", objData["firstName"])
	assert.Equal(t, "Wansley", objData["lastName"])
}

func TestUpsertCreate(t *testing.T) {
	appDB := db.NewTest()
	eventbus := bus.New()
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.upsert", strings.NewReader(
		`{
		"create":{
			"firstName":"Bob"
		},
		"update":{
			"firstName":"Bob"
		},
		"where": {
			"firstName": "Bob"
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := UpsertHandler{db: appDB, bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(bytes, &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Bob", objData["firstName"])
}