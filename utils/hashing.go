package utils

import "golang.org/x/crypto/bcrypt"

var HashPass = func(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

var MatchPass = func(hashedPass, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
