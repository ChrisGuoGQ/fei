package controllers

import (
	"fei/models"
	"fei/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func IndexBase(c *gin.Context) {
	bases, _ := models.IndexBase()
	c.JSON(http.StatusOK, *bases)
}

func AddBases(c *gin.Context) {
	con := c.DefaultQuery("con", "10")
	group := c.Query("group")
	files, err := ioutil.ReadDir("pics")
	if err != nil {
		log.Fatal(err)
	}
	cfg := utils.GetCfg()
	conNum, _ := strconv.Atoi(con)
	requests := make(chan string, conNum)
	for _, f := range files {
		requests <- f.Name()
		go func() {
			tag := f.Name()
			file, _ := os.Open(path.Join("pics", f.Name()))
			params := map[string]io.Reader{
				"image": file,
			}
			r, _ := utils.MultipartReq(cfg.Face.Minions+"/extract_with_detect", params)
			body, _ := ioutil.ReadAll(r.Body)
			// log.Printf("extractBody: %s", body)
			featResult, _ := jsonparser.GetString(body, "faces", "[0]", "Feature")
			base := models.Base{Name: tag, Tag: tag, GroupName: group, Feature: featResult}
			models.CreateBase(&base)
			addParams := map[string]io.Reader{
				"groupname": strings.NewReader(group),
				"featureid": strings.NewReader(strconv.Itoa(int(base.ID))),
				"tag":       strings.NewReader(tag),
				"feature":   strings.NewReader(featResult),
			}
			addResult, _ := utils.MultipartReq(cfg.Face.Group+"/group/add", addParams)
			addBody, _ := ioutil.ReadAll(addResult.Body)
			log.Printf("addBody: %s", addBody)
			done := <-requests
			log.Info(done)
		}()
	}
	c.Status(http.StatusOK)
}
func SyncBases(c *gin.Context) {
	cfg := utils.GetCfg()
	groups, _ := models.IndexGroup()
	for _, group := range *groups {
		freeParams := map[string]io.Reader{
			"groupname": strings.NewReader(group.Tag),
		}
		utils.MultipartReq(cfg.Face.Group+"/group/free", freeParams)
		initParams := map[string]io.Reader{
			"groupname": strings.NewReader(group.Tag),
		}
		utils.MultipartReq(cfg.Face.Group+"/group/init", initParams)
	}
	bases, _ := models.IndexBase()
	for _, base := range *bases {
		addParams := map[string]io.Reader{
			"groupname": strings.NewReader(base.GroupName),
			"tag":       strings.NewReader(base.Tag),
			"featureid": strings.NewReader(strconv.Itoa(int(base.ID))),
			"feature":   strings.NewReader(base.Feature),
		}
		addResult, _ := utils.MultipartReq(cfg.Face.Group+"/group/add", addParams)
		addBody, _ := ioutil.ReadAll(addResult.Body)
		log.Printf("addBody: %s", addBody)
	}
	c.Status(http.StatusOK)
}
