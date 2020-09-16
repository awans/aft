package db

import (
	"errors"

	"github.com/google/uuid"
)

// Utility methods
// filterEmtpy -- copies a []*QR into a new slice, removing any that are isEmpty()

func filterEmpty(results []*QueryResult) []*QueryResult {
	filtered := []*QueryResult{}
	for _, qr := range results {
		if !qr.isEmpty() {
			filtered = append(filtered, qr)
		}
	}
	return filtered
}

func copyResultsShallow(results []*QueryResult) []*QueryResult {
	var copied []*QueryResult
	for _, r := range results {
		copied = append(copied, &QueryResult{Record: r.Record})
	}
	return copied
}

func applyMatcher(results []*QueryResult, matcher Matcher) []*QueryResult {
	for _, result := range results {
		if !result.isEmpty() {
			match, _ := matcher.Match(result.Record)
			if !match {
				result.Empty()
			}
		}
	}
	return results
}

// Entrypoint

func (qb Q) runBlockRoot(tx *holdTx) []*QueryResult {
	matchers := qb.Filters[qb.Root.AliasID]
	outer := qb.performScan(tx, qb.Root.InterfaceID, And(matchers...))
	results := qb.runBlock(tx, outer, qb.Root.AliasID)
	results = filterEmpty(results)
	qb.projectFields(qb.Selections[qb.Root.AliasID], results)
	return results
}

func (qb Q) projectFields(selection Selection, qrs []*QueryResult) {
	if selection.selecting {
		for _, qr := range qrs {
			qr.HideAll()
			for k, _ := range selection.fields {
				qr.Show(k)
			}
		}
	}
}

func (qb Q) runBlockNested(tx *holdTx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	matchers, ok := qb.Filters[aliasID]
	if ok {
		outer = applyMatcher(outer, And(matchers...))
	}
	qb.projectFields(qb.Selections[aliasID], outer)
	return qb.runBlock(tx, outer, aliasID)
}

func (qb Q) runBlock(tx *holdTx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	results := qb.performJoins(tx, outer, aliasID)
	results = qb.performSetOps(tx, results, aliasID)
	results = qb.performCases(tx, results, aliasID)
	return results
}

func (qb Q) performCases(tx *holdTx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	cases, ok := qb.Cases[aliasID]
	if ok {
		for _, c := range cases {
			outer = qb.performCase(tx, outer, c, aliasID)
		}
	}
	return outer
}

func (qb Q) performCase(tx *holdTx, outer []*QueryResult, c CaseOperation, aliasID uuid.UUID) []*QueryResult {
	mid := c.Of.InterfaceID
	inner := []*QueryResult{}
	for _, qr := range outer {
		if !qr.isEmpty() && qr.Record.Interface().ID() == mid {
			inner = append(inner, qr)
		}
	}
	qb.runBlock(tx, inner, c.Of.AliasID)
	return outer
}

func (qb Q) performSetOps(tx *holdTx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	setops, ok := qb.SetOps[aliasID]
	if ok {
		for _, s := range setops {
			outer = qb.performSetOp(tx, outer, s, aliasID)
		}
	}
	return outer
}

func orResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {

	for i, o := range original {
		any := false
		for j := range set {
			r := set[j][i]
			if !r.isEmpty() {
				any = true
				break
			}
		}
		if !any {
			o.Empty()
		}
	}
	return original
}

func andResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {
	for i, o := range original {
		all := true
		for j := range set {
			r := set[j][i]
			if r.isEmpty() {
				all = false
			}
		}
		if !all {
			o.Empty()
		}
	}
	return original
}

func notResults(original []*QueryResult, set [][]*QueryResult) []*QueryResult {
	for i, o := range original {
		any := false
		for j := range set {
			r := set[j][i]
			if !r.isEmpty() {
				any = false
			}
		}
		if any {
			o.Empty()
		}
	}
	return original
}

func (qb Q) performSetOp(tx *holdTx, outer []*QueryResult, op SetOperation, aliasID uuid.UUID) []*QueryResult {
	original := copyResultsShallow(outer)
	var set [][]*QueryResult
	for _, b := range op.Branches {
		branchCopy := copyResultsShallow(outer)
		branchResults := b.runBlockNested(tx, branchCopy, aliasID)
		set = append(set, branchResults)
	}
	switch op.op {
	case or:
		return orResults(original, set)
	case and:
		return andResults(original, set)
	case not:
		return notResults(original, set)
	default:
		panic("invalid set op")
	}
}

