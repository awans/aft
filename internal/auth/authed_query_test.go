package auth

import (
	"context"
	"testing"

	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var suser = MakeUser(
	db.MakeID("293ee4d6-b846-4d21-b8b7-24faef34bc81"),
	"signedin@gmail.com",
	"coolpass2",
	signedin,
)

var auser = MakeUser(
	db.MakeID("5edcd0d0-fab3-4fc6-b998-ea0eae1fbe88"),
	"admin@gmail.com",
	"coolpass1",
	admin,
)

var tuser = MakeUser(
	db.MakeID("a6c34811-28f1-4cb1-bbea-56c85009011a"),
	"tech@gmail.com",
	"coolpass3",
	tech,
)

var signedin = RoleL{
	ID_:      db.MakeID("9cdf547a-03ae-4c88-aee2-fe5647c3252d"),
	Name:     "signedin",
	Policies: []PolicyL{signedinPolicy},
}

var admin = RoleL{
	ID_:      db.MakeID("b928a9a9-7760-4ece-81ff-1bc74ae51a63"),
	Name:     "admin",
	Policies: []PolicyL{adminPolicy},
}

var tech = RoleL{
	ID_:      db.MakeID("b0c63d48-dd5f-4b2a-809b-8c0e97e9dd05"),
	Name:     "tech",
	Policies: []PolicyL{techPolicy},
}

var signedinPolicy = PolicyL{
	ID_:       db.MakeID("f4884fb0-9fef-4af8-82cc-3592591b035d"),
	AllowRead: true,
	ReadWhere: `{ "email": "signedin@gmail.com" }`,
	For_:      UserModel,
}

var adminPolicy = PolicyL{
	ID_:       db.MakeID("09d881bc-5246-4197-847b-037b55c2e5b0"),
	AllowRead: true,
	ReadWhere: `{}`,
	For_:      UserModel,
}

var techPolicy = PolicyL{
	ID_:       db.MakeID("bc387bc8-90fe-4749-b7d0-3bf74bfd0eac"),
	AllowRead: true,
	ReadWhere: `{"role": {"name": "admin"}}`,
	For_:      UserModel,
}

var testData = []db.Literal{
	signedin, admin, tech,
	tuser, auser, suser,
	signedinPolicy, techPolicy, adminPolicy,
}

func pluckIDs(recs []db.Record) (ids []db.ID) {
	for _, r := range recs {
		ids = append(ids, r.ID())
	}
	return
}

func TestAuthedQuery(t *testing.T) {
	appDB := db.NewTest()
	rwtx := appDB.NewRWTx()
	appDB.AddLiteral(rwtx, Password)
	appDB.RegisterNativeFunction(passwordValidator)
	appDB.AddLiteral(rwtx, PolicyModel)
	appDB.AddLiteral(rwtx, UserModel)
	appDB.AddLiteral(rwtx, RoleModel)
	appDB.RegisterNativeFunction(emailAddressValidator)
	appDB.AddLiteral(rwtx, EmailAddress)
	rwtx.Commit()

	rwtx = appDB.NewRWTx()
	for _, t := range testData {
		appDB.AddLiteral(rwtx, t)
	}
	rwtx.Commit()
	authDB := AuthedDB(appDB)

	var cases = []struct {
		user    UserL
		results []db.ID
	}{
		{
			tuser,
			[]db.ID{auser.ID()},
		},
		{
			auser,
			[]db.ID{auser.ID(), tuser.ID(), suser.ID()},
		},
		{
			suser,
			[]db.ID{suser.ID()},
		},
	}

	opt := cmpopts.SortSlices(func(a, b db.ID) bool {
		return a.String() < b.String()
	})

	for _, c := range cases {
		tx := authDB.NewTx()
		deesc := Escalate(tx)
		users := tx.Ref(UserModel.ID())
		uRec, _ := tx.Query(users, db.Filter(users, db.EqID(c.user.ID()))).OneRecord()
		ctx := withUser(context.Background(), uRec)
		role, err := RoleForUser(tx, uRec)
		deesc()
		if err != nil {
			t.Errorf("error getting role: %v", err)
		}
		ctx = withRole(ctx, role)

		aTx := authDB.NewTxWithContext(ctx)
		users = aTx.Ref(UserModel.ID())
		results := aTx.Query(users).Records()
		ids := pluckIDs(results)

		diff := cmp.Diff(c.results, ids, opt)
		if diff != "" {
			t.Errorf("case: %v\n(-want +got):\n%s", c.user, diff)
		}
	}
}
