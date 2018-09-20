package controllers

import (
	"gintest/models"
	"gintest/utils"
	"io"
	"strings"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexGroup(c *gin.Context) {
	groups, _ := models.IndexGroup()
	c.JSON(http.StatusOK, *groups)
}
func CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.CreateGroup(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := map[string]io.Reader{
		"groupname": strings.NewReader(group.Tag),
	}
	cfg := utils.GetCfg()
	_, err := utils.MultipartReq(cfg.Face.Group+"/group/init", params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Printf("%+v", res)
	c.JSON(http.StatusCreated, group)
	return
}
func UpdateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	group.ID = uint(id)
	if err := models.UpdateGroup(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}
func DestroyGroup(c *gin.Context) {
	var group models.Group
	id, _ := strconv.Atoi(c.Param("id"))
	group.ID = uint(id)
	if err := models.DestroyGroup(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := map[string]io.Reader{
		"groupname": strings.NewReader(group.Tag),
	}
	cfg := utils.GetCfg()
	_, err := utils.MultipartReq(cfg.Face.Group+"/group/free", params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
	return
}
