package db

import (
	"awans.org/aft/internal/model"
	"encoding/json"
	"github.com/go-test/deep"
	"testing"
)

func makeRecord(tx Tx, modelName string, jsonValue string) model.Record {
	st := tx.MakeRecord(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

type FindCase struct {
	st        model.Record
	modelName string
}

func TestCreateApply(t *testing.T) {
	appDB := New()
	AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := makeRecord(tx, "user", `{ 
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`)
	p := makeRecord(tx, "profile", `{
		"type":"profile",
		"text": "My bio.."}`)

	var createTests = []struct {
		operations []CreateOperation
		output     []FindCase
	}{
		// Simple Create
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{},
				},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
			},
		},
		// Nested Create
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{
						NestedCreateOperation{
							Relationship: User.Relationships["profile"],
							Record:       p,
							Nested:       []NestedOperation{},
						},
					},
				},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
				FindCase{
					st:        p,
					modelName: "profile",
				},
			},
		},
		// Nested connect
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: p,
					Nested: []NestedOperation{},
				},
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{
						NestedConnectOperation{
							Relationship: User.Relationships["profile"],
							// eventually need this to be a unique prop
							UniqueQuery: UniqueQuery{
								Key: "Text",
								Val: "My bio..",
							},
						},
					},
				},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
				FindCase{
					st:        p,
					modelName: "profile",
				},
			},
		},
	}
	for _, testCase := range createTests {
		// start each test on a fresh db
		appDB = New()
		AddSampleModels(appDB)
		tx = appDB.NewRWTx()
		for _, op := range testCase.operations {
			op.Apply(tx)
		}
		for _, findCase := range testCase.output {
			found, _ := findOneById(tx, findCase.modelName, findCase.st.Id())
			if diff := deep.Equal(found, findCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}