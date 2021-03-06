package db

import (
	"fmt"

	"github.com/google/uuid"
)

type Matcher interface {
	Match(Record) (bool, error)
}

type op int

const (
	eq op = iota
	neq
	gt // not implemented
	lt // not implemented
)

type FieldMatcher struct {
	field string
	val   interface{}
	op    op
}

func (fm FieldMatcher) String() string {
	return fmt.Sprintf("match{%v == %v}", fm.field, fm.val)
}

// could be faster probably
func (fm FieldMatcher) Match(st Record) (bool, error) {
	candidate := st.MustGet(fm.field)
	comparison := fm.val
	return candidate == comparison, nil
}

func EqID(val ID) Matcher {
	u := uuid.UUID(val)
	return FieldMatcher{field: "id", val: u, op: eq}
}

func Eq(field string, val interface{}) Matcher {
	return FieldMatcher{field: field, val: val, op: eq}
}

type void struct{}

type idSetMatcher struct {
	ids map[uuid.UUID]void
}

func (im idSetMatcher) Match(r Record) (bool, error) {
	id := r.MustGet("id")
	_, ok := im.ids[id.(uuid.UUID)]
	return ok, nil
}

func (im idSetMatcher) String() string {
	return fmt.Sprintf("match{id in %v}", im.ids)
}

func IDIn(ids []ID) Matcher {
	hash := make(map[uuid.UUID]void)
	for _, id := range ids {
		u := uuid.UUID(id)
		hash[u] = void{}
	}

	return idSetMatcher{ids: hash}
}

type AndMatcher struct {
	inner []Matcher
}

func (am AndMatcher) Match(st Record) (bool, error) {
	for _, m := range am.inner {
		match, err := m.Match(st)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (am AndMatcher) String() string {
	return fmt.Sprintf("match{and(%v)}", am.inner)
}

func And(matchers ...Matcher) Matcher {
	return AndMatcher{inner: matchers}
}

type FalseMatcher struct {
}

func (FalseMatcher) Match(Record) (bool, error) {
	return false, nil
}

func False() Matcher {
	return FalseMatcher{}
}

func IsModel(modelID ID) Matcher {
	return ModelMatcher{
		modelID: modelID,
	}
}

func IsNotModel(modelID ID) Matcher {
	return ModelNonMatcher{
		modelID: modelID,
	}
}

type ModelNonMatcher struct {
	modelID ID
}

func (m ModelNonMatcher) Match(rec Record) (bool, error) {
	return m.modelID != rec.InterfaceID(), nil
}

func (m ModelNonMatcher) String() string {
	return fmt.Sprintf("match{notModel(%v)}", m.modelID)
}

type ModelMatcher struct {
	modelID ID
}

func (m ModelMatcher) Match(rec Record) (bool, error) {
	return m.modelID == rec.InterfaceID(), nil
}

func (m ModelMatcher) String() string {
	return fmt.Sprintf("match{isModel(%v)}", m.modelID)
}
