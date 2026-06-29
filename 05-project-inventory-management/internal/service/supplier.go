package service

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/repository"
)

var (
	ErrSupplierNotFound = errors.New("supplier not found")
)

type SupplierService interface {
	Create(ctx context.Context, name, contactName, email, phone string) (*entity.Supplier, error)
	GetByID(ctx context.Context, id uint) (*entity.Supplier, error)
	GetAll(ctx context.Context) ([]entity.Supplier, error)
	Update(ctx context.Context, id uint, name, contactName, email, phone string) (*entity.Supplier, error)
	Delete(ctx context.Context, id uint) error
}

type supplierService struct {
	repo repository.SupplierRepository
}

func NewSupplierService(repo repository.SupplierRepository) SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) Create(ctx context.Context, name, contactName, email, phone string) (*entity.Supplier, error) {
	if name == "" {
		return nil, errors.New("supplier name cannot be empty")
	}

	sup := &entity.Supplier{
		Name:        name,
		ContactName: contactName,
		Email:       email,
		Phone:       phone,
	}
	if err := s.repo.Create(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *supplierService) GetByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	sup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sup == nil {
		return nil, ErrSupplierNotFound
	}
	return sup, nil
}

func (s *supplierService) GetAll(ctx context.Context) ([]entity.Supplier, error) {
	return s.repo.GetAll(ctx)
}

func (s *supplierService) Update(ctx context.Context, id uint, name, contactName, email, phone string) (*entity.Supplier, error) {
	if name == "" {
		return nil, errors.New("supplier name cannot be empty")
	}

	sup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sup == nil {
		return nil, ErrSupplierNotFound
	}

	sup.Name = name
	sup.ContactName = contactName
	sup.Email = email
	sup.Phone = phone

	if err := s.repo.Update(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *supplierService) Delete(ctx context.Context, id uint) error {
	sup, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sup == nil {
		return ErrSupplierNotFound
	}
	return s.repo.Delete(ctx, id)
}
