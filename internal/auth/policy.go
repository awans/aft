package auth

import (
	"encoding/json"
	"fmt"

	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/api/parsers"
	"awans.org/aft/internal/db"
)

var PolicyModel = db.MakeModel(
	db.MakeID("ea5eda03-6780-4a31-8b9b-e5f16a98d8b3"),
	"policy",
	[]db.AttributeL{
		pText,
		pRead,
		pWrite,
	},
	// set in init
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

var pText = db.MakeConcreteAttribute(
	db.MakeID("55cfda72-c7f2-47aa-85ab-e54b98f1eda0"),
	"text",
	db.String,
)

var pRead = db.MakeConcreteAttribute(
	db.MakeID("c14bf00b-2d76-4e3b-8a54-f20d4064784f"),
	"read",
	db.Bool,
)

var pWrite = db.MakeConcreteAttribute(
	db.MakeID("4a5b3ccb-6d30-4bbd-91e0-524e6d8ce445"),
	"write",
	db.Bool,
)

var PolicyFor = db.MakeConcreteRelationship(
	db.MakeID("be24d5ca-48f4-4d6f-a550-5b969703f440"),
	"interface",
	false,
	db.InterfaceInterface,
)

var InterfacePolicies = db.MakeReverseRelationship(
	db.MakeID("09579552-6982-4732-9d69-585f2e6a74b1"),
	"policies",
	PolicyFor,
)

var PolicyRole = db.MakeReverseRelationship(
	db.MakeID("e7bb2583-ce26-4369-86dc-9a8f6952ad2e"),
	"role",
	RolePolicy,
)

type PolicyL struct {
	ID_   db.ID  `record:"id"`
	Text_ string `record:"text"`
	For_  db.ModelL
}

func (lit PolicyL) ID() db.ID {
	return lit.ID_
}

func (lit PolicyL) MarshalDB() (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(lit, PolicyModel)
	recs = append(recs, rec)
	links = append(links, db.Link{rec.ID(), lit.For_.ID(), PolicyFor})
	return
}

type policy struct {
	rec db.Record
	tx  db.Tx
}

func (p *policy) String() string {
	return fmt.Sprintf("policy{\"%v\"}", p.Text())
}

func (p *policy) Text() string {
	return pText.MustGet(p.rec).(string)
}

func (p *policy) Interface() db.Interface {
	tx := p.tx
	policies := tx.Ref(PolicyModel.ID())
	ifaces := tx.Ref(db.InterfaceInterface.ID())
	ifrec, err := tx.Query(ifaces, db.Join(policies, ifaces.Rel(InterfacePolicies)), db.Filter(policies, db.EqID(p.rec.ID()))).OneRecord()
	if err != nil {
		panic("No model")
	}
	// this is awkward and inefficient
	i, err := tx.Schema().GetInterfaceByID(ifrec.ID())
	if err != nil {
		panic("No model")
	}
	return i
}

func (p *policy) Apply(tx db.Tx, ref db.ModelRef) []db.QueryClause {
	iface, err := tx.Schema().GetInterfaceByID(ref.InterfaceID)
	if err != nil {
		panic("bad")
	}
	var data map[string]interface{}

	json.Unmarshal([]byte(p.Text()), &data)
	w, err := parsers.Parser{tx}.ParseWhere(iface, data)
	if err != nil {
		panic("bad")
	}
	clauses := operations.HandleWhere(tx, ref, w)
	return clauses
}
