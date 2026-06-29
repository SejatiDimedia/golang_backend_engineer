package service

import (
	"context"
	"errors"

	"github.com/timurdian/file-management/internal/entity"
	"github.com/timurdian/file-management/internal/repository"
	"github.com/timurdian/file-management/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists       = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry int) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(hashedPass),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
}
