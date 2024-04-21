package utils

import "golang.org/x/crypto/bcrypt"

func HashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ComparePassword(password []byte, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