func (qb Q) performScan(tx *holdTx, modeID ID, matcher Matcher) []*QueryResult {
	recs, _ := tx.FindMany(qb.Root.InterfaceID, matcher)
	var results []*QueryResult
	for _, rec := range recs {
		results = append(results, &QueryResult{Record: rec})
	}
	return results
}

func (qb Q) performJoins(tx *holdTx, outer []*QueryResult, aliasID uuid.UUID) []*QueryResult {
	for _, j := range qb.Joins[aliasID] {
		toOne := j.IsToOne()

		if toOne {
			outer = qb.performJoinOne(tx, outer, j)
		} else {
			outer = qb.performJoinMany(tx, outer, j)
		}
	}
	return outer
}

func (qb Q) performJoinOne(tx *holdTx, outer []*QueryResult, j JoinOperation) []*QueryResult {
	var inner []*QueryResult
	key := j.Key()
	matchers := qb.Filters[j.To.AliasID]

	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedOne(tx, r.Record, j, And(matchers...))
			inner = append(inner, qr)
		} else {
			inner = append(inner, &QueryResult{})
		}
	}

	qb.runBlock(tx, inner, j.To.AliasID)

	if j.jt == innerJoin {
		// inner join
		for i := range outer {
			if !inner[i].isEmpty() {
				toOneMap := outer[i].ToOne
				if toOneMap == nil {
					outer[i].ToOne = map[string]*QueryResult{key: inner[i]}
				} else {
					outer[i].ToOne[key] = inner[i]
				}

			} else {
				outer[i].Empty()
			}
		}
		qb.projectFields(qb.Selections[j.To.AliasID], outer)
		return outer
	} else {
		// left join
		for i := range outer {
			if !inner[i].isEmpty() {
				toOneMap := outer[i].ToOne
				if toOneMap == nil {
					outer[i].ToOne = map[string]*QueryResult{key: inner[i]}
				} else {
					outer[i].ToOne[key] = inner[i]
				}
			}
		}
		qb.projectFields(qb.Selections[j.To.AliasID], outer)
		return outer
	}
}

func getRelatedOne(tx *holdTx, rec Record, j JoinOperation, matcher Matcher) *QueryResult {
	hit, err := tx.h.GetLinkedOne(rec.ID(), j.on.rel.ID())
	// not too sure about this..
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return &QueryResult{Record: nil}
		} else {
			panic(err)
		}
	}
	ok, err := matcher.Match(hit)
	if err != nil {
		panic(err)
	}
	if ok {
		return &QueryResult{Record: hit}
	}
	return &QueryResult{Record: nil}
}

