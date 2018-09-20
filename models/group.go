package models

import (
	"github.com/jinzhu/gorm"
)

type Group struct {
	gorm.Model
	Name      string `json:"name" binding:"required"`
	Tag       string `json:"tag" binding:"required"`
	Threshold int    `json:"threshold" binding:"required"`
}

func IndexGroup() (*[]*Group, error) {
	var groups []*Group
	err := db.Find(&groups).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &groups, nil
}

func CreateGroup(group *Group) error {
	if err := db.Create(&group).Error; err != nil {
		return err
	}
	return nil
}
func UpdateGroup(group *Group) error {
	if err := db.Omit("created_at").Save(&group).Error; err != nil {
		return err
	}
	return nil
}
func DestroyGroup(group *Group) error {
	if err := db.Delete(&group).Error; err != nil {
		return err
	}
	return nil
}
