package operations

import (
	"fmt"
	"testing"

	"awans.org/aft/internal/db"
)

func TestCase(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()

	operation := FindManyOperation{
		ModelID: db.RelationshipInterface.ID(),
		FindArgs: FindArgs{
			Where: Where{},
			Case: Case{
				Entries: []CaseEntry{
					CaseEntry{
						ModelID: db.ConcreteRelationshipModel.ID(),
						Include: Include{
							[]Inclusion{
								Inclusion{
									Relationship:   db.ConcreteRelationshipTarget.Load(tx),
									NestedFindMany: FindArgs{},
								},
							},
						},
					},
					CaseEntry{
						ModelID: db.ReverseRelationshipModel.ID(),
						Include: Include{
							[]Inclusion{
								Inclusion{
									Relationship:   db.ReverseRelationshipReferencing.Load(tx),
									NestedFindMany: FindArgs{},
								},
							},
						},
					},
				},
			},
		},
	}

	results, _ := operation.Apply(tx)
	for _, res := range results {
		if res.Record.InterfaceID() == db.ConcreteRelationshipModel.ID() {
			_, ok := res.ToOne["target"]
			if !ok {
				err := fmt.Errorf("Didn't get target for %v\n", res)
				t.Error(err)
			}
			_, ok = res.ToOne["referencing"]
			if ok {
				err := fmt.Errorf("Got referencing for %v\n", res)
				t.Error(err)
			}
		} else if res.Record.InterfaceID() == db.ReverseRelationshipModel.ID() {
			_, ok := res.ToOne["target"]
			if ok {
				err := fmt.Errorf("Get target for %v\n", res)
				t.Error(err)
			}
			_, ok = res.ToOne["referencing"]
			if !ok {
				err := fmt.Errorf("Didn't get referencing for %v\n", res)
				t.Error(err)
			}
		}
	}

}
