package controllers

import (
	"gintest/models"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexCamera(c *gin.Context) {
	cameras, _ := models.IndexCamera()
	c.JSON(http.StatusOK, *cameras)
}
func CreateCamera(c *gin.Context) {
	var camera models.Camera
	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.CreateCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, camera)
	return
}

func UpdateCamera(c *gin.Context) {
	var camera models.Camera
	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	camera.ID = uint(id)
	if err := models.UpdateCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}

func DestroyCamera(c *gin.Context) {
	var camera models.Camera
	id, _ := strconv.Atoi(c.Param("id"))
	camera.ID = uint(id)
	if err := models.DestroyCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
	return
}
