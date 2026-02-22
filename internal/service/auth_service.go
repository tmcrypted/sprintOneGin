package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo        UserRepository
	refreshRepo     RefreshSessionRepository
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(userRepo UserRepository, refreshRepo RefreshSessionRepository, jwtSecret string, accessTokenTTL, refreshTokenTTL time.Duration) *authService {
	return &authService{
		userRepo:        userRepo,
		refreshRepo:     refreshRepo,
		jwtSecret:       []byte(jwtSecret),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *authService) Login(ctx context.Context, body dto.LoginRequest) (*dto.AuthResponse, error) {
	if body.Email == "" || body.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, body.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createRefreshSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *authService) Register(ctx context.Context, body dto.CreateUserRequest) (*dto.AuthResponse, error) {
	if body.Email == "" {
		return nil, errors.New("email is required")
	}
	if body.Password == "" {
		return nil, errors.New("password is required")
	}
	if body.FIO == "" {
		return nil, errors.New("fio is required")
	}
	if body.Role == "" {
		body.Role = model.RoleWorker
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        body.Email,
		PasswordHash: string(hash),
		Role:         body.Role,
		FIO:          body.FIO,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.createRefreshSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	tokenHash := hashToken(refreshToken)

	session, err := s.refreshRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.refreshRepo.DeleteByID(ctx, session.ID)
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshRepo.DeleteByID(ctx, session.ID); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := s.createRefreshSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *authService) generateAccessToken(user *model.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  now.Add(s.accessTokenTTL).Unix(),
		"iat":  now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) createRefreshSession(ctx context.Context, userID int64) (string, error) {
	// генерируем случайный токен
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := hex.EncodeToString(raw)
	tokenHash := hashToken(token)

	session := &model.RefreshSession{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}

	if err := s.refreshRepo.Create(ctx, session); err != nil {
		return "", err
	}

	return token, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
