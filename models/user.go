package models

import (
	"github.com/zebresel-com/mongodm"
)

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Name                 string `json:"name" bson:"name" required:"true"`
	Email                string `json:"email" bson:"email" validation:"email" required:"true"`
	PasswordRaw          string `json:"password,omitempty" bson:"-"`
	PasswordHash         string `json:"-" bson:"password"`
}
