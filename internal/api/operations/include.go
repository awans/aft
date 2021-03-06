package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
)

type Include struct {
	Includes []Inclusion
}

type Inclusion struct {
	Relationship   db.Relationship
	NestedFindMany FindArgs
}

func (i Inclusion) String() string {
	return fmt.Sprintf("inclusion{%v}\n", i.Relationship.Name())
}

func handleIncludes(tx db.Tx, parent db.ModelRef, i Include) []db.QueryClause {
	clauses := []db.QueryClause{}
	for _, inclusion := range i.Includes {
		clauses = append(clauses, handleInclusion(tx, parent, inclusion)...)
	}
	return clauses
}

func handleInclusion(tx db.Tx, parent db.ModelRef, i Inclusion) (clauses []db.QueryClause) {
	child := tx.Ref(i.Relationship.Target(tx).ID())
	j := db.LeftJoin(child, parent.Rel(i.Relationship))
	clauses = append(clauses, j)
	if i.Relationship.Multi() {
		a := db.Aggregate(child, db.Include)
		clauses = append(clauses, a)
	}
	clauses = append(clauses, handleFindMany(tx, child, i.NestedFindMany)...)
	return clauses
}
