package oplog

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

type DBOp int

const (
	Connect DBOp = iota
	Create
	Update
	Delete
)

type TxEntry struct {
	Ops []DBOpEntry
}

func (txe TxEntry) Replay(rwtx db.RWTx) {
	for _, op := range txe.Ops {
		op.Replay(rwtx)
	}
}

type DBOpEntry struct {
	OpType DBOp
	Op     interface{}
}

func (oe DBOpEntry) Replay(rwtx db.RWTx) {
	switch oe.OpType {
	case Create:
		cro := oe.Op.(CreateOp)
		cro.Replay(rwtx)
	case Connect:
		cno := oe.Op.(ConnectOp)
		cno.Replay(rwtx)
	case Update:
		uo := oe.Op.(UpdateOp)
		uo.Replay(rwtx)
	case Delete:
		panic("Not implemented")
	}
}

type CreateOp struct {
	RecordFields interface{}
	ModelID      uuid.UUID
}

func (cro CreateOp) Replay(rwtx db.RWTx) {
	st := cro.RecordFields
	m, err := rwtx.GetModelByID(cro.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Insert(db.RecordFromParts(st, m))
}

type ConnectOp struct {
	Left  uuid.UUID
	Right uuid.UUID
	RelID uuid.UUID
}

func (cno ConnectOp) Replay(rwtx db.RWTx) {
	relRec, err := rwtx.FindOne(db.RelationshipModel.ID, db.Eq("id", cno.RelID))
	if err != nil {
		panic("couldn't find one on replay")
	}
	rel, err := db.LoadRel(relRec)
	if err != nil {
		panic("couldn't find one on replay")
	}
	left, err := rwtx.FindOne(rel.LeftModelID, db.Eq("id", cno.Left))
	if err != nil {
		panic("couldn't find one on replay")
	}
	right, err := rwtx.FindOne(rel.RightModelID, db.Eq("id", cno.Right))
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Connect(left, right, rel)
}

type UpdateOp struct {
	OldFields interface{}
	NewFields interface{}
	ModelID   uuid.UUID
}

func (uo UpdateOp) Replay(rwtx db.RWTx) {
	Ost := uo.OldFields
	Nst := uo.NewFields
	m, err := rwtx.GetModelByID(uo.ModelID)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Update(db.RecordFromParts(Ost, m), db.RecordFromParts(Nst, m))
}

type DeleteOp struct {
	ID uuid.UUID
}

type loggedDB struct {
	inner db.DB
	l     OpLog
}

type loggedTx struct {
	inner db.RWTx
	txe   TxEntry
	l     OpLog
}

func DBFromLog(db db.DB, l OpLog) error {
	iter := l.Iterator()
	rwtx := db.NewRWTx()
	for iter.Next() {
		val := iter.Value()
		txe := val.(TxEntry)
		txe.Replay(rwtx)
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	err := rwtx.Commit()
	return err
}

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{inner: d, l: l}
}

func (l *loggedDB) NewTx() db.Tx {
	return l.inner.NewTx()
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{inner: l.inner.NewRWTx(), l: l.l}
}

func (l *loggedDB) DeepEquals(o db.DB) bool {
	return l.inner.DeepEquals(o)
}

func (l *loggedDB) Iterator() db.Iterator {
	return l.inner.Iterator()
}

func (tx *loggedTx) GetModelByID(id uuid.UUID) (db.Model, error) {
	return tx.inner.GetModelByID(id)
}

func (tx *loggedTx) GetModel(modelName string) (db.Model, error) {
	return tx.inner.GetModel(modelName)
}

func (tx *loggedTx) SaveModel(m db.Model) error {
	return tx.inner.SaveModel(m)
}

func (tx *loggedTx) Ref(u uuid.UUID) db.ModelRef {
	return tx.inner.Ref(u)
}

func (tx *loggedTx) Query(m db.ModelRef) db.Q {
	return tx.inner.Query(m)
}

func (tx *loggedTx) MakeRecord(modelID uuid.UUID) db.Record {
	return tx.inner.MakeRecord(modelID)
}

func (tx *loggedTx) Insert(rec db.Record) error {
	co := CreateOp{RecordFields: rec.RawData(), ModelID: rec.Model().ID}
	dboe := DBOpEntry{Create, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Insert(rec)
}

func (tx *loggedTx) Connect(left, right db.Record, rel db.Relationship) error {
	co := ConnectOp{Left: left.ID(), Right: right.ID(), RelID: rel.ID}
	dboe := DBOpEntry{Connect, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Connect(left, right, rel)
}

func (tx *loggedTx) Update(oldRec, newRec db.Record) error {
	uo := UpdateOp{OldFields: oldRec.RawData(), NewFields: newRec.RawData(), ModelID: oldRec.Model().ID}
	dboe := DBOpEntry{Update, uo}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	return tx.inner.Update(oldRec, newRec)
}

func (tx *loggedTx) FindOne(modelID uuid.UUID, matcher db.Matcher) (db.Record, error) {
	return tx.inner.FindOne(modelID, matcher)
}

func (tx *loggedTx) FindMany(modelID uuid.UUID, matcher db.Matcher) ([]db.Record, error) {
	return tx.inner.FindMany(modelID, matcher)
}

func (tx *loggedTx) Commit() (err error) {
	err = tx.l.Log(tx.txe)
	if err != nil {
		return
	}
	err = tx.inner.Commit()
	return
}
