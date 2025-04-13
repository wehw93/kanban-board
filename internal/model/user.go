package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID                 int
	Name               string
	Email              string
	Password           string
	Encrypted_password string
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptstring(u.Password)
		if err != nil {
			return err
		}
		u.Encrypted_password = enc
	}
	return nil
}

func encryptstring(pass string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
