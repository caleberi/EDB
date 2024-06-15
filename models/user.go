package models

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	omitList = []string{
		"password",
		"createdAt",
		"updatedAt",
	}
)

type User struct {
	ID                 primitive.ObjectID `validate:"mongodb" json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName          string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName           string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email              string             `validate:"required,email" json:"email,omitempty" bson:"email,omitempty"`
	Password           string             `validate:"required" json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt          time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty" validate:"required"`
	UpdatedAt          time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" validate:"required"`
	MiddleName         string             `json:"middleName,omitempty" bson:"middleName,omitempty"`
	BVN                string             `json:"bvn,omitempty" bson:"bvn,omitempty"`
	DOB                string             `json:"dob,omitempty" bson:"dob,omitempty"`
	Address            string             `json:"address,omitempty" bson:"address,omitempty"`
	Phone              string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Country            string             `json:"country,omitempty" bson:"country,omitempty"`
	IdNumber           string             `json:"idNumber,omitempty" bson:"idNumber,omitempty"`
	IdType             string             `json:"idType,omitempty" bson:"idType,omitempty"`
	AdditionalIdType   string             `json:"additionalIdType,omitempty" bson:"additionalIdType,omitempty"`
	AdditionalIdNumber string             `json:"additionalIdNumber,omitempty" bson:"additionalIdNumber,omitempty"`
}

func (u *User) Omit() (User, error) {
	copiedUser := *u

	userType := reflect.TypeOf(copiedUser)
	fieldValues := make(map[string]interface{})

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		if !lo.Contains(omitList, field.Name) {
			value := reflect.ValueOf(*u).FieldByName(field.Name).Interface()
			fieldValues[field.Name] = value
		}
	}

	filteredJSON, err := json.Marshal(fieldValues)
	if err != nil {
		return User{}, err
	}

	var filteredUser User
	err = json.Unmarshal(filteredJSON, &filteredUser)
	if err != nil {
		return User{}, err
	}

	filteredUser.ID = u.ID
	return filteredUser, nil
}
