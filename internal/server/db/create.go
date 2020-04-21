package db

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
)

func getId(st interface{}) uuid.UUID {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id
}

func (op CreateOperation) Apply(db DB) interface{} {
	db.h.Insert(op.Struct)
	for _, no := range op.Nested {
		no.ApplyNested(db, op.Struct)
	}
	return op.Struct
}

func getBackref(db DB, rel model.Relationship) model.Relationship {
	m := db.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func setFK(st interface{}, key string, id uuid.UUID) {
	fieldName := model.JsonKeyToRelFieldName(key)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	v := reflect.ValueOf(id)
	field.Set(v)
}

func connect(db DB, from interface{}, fromRel model.Relationship, to interface{}) {
	toRel := getBackref(db, fromRel)
	if fromRel.Type == model.BelongsTo && (toRel.Type == model.HasOne || toRel.Type == model.HasMany) {
		// FK from
		setFK(from, toRel.TargetRel, getId(to))
	} else if toRel.Type == model.BelongsTo && (fromRel.Type == model.HasOne || fromRel.Type == model.HasMany) {
		// FK to
		setFK(to, fromRel.TargetRel, getId(from))
	} else if toRel.Type == model.HasManyAndBelongsToMany && fromRel.Type == model.HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		panic("Trying to connect invalid relationship")
	}
}

func (op NestedCreateOperation) ApplyNested(db DB, parent interface{}) {
	connect(db, parent, op.Relationship, op.Struct)
	db.h.Insert(op.Struct)
	db.h.Insert(parent)
}

func findOne(db DB, modelName string, uq UniqueQuery) interface{} {
	val, err := db.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	if err != nil {
		panic("FindOne failed")
	}
	return val
}

func findOneById(db DB, modelName string, id uuid.UUID) interface{} {
	return findOne(db, modelName, UniqueQuery{Key: "Id", Val: id})
}

func (op NestedConnectOperation) ApplyNested(db DB, parent interface{}) {
	modelName := op.Relationship.TargetModel
	st := findOne(db, modelName, op.UniqueQuery)
	connect(db, parent, op.Relationship, st)
	db.h.Insert(st)
	db.h.Insert(parent)
}