package db

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrData         = errors.New("data-error")
	ErrConstraint   = errors.New("constraint-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

type Tx interface {
	AsOfStart() Tx
	Schema() *Schema
	Context() context.Context
	WithContext(context.Context) Tx

	getRelatedOne(ID, ID) (Record, error)
	getRelatedMany(ID, ID) ([]Record, error)
	getRelatedManyReverse(ID, ID) ([]Record, error)
	getRelatedOneReverse(ID, ID) (Record, error)

	MakeRecord(ID) (Record, error)
	Ref(ID) ModelRef
	Query(ModelRef, ...QueryClause) Q
	Subquery(...QueryClause) Q
	Operations() []Operation

	Commit() error
	Abort(error)
}

type RWTx interface {
	Tx
	RWWithContext(context.Context) RWTx
	Insert(Record) error
	Update(oldRec, newRec Record) error
	Delete(Record) error
	Connect(source, target, rel ID) error
	Disconnect(source, target, rel ID) error

	unloggedUpdate(Record, Record) error
	unloggedDelete(Record) error
	unloggedDisconnect(ID, ID, ID) error

	dropImplements(ID, ID) error
	addImplements(ID, ID) error
	dropRel(ID, ID, ID) error
	addRel(ID, ID, ID) error
}

type holdTx struct {
	initH   *hold
	h       *hold
	db      *holdDB
	rw      bool
	ops     []Operation
	aborted error
}

type txWithContext struct {
	*holdTx
	ctx context.Context
}

func (tx *holdTx) Abort(err error) {
	tx.ensureWrite()
	tx.db.writer.Unlock()
	tx.aborted = err
}

func (twc *txWithContext) Context() context.Context {
	return twc.ctx
}

func (tx *txWithContext) WithContext(c context.Context) Tx {
	return &txWithContext{holdTx: tx.holdTx, ctx: c}
}

func (tx *txWithContext) RWWithContext(c context.Context) RWTx {
	return &txWithContext{holdTx: tx.holdTx, ctx: c}
}

func (tx *holdTx) Operations() []Operation {
	return tx.ops
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (tx *txWithContext) AsOfStart() Tx {
	h := &holdTx{initH: tx.initH, h: tx.initH, db: tx.db, rw: false, ops: nil, aborted: nil}
	return &txWithContext{holdTx: h, ctx: tx.ctx}
}

func (tx *holdTx) getRelatedOne(id, rel ID) (Record, error) {
	r, err := tx.h.GetLinkedOne(id, rel)
	return r, err
}

func (tx *holdTx) getRelatedMany(id, rel ID) ([]Record, error) {
	return tx.h.GetLinkedMany(id, rel)
}

func (tx *holdTx) getRelatedManyReverse(id, rel ID) ([]Record, error) {
	return tx.h.GetLinkedManyReverse(id, rel)
}

func (tx *holdTx) getRelatedOneReverse(id, rel ID) (Record, error) {
	return tx.h.GetLinkedOneReverse(id, rel)
}

func (tx *holdTx) Insert(rec Record) error {
	tx.ensureWrite()
	tx.h = tx.h.Insert(rec)
	co := &CreateOp{Record: rec}
	tx.ops = append(tx.ops, co)
	return nil
}

func (tx *holdTx) Update(oldRec, newRec Record) error {
	err := tx.unloggedUpdate(oldRec, newRec)
	if err != nil {
		return err
	}
	uo := &UpdateOp{OldRecord: oldRec, NewRecord: newRec}
	tx.ops = append(tx.ops, uo)
	return nil
}

func (tx *holdTx) unloggedUpdate(oldRec, newRec Record) error {
	tx.ensureWrite()
	if oldRec.ID() != newRec.ID() {
		return fmt.Errorf("Can't update ID field on a record")
	}
	tx.h = tx.h.Insert(newRec)
	return nil
}

func (tx *txWithContext) Connect(source, target, relID ID) error {
	tx.ensureWrite()
	rel, err := tx.Schema().GetRelationshipByID(relID)
	if err != nil {
		return err
	}
	if !rel.Multi() {
		v, err := tx.getRelatedOne(source, relID)
		if err == nil {
			return fmt.Errorf("%w: can't connect already-connected (%v) non-multi relationship %v", ErrConstraint, v.ID(), rel.Name())
		}
		if !errors.Is(ErrNotFound, err) {
			return err
		}
	}
	tx.h = tx.h.Link(source, target, relID)

	co := &ConnectOp{Source: source, Target: target, RelID: relID}
	tx.ops = append(tx.ops, co)
	return nil
}

func (tx *holdTx) Disconnect(source, target, relID ID) error {
	err := tx.unloggedDisconnect(source, target, relID)
	if err != nil {
		return err
	}
	do := &DisconnectOp{Source: source, Target: target, RelID: relID}
	tx.ops = append(tx.ops, do)
	return nil
}

func (tx *holdTx) unloggedDisconnect(source, target, relID ID) error {
	tx.ensureWrite()
	tx.h = tx.h.Unlink(source, target, relID)
	return nil
}

func (tx *holdTx) unloggedDelete(rec Record) error {
	tx.ensureWrite()

	// cascading to rels is handled by the Hold
	tx.h = tx.h.Delete(rec)
	return nil
}

func (tx *holdTx) Delete(rec Record) error {
	err := tx.unloggedDelete(rec)
	if err != nil {
		return err
	}
	do := &DeleteOp{Record: rec}
	tx.ops = append(tx.ops, do)
	return nil
}

func (tx *txWithContext) MakeRecord(interfaceID ID) (rec Record, err error) {
	i, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		return
	}
	rec, err = tx.db.b.RecordForInterface(tx, i)
	return
}

func (tx *txWithContext) Schema() *Schema {
	return &Schema{tx, tx.db}
}

func (tx *txWithContext) Commit() error {
	tx.ensureWrite()
	if tx.aborted != nil {
		return tx.aborted
	}

	tx.db.bus.Publish(BeforeCommit{tx})
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
	tx.db.writer.Unlock()

	return nil
}

func (tx *holdTx) String() string {
	return tx.h.String()
}

func (tx *holdTx) dropImplements(modelID, interfaceID ID) error {
	tx.h = tx.h.dropImplements(modelID, interfaceID)
	return nil
}

func (tx *holdTx) addImplements(modelID, interfaceID ID) error {
	tx.h = tx.h.addImplements(modelID, interfaceID)
	return nil
}

func (tx *holdTx) dropRel(sourceInterfaceID, targetInterfaceID, relID ID) error {
	h, err := tx.h.dropRel(sourceInterfaceID, targetInterfaceID, relID)
	if err != nil {
		return err
	}
	tx.h = h
	return nil

}

func (tx *holdTx) addRel(sourceInterfaceID, targetInterfaceID, relID ID) error {
	tx.h = tx.h.addRel(sourceInterfaceID, targetInterfaceID, relID)
	return nil
}
