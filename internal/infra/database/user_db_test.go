package database

import (
	"github.com/carloseduribeiro/crud-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestNewUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	if err = db.AutoMigrate(&entity.User{}); err != nil {
		t.Error(err)
	}
	user, _ := entity.NewUser("Jhon", "j@j.com", "123456")
	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.Nil(t, err)
	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Password, userFound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	if err = db.AutoMigrate(&entity.User{}); err != nil {
		t.Error(err)
	}
	user, _ := entity.NewUser("Jhon", "j@j.com", "123456")
	userDB := NewUser(db)
	err = userDB.Create(user)
	assert.Nil(t, err)
	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Password, userFound.Password)
}
