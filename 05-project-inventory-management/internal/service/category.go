package service

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/repository"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService interface {
	Create(ctx context.Context, name string) (*entity.Category, error)
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, id uint, name string) (*entity.Category, error)
	Delete(ctx context.Context, id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, name string) (*entity.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	cat := &entity.Category{Name: name}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}
	return cat, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]entity.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) Update(ctx context.Context, id uint, name string) (*entity.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	cat.Name = name
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrCategoryNotFound
	}
	return s.repo.Delete(ctx, id)
}
