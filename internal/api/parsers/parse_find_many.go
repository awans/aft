package parsers

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/api/operations"
)

func (p Parser) ParseFindMany(modelName string, data map[string]interface{}) (op operations.FindManyOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}
	q, err := p.ParseWhere(modelName, data)
	if err != nil {
		return
	}

	op = operations.FindManyOperation{
		Where:   q,
		ModelID: m.ID(),
	}
	return op, nil
}

func (p Parser) parseCompositeQueryList(modelName string, opVal interface{}) (ql []operations.Where, err error) {
	opList := opVal.([]interface{})
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		var opQ operations.Where
		opQ, err = p.ParseWhere(modelName, opMap)
		if err != nil {
			return
		}
		ql = append(ql, opQ)
	}
	return
}

func (p Parser) ParseWhere(modelName string, data map[string]interface{}) (q operations.Where, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}
	q = operations.Where{}
	fc, err := parseFieldCriteria(m, data)
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
		orQL, err = p.parseCompositeQueryList(modelName, orVal)
		if err != nil {
			return
		}
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		var andQL []operations.Where
		andQL, err = p.parseCompositeQueryList(modelName, andVal)
		if err != nil {
			return
		}
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		var notQL []operations.Where
		notQL, err = p.parseCompositeQueryList(modelName, notVal)
		if err != nil {
			return
		}
		q.Not = notQL
	}
	return
}

func (p Parser) parseSingleRelationshipCriteria(m db.Interface, data map[string]interface{}) (rcl []operations.RelationshipCriterion, err error) {
	rels, err := m.Relationships()
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
	rels, err := m.Relationships()
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

func parseFieldCriteria(m db.Interface, data map[string]interface{}) (fieldCriteria []operations.FieldCriterion, err error) {
	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if value, ok := data[attr.Name()]; ok {
			var fc operations.FieldCriterion
			fc, err = parseFieldCriterion(attr, value)
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	return
}

func parseFieldCriterion(a db.Attribute, value interface{}) (fc operations.FieldCriterion, err error) {
	fieldName := db.JSONKeyToFieldName(a.Name())

	d := a.Datatype()
	f, err := d.FromJSON()
	if err != nil {
		return
	}

	parsedValue, err := f.Call(value)

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
			panic("Bad aggregation")
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
	m := r.Target()
	fc, err := parseFieldCriteria(m, mapValue)
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