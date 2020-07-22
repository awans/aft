package operations

import (
	"awans.org/aft/internal/db"
)

type FindManyArgs struct {
	Where   Where
	Include Include
	// Add Select
}

type FindOneOperation struct {
	ModelID db.ID
	FindManyArgs FindManyArgs
}

type FindManyOperation struct {
	ModelID db.ID
	FindManyArgs FindManyArgs
}

type CreateOperation struct {
	Record  db.Record
	FindManyArgs FindManyArgs
	Nested  []NestedOperation
}

type UpdateOperation struct {
	ModelID db.ID
	FindManyArgs FindManyArgs
	Data    map[string]interface{}
	Nested []NestedOperation
}

type UpsertOperation struct {
	ModelID      db.ID
	FindManyArgs FindManyArgs
	Create       db.Record
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
}

type DeleteOperation struct {
	ModelID db.ID
	FindManyArgs FindManyArgs
	Nested  []NestedOperation
}

type UpdateManyOperation struct {
	ModelID db.ID
	Where   Where
	Data    map[string]interface{}
	Nested  []NestedOperation
}

type DeleteManyOperation struct {
	ModelID db.ID
	Where   Where
	Nested  []NestedOperation
}

type CountOperation struct {
	ModelID db.ID
	Where   Where
}

<<<<<<< HEAD
=======
type CountOperation struct {
	ModelID db.ID
	Where   Where
}

>>>>>>> 6b57a71... Refactor API to use common code
//Nested operations
type NestedOperation interface {
	ApplyNested(db.RWTx) error
}

type NestedCreateOperation struct {
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedConnectOperation struct {
	Relationship db.Relationship
	Where        Where
}

type NestedUpdateOperation struct {
	Where        Where
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedDeleteOperation struct {
	Where        Where
	Relationship db.Relationship
	Nested       []NestedOperation
}

type NestedUpdateManyOperation struct {
	Where        Where
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedDeleteManyOperation struct {
	Where        Where
	Relationship db.Relationship
	Nested       []NestedOperation
}

type NestedUpsertOperation struct {
	Relationship db.Relationship
	Where        Where
	Create       map[string]interface{}
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
}
