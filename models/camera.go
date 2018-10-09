package models

import (
	"github.com/jinzhu/gorm"
)

type Camera struct {
	gorm.Model
	Name     string  `json:"name" binding:"required"`
	Proxy    string  `json:"proxy" binding:"required"`
	Ip       string  `json:"ip" binding:"required"`
	Admin    string  `json:"admin" binding:"required"`
	Password string  `json:"password" binding:"required"`
	State    bool    `sql:"default:false" json:"state"`
	Groups   []Group `gorm:"many2many:camera_groups;" json:"groups"`
}

func IndexCamera() (*[]*Camera, error) {
	var cameras []*Camera
	err := db.Preload("Groups").Find(&cameras).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &cameras, nil
}

func CreateCamera(camera *Camera) error {
	if err := db.Omit("groups").Create(&camera).Error; err != nil {
		return err
	}
	if err := db.Model(&camera).Association("Groups").Append(camera.Groups).Error; err != nil {
		return err
	}
	return nil
}
func UpdateCamera(camera *Camera) error {
	if err := db.Omit("groups", "created_at").Save(&camera).Error; err != nil {
		return err
	}
	if err := db.Model(&camera).Association("Groups").Replace(camera.Groups).Error; err != nil {
		return err
	}
	return nil
}
func UpdateCameraByAttr(camera *Camera) error {
	if err := db.Omit("groups", "created_at").Save(&camera).Error; err != nil {
		return err
	}
	if err := db.Model(&camera).Association("Groups").Replace(camera.Groups).Error; err != nil {
		return err
	}
	return nil
}
func DestroyCamera(camera *Camera) error {
	if err := db.Model(&camera).Association("Groups").Clear().Error; err != nil {
		return err
	}
	if err := db.Delete(&camera).Error; err != nil {
		return err
	}
	return nil
}
func GetCamera(camera *Camera) error {
	if err := db.Preload("Groups").First(&camera, camera.ID).Error; err != nil {
		return err
	}
	return nil
}
