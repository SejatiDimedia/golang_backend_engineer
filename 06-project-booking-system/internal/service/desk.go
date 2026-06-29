package service

import (
	"context"
	"errors"

	"github.com/timurdian/booking-system/internal/entity"
	"github.com/timurdian/booking-system/internal/repository"
)

var (
	ErrDeskNotFound = errors.New("desk not found")
)

type DeskService interface {
	Create(ctx context.Context, name, deskType string) (*entity.Desk, error)
	GetByID(ctx context.Context, id uint) (*entity.Desk, error)
	GetAllActive(ctx context.Context) ([]entity.Desk, error)
	GetAll(ctx context.Context) ([]entity.Desk, error)
	Update(ctx context.Context, id uint, name, deskType string, isActive bool) (*entity.Desk, error)
	Delete(ctx context.Context, id uint) error
}

type deskService struct {
	repo repository.DeskRepository
}

func NewDeskService(repo repository.DeskRepository) DeskService {
	return &deskService{repo: repo}
}

func (s *deskService) Create(ctx context.Context, name, deskType string) (*entity.Desk, error) {
	if name == "" {
		return nil, errors.New("desk name cannot be empty")
	}
	if deskType != "hot-desk" && deskType != "meeting-room" {
		return nil, errors.New("invalid desk type")
	}

	desk := &entity.Desk{
		Name:     name,
		Type:     deskType,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, desk); err != nil {
		return nil, err
	}
	return desk, nil
}

func (s *deskService) GetByID(ctx context.Context, id uint) (*entity.Desk, error) {
	desk, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if desk == nil {
		return nil, ErrDeskNotFound
	}
	return desk, nil
}

func (s *deskService) GetAllActive(ctx context.Context) ([]entity.Desk, error) {
	return s.repo.GetAllActive(ctx)
}

func (s *deskService) GetAll(ctx context.Context) ([]entity.Desk, error) {
	return s.repo.GetAll(ctx)
}

func (s *deskService) Update(ctx context.Context, id uint, name, deskType string, isActive bool) (*entity.Desk, error) {
	desk, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if desk == nil {
		return nil, ErrDeskNotFound
	}

	if name == "" {
		return nil, errors.New("desk name cannot be empty")
	}
	if deskType != "hot-desk" && deskType != "meeting-room" {
		return nil, errors.New("invalid desk type")
	}

	desk.Name = name
	desk.Type = deskType
	desk.IsActive = isActive

	if err := s.repo.Update(ctx, desk); err != nil {
		return nil, err
	}
	return desk, nil
}

func (s *deskService) Delete(ctx context.Context, id uint) error {
	desk, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if desk == nil {
		return ErrDeskNotFound
	}
	return s.repo.Delete(ctx, id)
}
