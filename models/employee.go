package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	FirstName        string             `bson:"firstName,omitempty" json:"firstName,omitempty" validate:"required"`
	LastName         string             `bson:"lastName,omitempty" json:"lastName,omitempty" validate:"required"`
	MiddleName       string             `bson:"middleName,omitempty" json:"middleName,omitempty"`
	Email            string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"required"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"required"`
	BVN              string             `bson:"bvn,omitempty" json:"bvn,omitempty"`
	DOB              string             `bson:"dob,omitempty" json:"dob,omitempty"`
	Address          string             `bson:"address,omitempty" json:"address,omitempty"`
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Country          string             `bson:"country,omitempty" json:"country,omitempty"`
	IDNumber         string             `bson:"idNumber,omitempty" json:"idNumber,omitempty"`
	IDType           string             `bson:"idType,omitempty" json:"idType,omitempty"`
	AdditionalIDType string             `bson:"additionalIdType,omitempty" json:"additionalIdType,omitempty"`
	Salary           float64            `bson:"salary,omitempty" json:"salary,omitempty" validate:"required"`
	UserID           primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty" validate:"required"`
	AccountName      string             `bson:"account_name,omitempty" json:"account_name,omitempty" validate:"required"`
	AccountType      string             `bson:"account_type,omitempty" json:"account_type,omitempty" validate:"required"`
	BankName         string             `bson:"bank_name,omitempty" json:"bank_name,omitempty" validate:"required"`
}
