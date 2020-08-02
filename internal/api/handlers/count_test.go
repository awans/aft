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

func TestCountSimple(t *testing.T) {
	appDB := db.NewTest()
	eventbus := bus.New()
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{ "id": "f90e1855-dbaa-4385-9929-20efe86cccb2", "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	jsonString2 := `{ "id": "9dd0a0c6-7e41-4107-9529-e75a5c7135cf", "firstName":"Chase", "lastName":"Hensel", "age": 33, "emailAddress":"chase.hensel@gmail.com"}`
	u2 := api.MakeRecord(appDB.NewTx(), "user", jsonString2)
	tx.Insert(u2)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.count", strings.NewReader(
		`{"where": {
			"age": 32
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := UpdateManyHandler{db: appDB, bus: eventbus}
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
	assert.Equal(t, 1.0, data["count"])
}
