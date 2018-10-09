package models

import (
	"github.com/jinzhu/gorm"
)

type Base struct {
	gorm.Model
	Name      string `json:"name" binding:"required"`
	Tag       string `json:"tag" binding:"required"`
	GroupName string `json:"groupName" binding:"required"`
	Feature   string `json:"feature" binding:"required"`
}

func IndexBase() (*[]*Base, error) {
	var bases []*Base
	err := db.Find(&bases).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &bases, nil
}

func CreateBase(base *Base) error {
	if err := db.Create(&base).Error; err != nil {
		return err
	}
	return nil
}
func UpdateBase(base *Base) error {
	if err := db.Omit("created_at").Save(&base).Error; err != nil {
		return err
	}
	return nil
}
func DestroyBase(base *Base) error {
	if err := db.Delete(&base).Error; err != nil {
		return err
	}
	return nil
}
