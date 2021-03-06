package parsers

import (
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
)

func (p Parser) consumeWhere(m db.Interface, keys api.Set, data map[string]interface{}) (operations.Where, error) {
	var w map[string]interface{}
	if v, ok := data["where"]; ok {
		w = v.(map[string]interface{})
		delete(keys, "where")
	}
	return p.ParseWhere(m, w)
}

func (p Parser) ParseWhere(m db.Interface, data map[string]interface{}) (q operations.Where, err error) {
	q = operations.Where{}
	fc, err := p.parseFieldCriteria(m, data)
	if err != nil {
		return
	}
	q.FieldCriteria = fc
	rc, err := p.parseSingleRelationshipCriteria(m, data)
	if err != nil {
		return
	}
	q.RelationshipCriteria = rc
	arc, err := p.parseAggregateRelationshipCriteria(m, data)
	if err != nil {
		return
	}
	q.AggregateRelationshipCriteria = arc

	if orVal, ok := data["OR"]; ok {
		var orQL []operations.Where
		orQL, err = p.parseCompositeQueryList(m, orVal)
		if err != nil {
			return
		}
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		var andQL []operations.Where
		andQL, err = p.parseCompositeQueryList(m, andVal)
		if err != nil {
			return
		}
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		var notQL []operations.Where
		notQL, err = p.parseCompositeQueryList(m, notVal)
		if err != nil {
			return
		}
		q.Not = notQL
	}
	return
}

func (p Parser) parseCompositeQueryList(m db.Interface, opVal interface{}) (ql []operations.Where, err error) {
	opList, ok := opVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected list as argument; got %T\n", opVal)
	}
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		var opQ operations.Where
		opQ, err = p.ParseWhere(m, opMap)
		if err != nil {
			return
		}
		ql = append(ql, opQ)
	}
	return
}

func (p Parser) parseSingleRelationshipCriteria(m db.Interface, data map[string]interface{}) (rcl []operations.RelationshipCriterion, err error) {
	rels, err := m.Relationships(p.Tx)
	if err != nil {
		return
	}
	for _, r := range rels {
		if !r.Multi() {
			if value, ok := data[r.Name()]; ok {
				var rc operations.RelationshipCriterion
				rc, err = p.parseRelationshipCriterion(r, value)
				if err != nil {
					return
				}
				rcl = append(rcl, rc)
			}
		}
	}
	return rcl, nil
}

func (p Parser) parseAggregateRelationshipCriteria(m db.Interface, data map[string]interface{}) (arcl []operations.AggregateRelationshipCriterion, err error) {
	rels, err := m.Relationships(p.Tx)
	if err != nil {
		return
	}
	for _, r := range rels {
		if r.Multi() {
			if value, ok := data[r.Name()]; ok {
				var arc operations.AggregateRelationshipCriterion
				arc, err = p.parseAggregateRelationshipCriterion(r, value)
				if err != nil {
					return
				}
				arcl = append(arcl, arc)
			}
		}
	}
	return arcl, nil
}

func (p Parser) parseFieldCriteria(m db.Interface, data map[string]interface{}) (fieldCriteria []operations.FieldCriterion, err error) {
	attrs, err := m.Attributes(p.Tx)
	if err != nil {
		return
	}
	rec, err := p.Tx.MakeRecord(m.ID())
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if value, ok := data[attr.Name()]; ok {
			var fc operations.FieldCriterion
			fc, err = parseFieldCriterion(p.Tx, attr, value, rec)
			if err != nil {
				return
			}
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	return
}

func parseFieldCriterion(tx db.Tx, a db.Attribute, value interface{}, rec db.Record) (fc operations.FieldCriterion, err error) {
	fieldName := db.JSONKeyToFieldName(a.Name())

	d := a.Datatype(tx)
	f, err := d.FromJSON(tx)
	if err != nil {
		err = fmt.Errorf("%w: err loading FromJSON of %v\n", err, d.Name())
		return
	}

	parsedValue, err := f.Call(tx.Context(), []interface{}{value, rec})

	fc = operations.FieldCriterion{
		// TODO handle function values like {startsWith}
		Key: fieldName,
		Val: parsedValue,
	}
	return
}

func (p Parser) parseAggregateRelationshipCriterion(r db.Relationship, value interface{}) (arc operations.AggregateRelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	if len(mapValue) > 1 {
		panic("too much data in parseAggregateRel")
	} else if len(mapValue) == 0 {
		panic("empty data in parseAggregateRel")
	}
	var ag db.Aggregation
	for k, v := range mapValue {
		switch k {
		case "some":
			ag = db.Some
		case "none":
			ag = db.None
		case "every":
			ag = db.Every
		default:
			return arc, fmt.Errorf("%w: %v", ErrAggregation, r.Name())
		}
		var rc operations.RelationshipCriterion
		rc, err = p.parseRelationshipCriterion(r, v)
		if err != nil {
			return
		}
		arc = operations.AggregateRelationshipCriterion{
			Aggregation:           ag,
			RelationshipCriterion: rc,
		}
	}
	return
}

func (p Parser) parseRelationshipCriterion(r db.Relationship, value interface{}) (rc operations.RelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	m := r.Target(p.Tx)
	fc, err := p.parseFieldCriteria(m, mapValue)
	if err != nil {
		return
	}
	rrc, err := p.parseSingleRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	arrc, err := p.parseAggregateRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	rc = operations.RelationshipCriterion{
		Relationship: r,
		Where: operations.Where{
			FieldCriteria:                 fc,
			RelationshipCriteria:          rrc,
			AggregateRelationshipCriteria: arrc,
		},
	}
	return
}
