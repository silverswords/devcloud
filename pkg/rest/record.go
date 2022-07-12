package rest

import (
	"github.com/gin-gonic/gin"
)

type Record map[string]interface{}

func getRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Param("record_name")

		c.Next()

	}
}

func (app *RestfulAPI) GetRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) ListRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) CreateRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) DeleteRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) UpdateRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
