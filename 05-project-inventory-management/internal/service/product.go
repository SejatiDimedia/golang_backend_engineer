package service

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductService interface {
	Create(ctx context.Context, prod *entity.Product) (*entity.Product, error)
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetAll(ctx context.Context, search string, categoryID uint, page, limit int) ([]entity.Product, int64, error)
	Update(ctx context.Context, id uint, prod *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
	supplierRepo repository.SupplierRepository
}

func NewProductService(
	repo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	supplierRepo repository.SupplierRepository,
) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
		supplierRepo: supplierRepo,
	}
}

func (s *productService) Create(ctx context.Context, prod *entity.Product) (*entity.Product, error) {
	// 1. Validasi CategoryID
	cat, err := s.categoryRepo.GetByID(ctx, prod.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("invalid category_id: category does not exist")
	}

	// 2. Validasi SupplierID
	sup, err := s.supplierRepo.GetByID(ctx, prod.SupplierID)
	if err != nil {
		return nil, err
	}
	if sup == nil {
		return nil, errors.New("invalid supplier_id: supplier does not exist")
	}

	prod.StockQuantity = 0 // Inisiasi awal stok selalu 0 (stok hanya berubah lewat mutasi)

	if err := s.repo.Create(ctx, prod); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, prod.ID)
}

func (s *productService) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	prod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if prod == nil {
		return nil, ErrProductNotFound
	}
	return prod, nil
}

func (s *productService) GetAll(ctx context.Context, search string, categoryID uint, page, limit int) ([]entity.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetAll(ctx, search, categoryID, page, limit)
}

func (s *productService) Update(ctx context.Context, id uint, prod *entity.Product) (*entity.Product, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}

	// Validasi CategoryID jika dirubah
	if prod.CategoryID != existing.CategoryID {
		cat, err := s.categoryRepo.GetByID(ctx, prod.CategoryID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, errors.New("invalid category_id: category does not exist")
		}
		existing.CategoryID = prod.CategoryID
	}

	// Validasi SupplierID jika dirubah
	if prod.SupplierID != existing.SupplierID {
		sup, err := s.supplierRepo.GetByID(ctx, prod.SupplierID)
		if err != nil {
			return nil, err
		}
		if sup == nil {
			return nil, errors.New("invalid supplier_id: supplier does not exist")
		}
		existing.SupplierID = prod.SupplierID
	}

	existing.Name = prod.Name
	existing.SKU = prod.SKU
	existing.Description = prod.Description
	existing.Price = prod.Price

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	prod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if prod == nil {
		return ErrProductNotFound
	}
	return s.repo.Delete(ctx, id)
}
