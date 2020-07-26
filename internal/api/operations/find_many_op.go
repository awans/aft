package operations

import (
	"awans.org/aft/internal/db"
)

type FieldCriterion struct {
	Key string
	Val interface{}
}

type AggregateRelationshipCriterion struct {
	RelationshipCriterion RelationshipCriterion
	Aggregation           db.Aggregation
}

type RelationshipCriterion struct {
	Relationship db.Relationship
	Where        Where
}

type Where struct {
	FieldCriteria                 []FieldCriterion
	RelationshipCriteria          []RelationshipCriterion
	AggregateRelationshipCriteria []AggregateRelationshipCriterion
	Or                            []Where
	And                           []Where
	Not                           []Where
}

type FindManyOperation struct {
	ModelID db.ID
	Where   Where
}

func (fc FieldCriterion) Matcher() db.Matcher {
	return db.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) []db.Record {
	q := buildQuery(tx, op)
	qrs := q.All()
	results := []db.Record{}
	for _, qr := range qrs {
		results = append(results, qr.Record)
	}
	return results
}

func buildQuery(tx db.Tx, op FindManyOperation) db.Q {
	root := tx.Ref(op.ModelID)
	clauses := handleWhere(tx, root, op.Where)
	q := tx.Query(root, clauses...)
	return q
}

func handleWhere(tx db.Tx, parent db.ModelRef, w Where) []db.QueryClause {
	clauses := []db.QueryClause{}
	for _, fc := range w.FieldCriteria {
		clauses = append(clauses, db.Filter(parent, fc.Matcher()))
	}
	for _, rc := range w.RelationshipCriteria {
		clauses = append(clauses, handleRC(tx, parent, rc)...)
	}
	for _, arc := range w.AggregateRelationshipCriteria {
		clauses = append(clauses, handleARC(tx, parent, arc)...)
	}

	var orBlocks []db.QBlock
	for _, or := range w.Or {
		orBlock := handleSetOpBranch(tx, parent, or)
		orBlocks = append(orBlocks, orBlock)
	}
	if len(orBlocks) > 0 {
		clauses = append(clauses, db.Or(parent, orBlocks...))
	}

	var andBlocks []db.QBlock
	for _, and := range w.And {
		andBlock := handleSetOpBranch(tx, parent, and)
		andBlocks = append(andBlocks, andBlock)
	}
	if len(andBlocks) > 0 {
		clauses = append(clauses, db.Union(parent, andBlocks...))
	}

	var notBlocks []db.QBlock
	for _, not := range w.Not {
		notBlock := handleSetOpBranch(tx, parent, not)
		notBlocks = append(notBlocks, notBlock)
	}
	if len(notBlocks) > 0 {
		clauses = append(clauses, db.Not(parent, notBlocks...))
	}
	return clauses
}

func handleSetOpBranch(tx db.Tx, parent db.ModelRef, w Where) db.QBlock {
	clauses := handleWhere(tx, parent, w)
	return db.Subquery(clauses...)
}

func handleRC(tx db.Tx, parent db.ModelRef, rc RelationshipCriterion) []db.QueryClause {
	child := tx.Ref(rc.Relationship.Target().ID())
	on := parent.Rel(rc.Relationship)
	j := db.Join(child, on)
	clauses := handleWhere(tx, child, rc.Where)
	clauses = append(clauses, j)
	return clauses
}

func handleARC(tx db.Tx, parent db.ModelRef, arc AggregateRelationshipCriterion) []db.QueryClause {
	child := tx.Ref(arc.RelationshipCriterion.Relationship.Target().ID())
	on := parent.Rel(arc.RelationshipCriterion.Relationship)

	j := db.Join(child, on)
	a := db.Aggregate(child, arc.Aggregation)
	clauses := handleWhere(tx, child, arc.RelationshipCriterion.Where)
	clauses = append(clauses, a)
	clauses = append(clauses, j)
	return clauses
}
