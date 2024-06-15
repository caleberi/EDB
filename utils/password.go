package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	Cost int
}

func NewHasher(cost int) Hasher {
	return Hasher{
		Cost: cost,
	}
}

func (h Hasher) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", fmt.Errorf("failed to successfully hash password : [%w] ", err)
	}
	return string(hashed), nil
}

func (h Hasher) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
