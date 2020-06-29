package auth

import (
	"awans.org/aft/internal/bizdatatypes"
	"awans.org/aft/internal/db"
)

var UserModel = db.MakeModel(
	db.MakeID("e52f8264-7b95-4a3a-bf76-a23b2229d65a"),
	"user",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("236e800d-c39d-4ef3-94e6-5e1f0fc38e62"),
			"email",
			bizdatatypes.EmailAddress,
		),
		db.MakeConcreteAttribute(
			db.MakeID("658f314a-4602-44a9-8d19-884bbd3ea267"),
			"password",
			db.String,
		),
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)
