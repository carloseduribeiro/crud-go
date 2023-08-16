package database

import (
	"fmt"
	"github.com/carloseduribeiro/crud-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"testing"
)

func prepareDBConnForProduct(t *testing.T) (db *gorm.DB, err error) {
	t.Helper()
	if db, err = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{}); err != nil {
		return
	}
	if err = db.AutoMigrate(&entity.Product{}); err != nil {
		return nil, err
	}
	return
}

func TestCreateNewProduct(t *testing.T) {
	db, err := prepareDBConnForProduct(t)
	if err != nil {
		t.Error(err)
	}
	product, _ := entity.NewProduct("Product 1", 10.00)
	productDB := NewProduct(db)
	err = productDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db, err := prepareDBConnForProduct(t)
	if err != nil {
		t.Error(err)
	}
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		err = db.Create(product).Error
		assert.NoError(t, err)
	}
	productDB := NewProduct(db)
	// pg 1
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)
	// pg 2
	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)
	// pg 3
	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)
}

func TestFindProductById(t *testing.T) {
	db, err := prepareDBConnForProduct(t)
	if err != nil {
		t.Error(err)
	}
	product, _ := entity.NewProduct("Product 1", 10.00)
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}
	productDB := NewProduct(db)
	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, product.Name, productFound.Name)
}

func TestUpdateProduct(t *testing.T) {
	db, err := prepareDBConnForProduct(t)
	if err != nil {
		t.Error(err)
	}
	product, _ := entity.NewProduct("Product 1", 10.00)
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}
	product.Name = "Product 2"
	productDB := NewProduct(db)
	err = productDB.Update(product)
	assert.NoError(t, err)
	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 2", productFound.Name)
}

func TestDeleteProduct(t *testing.T) {
	db, err := prepareDBConnForProduct(t)
	if err != nil {
		t.Error(err)
	}
	product, _ := entity.NewProduct("Product 1", 10.00)
	err = db.Create(product).Error
	if err != nil {
		t.Error(err)
	}
	productDB := NewProduct(db)
	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)
	productFound, err := productDB.FindByID(product.ID.String())
	assert.Error(t, err)
	assert.Nil(t, productFound)
}
