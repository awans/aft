package operations

import (
	"encoding/gob"

	"awans.org/aft/internal/db"
)

type FieldCriterion struct {
	Key string
	Val interface{}
}

func init() {
	gob.Register(FieldCriterion{})
	gob.Register(AggregateRelationshipCriterion{})
	gob.Register(RelationshipCriterion{})
	gob.Register(Where{})
}

type AggregateRelationshipCriterion struct {
	RelationshipCriterion RelationshipCriterion `json:",omitempty"`
	Aggregation           db.Aggregation
}

type RelationshipCriterion struct {
	Relationship db.Relationship
	Where        Where
}

type Where struct {
	FieldCriteria                 []FieldCriterion                 `json:",omitempty"`
	RelationshipCriteria          []RelationshipCriterion          `json:",omitempty"`
	AggregateRelationshipCriteria []AggregateRelationshipCriterion `json:",omitempty"`
	Or                            []Where                          `json:",omitempty"`
	And                           []Where                          `json:",omitempty"`
	Not                           []Where                          `json:",omitempty"`
}

func (fc FieldCriterion) Matcher() db.Matcher {
	return db.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) ([]*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleFindMany(tx, root, op.FindArgs)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	return qrs, nil
}

func handleFindMany(tx db.Tx, parent db.ModelRef, fm FindArgs) []db.QueryClause {
	clauses := HandleWhere(tx, parent, fm.Where)
	clauses = append(clauses, handleIncludes(tx, parent, fm.Include)...)
	clauses = append(clauses, handleSelects(tx, parent, fm.Select)...)
	clauses = append(clauses, handleCase(tx, parent, fm.Case)...)

	if fm.Take != -1 {
		clauses = append(clauses, db.Limit(parent, fm.Take))
	}
	clauses = append(clauses, db.Offset(fm.Skip, parent))
	clauses = append(clauses, db.Order(parent, fm.Order))
	return clauses
}

func HandleWhere(tx db.Tx, parent db.ModelRef, w Where) []db.QueryClause {
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

	var orBlocks []db.Q
	for _, or := range w.Or {
		orBlock := handleSetOpBranch(tx, parent, or)
		orBlocks = append(orBlocks, orBlock)
	}
	if len(orBlocks) > 0 {
		clauses = append(clauses, db.Or(parent, orBlocks...))
	}

	var andBlocks []db.Q
	for _, and := range w.And {
		andBlock := handleSetOpBranch(tx, parent, and)
		andBlocks = append(andBlocks, andBlock)
	}
	if len(andBlocks) > 0 {
		clauses = append(clauses, db.Intersection(parent, andBlocks...))
	}

	var notBlocks []db.Q
	for _, not := range w.Not {
		notBlock := handleSetOpBranch(tx, parent, not)
		notBlocks = append(notBlocks, notBlock)
	}
	if len(notBlocks) > 0 {
		clauses = append(clauses, db.Not(parent, notBlocks...))
	}
	return clauses
}

func handleSetOpBranch(tx db.Tx, parent db.ModelRef, w Where) db.Q {
	clauses := HandleWhere(tx, parent, w)
	return tx.Subquery(clauses...)
}

func handleRC(tx db.Tx, parent db.ModelRef, rc RelationshipCriterion) []db.QueryClause {
	child := tx.Ref(rc.Relationship.Target(tx).ID())
	on := parent.Rel(rc.Relationship)
	j := db.Join(child, on)
	clauses := HandleWhere(tx, child, rc.Where)
	clauses = append(clauses, j)
	return clauses
}

func handleARC(tx db.Tx, parent db.ModelRef, arc AggregateRelationshipCriterion) []db.QueryClause {
	child := tx.Ref(arc.RelationshipCriterion.Relationship.Target(tx).ID())
	on := parent.Rel(arc.RelationshipCriterion.Relationship)

	j := db.Join(child, on)
	a := db.Aggregate(child, arc.Aggregation)
	clauses := HandleWhere(tx, child, arc.RelationshipCriterion.Where)
	clauses = append(clauses, a)
	clauses = append(clauses, j)
	return clauses
}

//logic to get the appropriate results from a relationship where a filter potentially exists
func handleRelationshipWhere(tx db.Tx, parent db.ModelRef, parents []*db.QueryResult, rel db.Relationship, where Where) (outs []*db.QueryResult, child db.ModelRef) {
	child = tx.Ref(rel.Target(tx).ID())

	ids := []db.ID{}
	for _, rec := range parents {
		ids = append(ids, rec.Record.ID())
	}
	var clauses []db.QueryClause
	clauses = append(clauses, db.Filter(parent, db.IDIn(ids)))
	clauses = append(clauses, db.Join(child, parent.Rel(rel)))
	if rel.Multi() {
		clauses = append(clauses, db.Aggregate(child, db.Some))
	}
	clauses = append(clauses, HandleWhere(tx, child, where)...)

	q := tx.Query(parent, clauses...)
	parentsWithChildren := q.All()

	for _, o := range parentsWithChildren {
		if rel.Multi() {
			outs = append(outs, o.GetChildRelMany(rel)...)
		} else {
			outs = append(outs, o.GetChildRelOne(rel))
		}
	}
	return outs, child
}
