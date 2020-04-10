package db

import (
	"awans.org/aft/internal/model"
)

func (db DB) AddSampleModels() {
	db.db.Insert(User.Name, &User)
	db.db.Insert(Profile.Name, &Profile)
	db.db.Insert(Post.Name, &Post)
}

var User = model.Model{
	Name: "user",
	Attributes: map[string]model.Attribute{
		"id": model.Attribute{
			Type: model.UUID,
		},
		"firstName": model.Attribute{
			Type: model.String,
		},
		"lastName": model.Attribute{
			Type: model.String,
		},
		"age": model.Attribute{
			Type: model.Int,
		},
	},
	Relationships: map[string]model.Relationship{
		"posts": model.Relationship{
			Target: "Post",
			Type:   model.HasMany,
		},
		"profile": model.Relationship{
			Target: "Profile",
			Type:   model.HasOne,
		},
	},
}

var Profile = model.Model{
	Name: "profile",
	Attributes: map[string]model.Attribute{
		"id": model.Attribute{
			Type: model.UUID,
		},
		"text": model.Attribute{
			Type: model.String,
		},
	},
	Relationships: map[string]model.Relationship{
		"user": model.Relationship{
			Target: "User",
			Type:   model.BelongsTo,
		},
	},
}

var Post = model.Model{
	Name: "post",
	Attributes: map[string]model.Attribute{
		"id": model.Attribute{
			Type: model.UUID,
		},
		"text": model.Attribute{
			Type: model.String,
		},
	},
	Relationships: map[string]model.Relationship{
		"author": model.Relationship{
			Target: "User",
			Type:   model.BelongsTo,
		},
	},
}
