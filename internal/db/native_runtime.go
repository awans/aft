package db

import (
	"github.com/google/uuid"
)

type Func func(interface{}) (interface{}, error)

type NativeRuntime struct {
	fMap map[ID]Func
	db   DB
}

func NewNativeRuntime(db DB) *NativeRuntime {
	return &NativeRuntime{
		fMap: make(map[ID]Func),
		db:   db,
	}
}

func (nr *NativeRuntime) ProvideModel() ModelL {
	return NativeFunctionModel
}

func (nr *NativeRuntime) Load(tx Tx, rec Record) Function {
	return nativeFunction{rec, nr, tx}
}

func MakeNativeFunction(id ID, name string, functionSignature EnumValueL, function Func) NativeFunctionL {
	return NativeFunctionL{
		id, name, functionSignature, function,
	}
}

type NativeFunctionL struct {
	ID_                ID         `record:"id"`
	Name_              string     `record:"name"`
	FunctionSignature_ EnumValueL `record:"functionSignature"`
	Function           Func
}

func (lit NativeFunctionL) ID() ID {
	return lit.ID_
}

func (lit NativeFunctionL) GetID() ID {
	return lit.ID_
}

func (lit NativeFunctionL) Name() string {
	return lit.Name_
}

func (lit NativeFunctionL) FunctionSignature() EnumValue {
	return lit.FunctionSignature_
}

func (lit NativeFunctionL) Call(args interface{}) (interface{}, error) {
	return lit.Function(args)
}

func (lit NativeFunctionL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, NativeFunctionModel)
	recs = append(recs, rec)
	return
}

func (nr *NativeRuntime) Save(lit NativeFunctionL) {
	f := lit.Function
	tx := nr.db.NewRWTx()
	rec := MarshalRecord(lit, NativeFunctionModel)
	tx.Insert(rec)
	nr.fMap[lit.ID()] = f
	tx.Commit()
}

var NativeFunctionModel = ModelL{
	ID:   MakeID("8deaec0c-f281-4583-baf7-89c3b3b051f3"),
	Name: "code",
	Attributes: []AttributeL{
		nfName,
		nfFuncSig,
	},
}
var nfName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("c47bcd30-01ea-467f-ad02-114342070241"),
	Datatype: String,
}

var nfFuncSig = ConcreteAttributeL{
	Name:     "functionSignature",
	ID:       MakeID("ba29d820-ae50-4424-b807-1a1dbd8d2f4b"),
	Datatype: FunctionSignature,
}

type nativeFunction struct {
	rec Record
	nr  *NativeRuntime
	tx  Tx
}

func (nf nativeFunction) ID() ID {
	return nf.rec.ID()
}

func (nf nativeFunction) Name() string {
	return nfName.AsAttribute().MustGet(nf.rec).(string)
}

func (nf nativeFunction) FunctionSignature() EnumValue {
	u := nfFuncSig.AsAttribute().MustGet(nf.rec).(uuid.UUID)
	ev, _ := nf.tx.Schema().GetEnumValueByID(ID(u))
	return ev
}

func (nf nativeFunction) Call(args interface{}) (interface{}, error) {
	panic("Not Implemented")
}
