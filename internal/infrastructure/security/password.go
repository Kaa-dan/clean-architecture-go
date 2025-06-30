package security

import "golang.org/x/crypto/bcrypt"

type PasswordManager struct {
	cost int
}

func NewPasswordManger() *PasswordManager {
	return &PasswordManager{
		cost: bcrypt.DefaultCost,
	}
}

func (p *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (p *PasswordManager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
