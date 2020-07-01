package db

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-immutable-radix"
)

var (
	ErrNotFound = fmt.Errorf("%w: not found", ErrData)
)

var sep = []byte("/")

type Hold struct {
	t *iradix.Tree
}

func NewHold() *Hold {
	return &Hold{t: iradix.New()}
}

func (h *Hold) FindOne(modelID ID, q Matcher) (Record, error) {
	mb, _ := modelID.Bytes()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)

	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return nil, err
		}
		if match {
			return rec, nil
		}
	}
	return nil, ErrNotFound
}

type MatchIter struct {
	q  Matcher
	it *iradix.Iterator
}

func (mi MatchIter) Next() (Record, bool) {
	for _, val, ok := mi.it.Next(); ok; _, val, ok = mi.it.Next() {
		rec := val.(Record)
		match, err := mi.q.Match(rec)
		if err != nil {
			return nil, false
		}
		if match {
			return rec, true
		}
	}
	return nil, false
}

func (h *Hold) IterMatches(modelID ID, q Matcher) Iterator {
	mb, _ := modelID.Bytes()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)
	return MatchIter{q: q, it: it}
}

func (h *Hold) FindMany(modelID ID, q Matcher) ([]Record, error) {
	mb, _ := modelID.Bytes()
	it := h.t.Root().Iterator()
	it.SeekPrefix(mb)
	hits := []Record{}
	for _, val, ok := it.Next(); ok; _, val, ok = it.Next() {
		rec := val.(Record)
		match, err := q.Match(rec)
		if err != nil {
			return hits, err
		}
		if match {
			hits = append(hits, rec)
		}
	}
	return hits, nil
}

func makeKey(rec Record) []byte {
	rb, _ := rec.ID().Bytes()
	mb, _ := rec.Model().ID().Bytes()

	bytes := append(append(mb, sep...), rb...)
	return bytes
}

func (h *Hold) Insert(rec Record) *Hold {
	newTree, _, _ := h.t.Insert(makeKey(rec), rec)
	return &Hold{t: newTree}
}

func (h *Hold) Delete(rec Record) *Hold {
	newTree, _, _ := h.t.Delete(makeKey(rec))
	return &Hold{t: newTree}
}

// link/<relid>/<sourceid>/<targetid>
func linkKey(source, target, rel ID) []byte {
	sb, _ := source.Bytes()
	tb, _ := target.Bytes()
	rb, _ := rel.Bytes()

	link := append([]byte("link/"), rb...)
	link = append(append(link, sep...), sb...)
	link = append(append(link, sep...), tb...)
	return link
}

func linkKeyPrefix(id ID, rel Relationship) []byte {
	sb, _ := id.Bytes()
	rb, _ := rel.ID().Bytes()

	link := append([]byte("link/"), rb...)
	link = append(append(link, sep...), sb...)
	return link
}

// rlink/<relid>/<targetid>/<sourceid>
func rlinkKey(source, target, rel ID) []byte {
	sb, _ := source.Bytes()
	tb, _ := target.Bytes()
	rb, _ := rel.Bytes()

	rlink := append([]byte("rlink/"), rb...)
	rlink = append(append(rlink, sep...), tb...)
	rlink = append(append(rlink, sep...), sb...)
	return rlink
}

func rlinkKeyPrefix(id ID, rel Relationship) []byte {
	tb, _ := id.Bytes()
	rb, _ := rel.ID().Bytes()

	rlink := append([]byte("rlink/"), rb...)
	rlink = append(append(rlink, sep...), tb...)
	return rlink
}

func (h *Hold) Link(source, target, rel ID) *Hold {
	lk := linkKey(source, target, rel)
	rk := rlinkKey(source, target, rel)
	newTree, _, _ := h.t.Insert(lk, nil)
	newTree, _, _ = newTree.Insert(rk, nil)

	return &Hold{t: newTree}
}

func (h *Hold) Unlink(source, target, rel ID) *Hold {
	lk := linkKey(source, target, rel)
	rk := rlinkKey(source, target, rel)

	newTree, _, _ := h.t.Delete(lk)
	newTree, _, _ = newTree.Delete(rk)

	return &Hold{t: newTree}
}

func linkKeyComp(k []byte, ix int) []byte {
	switch ix {
	case 0:
		return k[5 : 5+16]
	case 1:
		return k[5+16+1 : 5+16*2+1]
	case 2:
		return k[5+16*2+2 : 5+16*3+2]
	default:
		panic("invalid component")
	}
}

func rlinkKeyComp(k []byte, ix int) []byte {
	switch ix {
	case 0:
		return k[6 : 5+16]
	case 1:
		return k[6+16+1 : 6+16*2+1]
	case 2:
		return k[6+16*2+2 : 6+16*3+2]
	default:
		panic("invalid component")
	}
}

func (h *Hold) followLinks(id ID, rel Relationship, reverse bool) ([]Record, error) {
	var prefix []byte
	if reverse {
		prefix = rlinkKeyPrefix(id, rel)
	} else {
		prefix = linkKeyPrefix(id, rel)
	}

	it := h.t.Root().Iterator()
	it.SeekPrefix(prefix)
	var ids [][]byte
	for k, _, ok := it.Next(); ok; k, _, ok = it.Next() {
		if !bytes.HasPrefix(k, prefix) {
			break
		}
		var targetID []byte
		if reverse {
			targetID = rlinkKeyComp(k, 2)
		} else {
			targetID = linkKeyComp(k, 2)
		}
		ids = append(ids, targetID)
	}

	var hits []Record
	for _, idbytes := range ids {
		id := MakeIDFromBytes(idbytes)
		var hit Record
		var err error
		if reverse {
			hit, err = h.FindOne(rel.Source().ID(), EqID(id))
		} else {
			hit, err = h.FindOne(rel.Target().ID(), EqID(id))
		}
		if err != nil {
			return nil, err
		}
		hits = append(hits, hit)
	}
	return hits, nil
}

func (h *Hold) GetLinkedMany(sID ID, rel Relationship) ([]Record, error) {
	return h.followLinks(sID, rel, false)
}

func (h *Hold) GetLinkedManyReverse(tID ID, rel Relationship) ([]Record, error) {
	return h.followLinks(tID, rel, true)
}

func (h *Hold) followLinksOne(id ID, rel Relationship, reverse bool) (Record, error) {
	hits, err := h.followLinks(id, rel, reverse)
	if err != nil {
		return nil, err
	}
	switch len(hits) {
	case 0:
		return nil, ErrNotFound
	case 1:
		return hits[0], nil
	default:
		panic("Multi found for to-one rel")
	}
}

func (h *Hold) GetLinkedOne(id ID, rel Relationship) (Record, error) {
	return h.followLinksOne(id, rel, false)
}

func (h *Hold) GetLinkedOneReverse(id ID, rel Relationship) (Record, error) {
	return h.followLinksOne(id, rel, true)
}

type RootIter struct {
	it *iradix.Iterator
}

func (ri RootIter) Next() (Record, bool) {
	for _, val, ok := ri.it.Next(); ok; _, val, ok = ri.it.Next() {
		rec, ok := val.(Record)
		if !ok {
			// skip links, i think
			continue
		}
		return rec, true
	}
	return nil, false
}

func (h *Hold) Iterator() Iterator {
	it := h.t.Root().Iterator()
	return RootIter{it: it}
}

func (h *Hold) PrintTree() {
	fmt.Printf("print tree:\n")
	it := h.t.Root().Iterator()
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		fmt.Printf("%v:%v\n", string(k), v)
	}
	fmt.Printf("done printing\n")
}
