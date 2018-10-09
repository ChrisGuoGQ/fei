package controllers

import (
	"encoding/base64"
	"fei/models"
	"fei/utils"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var cps map[uint]*websocket.Conn

func init() {
	cps = make(map[uint]*websocket.Conn)

}
func ConnectAll() {
	cameras, _ := models.IndexCamera()
	for _, camera := range *cameras {
		connectCP(camera)
	}
}
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

func StartCamera(c *gin.Context) {
	var camera models.Camera
	id, _ := strconv.Atoi(c.Param("id"))
	camera.ID = uint(id)
	if err := models.GetCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	connectCP(&camera)
	camera.State = true
	if err := models.UpdateCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}
func StopCamera(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	wc, _ := cps[uid]
	err := wc.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	delete(cps, uid)
	var camera models.Camera
	camera.ID = uid
	if err := models.GetCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	camera.State = false
	if err := models.UpdateCamera(&camera); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}
func connectCP(camera *models.Camera) {
	u := url.URL{Scheme: "ws", Host: camera.Proxy, Path: "/capture", RawQuery: "url=c3s-nat://" + camera.Admin + ":" + camera.Password + "@" + camera.Ip}
	wc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("read:", err)
		return
	}
	cfg := utils.GetCfg()

	go func() {
		for {
			_, message, err := wc.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			mt, err := jsonparser.GetString(message, "type")
			if err == nil && mt == "capture" {
				img, _ := jsonparser.GetString(message, "face", "crop", "image")
				params := map[string]io.Reader{
					"image": base64.NewDecoder(base64.StdEncoding, strings.NewReader(img)),
				}
				r, _ := utils.MultipartReq(cfg.Face.Minions+"/extract_with_detect", params)
				body, _ := ioutil.ReadAll(r.Body)
				featResult, _ := jsonparser.GetString(body, "faces", "[0]", "Feature")
				for _, group := range camera.Groups {
					searchParams := map[string]io.Reader{
						"groupname": strings.NewReader(group.Tag),
						"feature":   strings.NewReader(featResult),
						"limit":     strings.NewReader("1"),
					}
					searchResult, _ := utils.MultipartReq(cfg.Face.Group+"/group/search", searchParams)
					searchBody, _ := ioutil.ReadAll(searchResult.Body)
					log.Printf("searchBody: %s", searchBody)
					// searchId, _ := jsonparser.GetInt(searchBody, "ids", "[0]")
					score, _ := jsonparser.GetInt(searchBody, "scores", "[0]")
					// searchTag, _ := jsonparser.GetInt(searchBody, "tags", "[0]")
					if score > int64(group.Threshold) {
						break
					}
				}
			}
		}
	}()
	cps[camera.ID] = wc
}
