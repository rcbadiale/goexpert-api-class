package database

import (
	"fmt"
	"goexpert-api/internal/entity"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestCase(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	return db, func() {
		// teardown
	}
}

func TestCreateProduct(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Create(product)
	assert.Nil(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUserFindByIDWhenValidID(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Create(product)
	assert.Nil(t, err)

	productFound, err := productService.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUserFindByIDWhenInvalidID(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Create(product)
	assert.Nil(t, err)

	productFound, err := productService.FindByID("abc123")
	assert.Equal(t, "record not found", err.Error())
	assert.Nil(t, productFound)
}

func TestUserUpdateWhenProductExists(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Create(product)
	assert.Nil(t, err)

	product.Name = "Updated product 1"
	err = productService.Update(product)
	assert.Nil(t, err)

	productFound, err := productService.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUserUpdateWhenProductDoesntExists(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Update(product)
	assert.Equal(t, "record not found", err.Error())
}

func TestUserDeleteWhenProductExists(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Create(product)
	assert.Nil(t, err)

	err = productService.Delete(product.ID.String())
	assert.Nil(t, err)

	productFound, err := productService.FindByID(product.ID.String())
	assert.Equal(t, "record not found", err.Error())
	assert.Nil(t, productFound)
}

func TestUserDeleteWhenProductDoesntExists(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	product, err := entity.NewProduct("Product 1", 10)
	productService := NewProductService(db)

	err = productService.Delete(product.ID.String())
	assert.Equal(t, "record not found", err.Error())
}

func TestProductsFindAll(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	productService := NewProductService(db)

	var products []entity.Product
	for i := range 24 {
		name := fmt.Sprintf("Product %d", i+1)
		product, _ := entity.NewProduct(name, rand.Float64())
		productService.DB.Create(product)
		products = append(products, *product)
	}

	productsFound, err := productService.FindAll(0, 0, "")
	assert.Nil(t, err)
	assert.Len(t, productsFound, 24)
	for i := range 24 {
		assert.Equal(t, products[i].ID, productsFound[i].ID)
		assert.Equal(t, products[i].Name, productsFound[i].Name)
		assert.Equal(t, products[i].Price, productsFound[i].Price)
	}
}

func TestProductsFindAllWithPagination(t *testing.T) {
	db, teardownTest := setupTestCase(t)
	defer teardownTest()

	productService := NewProductService(db)

	var products []entity.Product
	items := 24
	for i := range items {
		name := fmt.Sprintf("Product %d", i+1)
		product, _ := entity.NewProduct(name, rand.Float64())
		productService.DB.Create(product)
		products = append(products, *product)
	}

	limit := 10
	pages := int(math.Ceil(float64(items) / float64(limit)))
	for page := range pages {
		productsFound, err := productService.FindAll(page+1, limit, "asc")
		assert.Nil(t, err)
		assert.LessOrEqual(t, len(productsFound), limit)
		for item := range len(productsFound) {
			assert.Equal(
				t,
				products[item+page*limit].ID,
				productsFound[item].ID,
			)
			assert.Equal(
				t,
				products[item+page*limit].Name,
				productsFound[item].Name,
			)
			assert.Equal(
				t,
				products[item+page*limit].Price,
				productsFound[item].Price,
			)
		}
	}
}
