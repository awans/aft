package db

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

var (
	ErrInvalidRelationship = fmt.Errorf("%w: invalid relationship", ErrData)
)

type Attribute struct {
	AttrType Datatype
	Id       uuid.UUID
}

type RelType int64

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Relationship struct {
	Id           uuid.UUID
	LeftBinding  RelType
	RightBinding RelType
	LeftModelId  uuid.UUID
	RightModelId uuid.UUID
	LeftName     string
	RightName    string
}

func (r Relationship) Left() Binding {
	return Binding{Relationship: r, Left: true}
}

func (r Relationship) Right() Binding {
	return Binding{Relationship: r, Left: false}
}

func JsonKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vId", strings.Title(strings.ToLower(key)))
}

type Model struct {
	Id                 uuid.UUID
	Name               string
	Attributes         map[string]Attribute
	LeftRelationships  []Relationship
	RightRelationships []Relationship
}

func (m Model) AttributeByName(name string) Attribute {
	a, ok := m.Attributes[name]
	if !ok {
		a, ok = SystemAttrs[name]
	}
	return a
}

func (m Model) GetBinding(name string) (Binding, error) {
	for _, b := range m.Bindings() {
		if b.Name() == name {
			return b, nil
		}
	}
	return Binding{}, ErrInvalidRelationship
}

func (m Model) Bindings() []Binding {
	var bs []Binding
	for _, r := range m.LeftRelationships {
		bs = append(bs, Binding{Relationship: r, Left: true})
	}
	for _, r := range m.RightRelationships {
		bs = append(bs, Binding{Relationship: r, Left: false})
	}
	return bs
}

type Binding struct {
	Relationship Relationship
	Left         bool
}

func (b Binding) HasField() bool {
	if b.Left {
		return (b.Relationship.LeftBinding == BelongsTo)
	} else {
		return (b.Relationship.RightBinding == BelongsTo)
	}
}

func (b Binding) Name() string {
	if b.Left {
		return b.Relationship.LeftName
	} else {
		return b.Relationship.RightName
	}
}

func (b Binding) ModelId() uuid.UUID {
	if b.Left {
		return b.Relationship.LeftModelId
	} else {
		return b.Relationship.RightModelId
	}
}

func (b Binding) Dual() Binding {
	return Binding{Relationship: b.Relationship, Left: !b.Left}
}

func (b Binding) RelType() RelType {
	if b.Left {
		return b.Relationship.LeftBinding
	} else {
		return b.Relationship.RightBinding
	}
}