func (qb Q) performJoinManySomeOrInclude(tx *holdTx, outer []*QueryResult, j JoinOperation, a Aggregation) []*QueryResult {
	key := j.Key()
	matchers := qb.Filters[j.To.AliasID]
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, And(matchers...))
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
		}
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[ID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins and any cases on this join
	// we're passing pointers so stuf'll get modified in-place in the dict
	qb.runBlock(tx, uniqValues, j.To.AliasID)

	// merge 'em back, filtering if none made it back
	for i := range outer {
		joinedSet := inner[i]
		var populatedJoinedSet []*QueryResult
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				populatedJoinedSet = append(populatedJoinedSet, populated)
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		isEmpty := true
		for _, v := range populatedJoinedSet {
			if !v.isEmpty() {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			// if this is a Some aggregation
			// blank out the parent record
			if a == Some {
				outer[i].Empty()
			} else {
				outer[i].SetChildRelMany(key, []*QueryResult{})
			}
		} else {
			dict := outer[i].ToMany
			if dict != nil {
				dict[key] = populatedJoinedSet
			} else {
				outer[i].SetChildRelMany(key, populatedJoinedSet)
			}
		}

	}
	qb.projectFields(qb.Selections[j.To.AliasID], outer)
	return outer
}

func (qb Q) performJoinManyNone(tx *holdTx, outer []*QueryResult, j JoinOperation, a Aggregation) []*QueryResult {
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, nil)
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
		}
	}

	matchers := qb.Filters[j.To.AliasID]
	matcher := And(matchers...)

	// apply the local filtering criteria
	// eagerly so we can maybe avoid doing some extra joins
	var filtered [][]*QueryResult
	for _, group := range inner {
		none := true
		for _, result := range group {
			match, _ := matcher.Match(result.Record)
			if match {
				none = false
				break
			}
		}
		if none {
			filtered = append(filtered, group)
		} else {
			filtered = append(filtered, []*QueryResult{})
		}
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[ID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins and any cases on this join
	// we're passing pointers so stuff'll get modified in-place in the dict
	qb.runBlock(tx, uniqValues, j.To.AliasID)

	// merge 'em back, filtering if any make it back
	for i := range outer {
		joinedSet := inner[i]
		none := true
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				if !populated.isEmpty() {
					none = false
					break
				}
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		if !none {
			// if this is an None aggregation
			// blank out the parent record
			if a == None {
				outer[i].Empty()
			}
		}
	}
	qb.projectFields(qb.Selections[j.To.AliasID], outer)
	return outer
}

func (qb Q) performJoinManyEvery(tx *holdTx, outer []*QueryResult, j JoinOperation, a Aggregation) []*QueryResult {
	key := j.Key()
	var inner [][]*QueryResult
	for _, r := range outer {
		if !r.isEmpty() {
			qr := getRelatedMany(tx, r.Record, j, nil)
			inner = append(inner, qr)
		} else {
			inner = append(inner, []*QueryResult{})
		}
	}

	matchers := qb.Filters[j.To.AliasID]
	matcher := And(matchers...)

	// apply the local filtering criteria
	var filtered [][]*QueryResult
	// eagerly so we can maybe avoid doing some extra joins
	for _, group := range inner {
		every := true
		for _, result := range group {
			match, _ := matcher.Match(result.Record)
			if !match {
				every = false
				break
			}
		}
		if every {
			filtered = append(filtered, group)
			qb.projectFields(qb.Selections[j.To.AliasID], group)
		} else {
			filtered = append(filtered, []*QueryResult{})
		}
	}

	// to prevent explosion, we first merge by unique records
	// and then expand out
	uniq := map[ID]*QueryResult{}
	for _, group := range inner {
		for _, result := range group {
			uniq[result.Record.ID()] = result
		}
	}

	// just copy out the unique values
	var uniqValues []*QueryResult
	for _, uniqVal := range uniq {
		uniqValues = append(uniqValues, uniqVal)
	}

	// do all of the child joins
	// we're passing pointers so stuf'll get modified in-place in the dict
	qb.runBlock(tx, uniqValues, j.To.AliasID)

	// merge 'em back, filtering if any didn't make it back
	for i := range outer {
		joinedSet := inner[i]
		every := true
		var populatedJoinedSet []*QueryResult
		for _, joined := range joinedSet {
			if !joined.isEmpty() {
				populated := uniq[joined.Record.ID()]
				if populated.isEmpty() {
					every = false
					break
				}
				populatedJoinedSet = append(populatedJoinedSet, populated)
			} else {
				every = false
				break
			}
		}

		// okay cool, now populated joined set contains fully joined values
		// for this input record
		// apply the aggregation!
		if !every {
			// if this is an Every aggregation
			// blank out the parent record
			if a == Every {
				outer[i].Empty()
			}
		} else {
			outer[i].SetChildRelMany(key, populatedJoinedSet)
		}

	}
	return outer
}

func (qb Q) performJoinMany(tx *holdTx, outer []*QueryResult, j JoinOperation) []*QueryResult {
	agg, ok := qb.Aggregations[j.To.AliasID]
	if !ok {
		agg = Include
	}
	switch agg {
	case Some, Include:
		return qb.performJoinManySomeOrInclude(tx, outer, j, agg)
	case Every:
		return qb.performJoinManyEvery(tx, outer, j, agg)
	case None:
		return qb.performJoinManyNone(tx, outer, j, agg)
	}
	panic("not implemented")
}

func getRelatedMany(tx *holdTx, rec Record, j JoinOperation, matcher Matcher) []*QueryResult {
	if rec == nil {
		panic("can't get related many of nil")
	}
	rel, err := tx.Schema().GetRelationshipByID(j.on.rel.ID())
	if err != nil || rel == nil {
		panic(j.on.rel.ID())
	}
	hits, err := rel.LoadMany(rec)

	results := []*QueryResult{}
	for _, h := range hits {
		if ok, _ := matcher.Match(h); ok {
			results = append(results, &QueryResult{Record: h})
		}
	}
	return results
}
