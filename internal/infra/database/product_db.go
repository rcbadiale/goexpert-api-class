package database

import (
	"goexpert-api/internal/entity"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (p *ProductService) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *ProductService) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductService) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *ProductService) Delete(id string) error {
	product, err := p.FindByID(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}

func (p *ProductService) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error
	if sort != "" || (sort != "asc" && sort != "desc") {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		// Busca com paginação
		err = p.DB.
			Limit(limit).
			Offset((page - 1) * limit).
			Order("created_at " + sort).
			Find(&products).
			Error
	} else {
		// Busca normal
		err = p.DB.
			Order("created_at " + sort).
			Find(&products).
			Error
	}
	return products, err
}
