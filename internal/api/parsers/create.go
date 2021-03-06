package parsers

import (
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
)

func (p Parser) ParseCreate(modelName string, args map[string]interface{}) (op operations.CreateOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeCreateRel(m, data)
	if err != nil {
		return
	}
	inc, sel, err := p.consumeIncludeOrSelect(m, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.CreateOperation{
		ModelID: m.ID(),
		Data:    data,
		Nested:  nested,
		FindArgs: operations.FindArgs{
			Include: inc,
			Select:  sel,
		},
	}, nil
}

// parseNestedCreate handles parsing a single nested create operation
func (p Parser) parseNestedCreate(rel db.Relationship, data map[string]interface{}) (op operations.NestedOperation, err error) {
	m, err := p.resolveInterface(rel.Target(p.Tx), data)
	if err != nil {
		return
	}
	nested, err := p.consumeCreateRel(m, data)
	if err != nil {
		return
	}
	return operations.NestedCreateOperation{Relationship: rel, Model: m, Data: data, Nested: nested}, nil
}

func (p Parser) resolveInterface(iface db.Interface, data map[string]interface{}) (m db.Model, err error) {
	m, ok := iface.(db.Model)
	if ok {
		return
	} else {
		typeValue, ok := data["type"]
		if !ok {
			err = fmt.Errorf("%w: %v", ErrInvalidModel, typeValue)
		}
		modelName, ok := typeValue.(string)
		if ok {
			m, err = p.Tx.Schema().GetModel(modelName)
		} else {
			err = fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
		}
	}
	return
}

// consumeCreateRel handles all nested operations on a given create (top-level or nested)
// so data contains keys for each attribute or relationship in the op
func (p Parser) consumeCreateRel(m db.Model, data map[string]interface{}) (nested []operations.NestedOperation, err error) {
	unusedKeys := make(api.Set)
	for k := range data {
		unusedKeys[k] = api.Void{}
	}

	// delete all attributes from unusedKeys
	// because we assume they've been handled by a peer function
	attrs, err := m.Attributes(p.Tx)
	if err != nil {
		return nil, err
	}
	for _, attr := range attrs {
		if _, ok := unusedKeys[attr.Name()]; ok {
			delete(unusedKeys, attr.Name())
		}
	}

	rels, err := m.Relationships(p.Tx)
	if err != nil {
		return
	}
	nested = []operations.NestedOperation{}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseNestedCreateRelationship(r, data)
		if err != nil {
			return nested, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
			delete(data, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return nested, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	return
}

// parseNestedCreateRelationship parses a single relationship key under a create, returning a slice of nested operations
// it takes the same 'data' as consumeCreateRel, so we start by indexing into it
func (p Parser) parseNestedCreateRelationship(r db.Relationship, data map[string]interface{}) ([]operations.NestedOperation, bool, error) {
	nestedOpMap, ok := data[r.Name()].(map[string]interface{})
	if !ok {
		_, isValue := data[r.Name()]
		if !isValue {
			return []operations.NestedOperation{}, false, nil
		}

		return []operations.NestedOperation{}, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, data)
	}

	var nested []operations.NestedOperation
	for k, val := range nestedOpMap {
		opList, err := listify(val)

		if err != nil {
			return nil, false, err
		}

		for _, op := range opList {
			nestedOp, ok := op.(map[string]interface{})
			if !ok {
				return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
			}
			switch k {
			case "set":
				nestedSet, err := p.parseNestedSet(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedSet)
			case "connect":
				nestedConnect, err := p.parseNestedConnect(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedConnect)
			case "create":
				nestedCreate, err := p.parseNestedCreate(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedCreate)
			}
		}
	}

	return nested, true, nil
}
