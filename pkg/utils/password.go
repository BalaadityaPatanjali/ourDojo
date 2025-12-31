package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword converts plain password → bcrypt hash
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hashedBytes), err
}

// CheckPassword compares plain password with stored hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
