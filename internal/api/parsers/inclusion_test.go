package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseInclude(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()
	p := Parser{Tx: tx}

	u, _ := tx.MakeRecord(db.User.ID())
	up, _ := u.Interface().RelationshipByName("profile")
	var inclusionTests = []struct {
		model      db.Interface
		jsonString string
		output     operations.Include
	}{
		// Simple Include
		{
			model: u.Interface(),
			jsonString: `{
			   "profile": true
			}`,
			output: operations.Include{
				Includes: []operations.Inclusion{
					operations.Inclusion{
						Relationship: up,
					},
				},
			},
		},
	}
	for _, testCase := range inclusionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.parseInclude(testCase.model, data)
		if err != nil {
			t.Error(err)
		}
		diff := cmp.Diff(testCase.output, parsedOp, CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
