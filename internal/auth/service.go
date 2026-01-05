package auth

import (
	"errors"
	"time"

	"github.com/quynhanh/internship-tracker/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Register
func (s *Service) Register(email, password string) error {
	// Check if email already exists
	var existing model.User
	err := s.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return errors.New("Email already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	user := model.User{
		Email:        email,
		PasswordHash: string(hashed),
	}

	return s.DB.Create(&user).Error
}

// Login
func (s *Service) Login(email, password string) (string, error) {
	// Check email and password
	var user model.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("Invalid email or password")
	}

	// Compare user's logged password with registered password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", errors.New("Invalid email or password")
	}

	// Return a newly created token via jwt.go layer
	return GenerateToken(user.ID)
}

// Mark tokens that no longer in use
func (s *Service) RevokeToken(token string, exp time.Time) error {
	return s.DB.Create(&model.RevokedToken{
		Token:     token,
		ExpiresAt: exp,
	}).Error
}
