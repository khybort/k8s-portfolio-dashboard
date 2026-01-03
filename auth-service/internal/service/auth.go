package service

import (
	"context"
	"errors"
	"time"
	"github.com/portfolio/auth-service/internal/jwt"
	"github.com/portfolio/auth-service/internal/model"
	"github.com/portfolio/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, string, *model.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	VerifyToken(ctx context.Context, token string) (*jwt.Claims, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtSecret string
	accessExpiry  int
	refreshExpiry int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, accessExpiry, refreshExpiry int) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*model.User, error) {
	// Check if user exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
		Role:         "user",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Role, s.jwtSecret, time.Duration(s.accessExpiry)*time.Minute)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String(), s.jwtSecret, time.Duration(s.refreshExpiry)*time.Hour)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := jwt.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return "", errors.New("invalid token type")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String(), user.Role, s.jwtSecret, time.Duration(s.accessExpiry)*time.Minute)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *authService) VerifyToken(ctx context.Context, token string) (*jwt.Claims, error) {
	return jwt.ValidateToken(token, s.jwtSecret)
}

