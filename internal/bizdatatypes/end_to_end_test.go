package bizdatatypes

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var Andrew = db.MakeCoreDatatype(
	db.MakeID("46c0ee11-3943-452d-9420-925dd9be8208"),
	"andrew",
	db.StringStorage,
	AndrewCode)

var testingRox = db.MakeCoreDatatype(
	db.MakeID("9d792f82-018e-47d1-a2e5-a1b5b4822fd9"),
	"testingRox",
	db.StringStorage,
	testingRoxCode,
)

var AndrewCode = starlark.MakeStarlarkFunction(
	db.MakeID("a4615a60-afed-4f29-b674-e24f35618847"),
	"andrew",
	db.FromJSON,
	`def main(arg):
     if str(arg) == "Andrew":
         return "testing rox"
     fail("arg should be Andrew!!!")`)

var testingRoxCode = starlark.MakeStarlarkFunction(
	db.MakeID("5b0cfd40-4f3d-4890-b3a9-923ab8740043"),
	"testingRox",
	db.FromJSON,
	`def main(arg):
	return "testing rox"`)

var UserStarlark = db.MakeModel(
	db.MakeID("c1da149d-8ba0-429a-ab66-a8f2973c9e1e"),
	"starlark",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("a6e4e877-3e8c-4b17-9e80-0b896c0a1086"),
			"firstName",
			Andrew),
		db.MakeConcreteAttribute(
			db.MakeID("9bd6a83c-b805-4daf-b56f-6824f51fdbca"),
			"lastName",
			testingRox),
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

var starlarkTests = []struct {
	in          string
	out         string
	field       string
	shouldError bool
}{
	{`{"data":{
			"firstName":"Chase",
			"lastName":"Wansley"
		}
	}`, "", "", true},
	{`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley"
		}
	}`, "testing rox", "lastName", false},
}

func TestError(t *testing.T) {
	tt := starlarkTests[0]
	runner(t, tt.in, tt.out, tt.field, tt.shouldError)
}

func TestNoError(t *testing.T) {
	tt := starlarkTests[1]
	runner(t, tt.in, tt.out, tt.field, tt.shouldError)
}

func runner(t *testing.T, in, out, field string, shouldError bool) {
	appDB := db.New()
	db.AddSampleModels(appDB)

	sr := starlark.NewStarlarkRuntime()
	appDB.RegisterRuntime(sr)

	appDB.AddLiteral(AndrewCode)
	appDB.AddLiteral(testingRoxCode)
	appDB.AddLiteral(Andrew)
	appDB.AddLiteral(testingRox)
	appDB.AddLiteral(UserStarlark)

	eventbus := bus.New()
	req, err := http.NewRequest("POST", "/starlark.create", strings.NewReader(in))
	req = mux.SetURLVars(req, map[string]string{"modelName": "starlark"})
	cs := api.CreateHandler{DB: appDB, Bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	if shouldError {
		assert.Error(t, err)
		return
	}
	var data map[string]interface{}
	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(bytes, &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, out, objData[field])
}
