package auth

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var User = db.Model{
	Id:   uuid.MustParse("e52f8264-7b95-4a3a-bf76-a23b2229d65a"),
	Name: "user",
	Attributes: map[string]db.Attribute{
		"email": db.Attribute{
			Id:       uuid.MustParse("236e800d-c39d-4ef3-94e6-5e1f0fc38e62"),
			AttrType: db.String,
		},
		"password": db.Attribute{
			Id:       uuid.MustParse("658f314a-4602-44a9-8d19-884bbd3ea267"),
			AttrType: db.String,
		},
	},
}
