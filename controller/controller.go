package controller

import (
	"ConcurrentSearch/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"err": "长度不够",
		})
	} else {
		results := server.InputMessage(name)
		if len(results) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"err":     "空",
				"results": []string{},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"err":     "nil",
				"results": results,
			})
		}
	}
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
	tireInit()
}

func tireInit() {
	server.Constructor()
	server.ToSearch()
}
