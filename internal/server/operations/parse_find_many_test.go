package operations

import (
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseFindMany(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	p := Parser{tx: appDB.NewTx()}

	var findManyTests = []struct {
		modelName  string
		jsonString string
		output     db.FindManyOperation
	}{
		// Basic String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew"
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					FieldCriteria: []db.FieldCriterion{
						db.FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
					},
				},
			},
		},
		// Multiple String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew",
				"lastName": "Wansley",
				"age": 32,
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					FieldCriteria: []db.FieldCriterion{
						db.FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
						db.FieldCriterion{
							Key: "Lastname",
							Val: "Wansley",
						},
						db.FieldCriterion{
							Key: "Age",
							Val: int64(32),
						},
					},
				},
			},
		},

		// Single Field To-One Relationship Criterion
		{
			modelName: "user",
			jsonString: `{ 
				"profile": { "text": "This is my bio.." }
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					RelationshipCriteria: []db.RelationshipCriterion{
						db.RelationshipCriterion{
							Relationship: db.User.Relationships["profile"],
							RelatedFieldCriteria: []db.FieldCriterion{
								db.FieldCriterion{
									Key: "Text",
									Val: "This is my bio..",
								},
							},
						},
					},
				},
			},
		},

		// Single Field To-One Relationship Criterion
		// with Nested Relationship Criterion
		{
			modelName: "user",
			jsonString: `{ 
				"profile": { 
					"text": "This is my bio..",
					"user": {
					  "firstName": "Andrew"
					}
				}
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					RelationshipCriteria: []db.RelationshipCriterion{
						db.RelationshipCriterion{
							Relationship: db.User.Relationships["profile"],
							RelatedRelationshipCriteria: []db.RelationshipCriterion{
								db.RelationshipCriterion{
									Relationship: db.Profile.Relationships["user"],
									RelatedFieldCriteria: []db.FieldCriterion{
										db.FieldCriterion{
											Key: "Firstname",
											Val: "Andrew",
										},
									},
								},
							},
							RelatedFieldCriteria: []db.FieldCriterion{
								db.FieldCriterion{
									Key: "Text",
									Val: "This is my bio..",
								},
							},
						},
					},
				},
			},
		},

		// Single Field To-Many "Some" Relationship Criterion
		{
			modelName:  "user",
			jsonString: `{ "posts": { "some": { "text": "This is my bio.." } } }`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					AggregateRelationshipCriteria: []db.AggregateRelationshipCriterion{
						db.AggregateRelationshipCriterion{
							Aggregation: db.Some,
							RelationshipCriterion: db.RelationshipCriterion{
								Relationship: db.User.Relationships["posts"],
								RelatedFieldCriteria: []db.FieldCriterion{
									db.FieldCriterion{
										Key: "Text",
										Val: "This is my bio..",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range findManyTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseFindMany(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		tFC := cmpopts.SortSlices(func(a, b db.FieldCriterion) bool {
			return a.Key < b.Key
		})
		tRC := cmpopts.SortSlices(func(a, b db.RelationshipCriterion) bool {
			return a.Relationship.TargetRel < b.Relationship.TargetRel
		})
		diff := cmp.Diff(testCase.output, parsedOp, tFC, tRC)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
