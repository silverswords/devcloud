package rest

import "github.com/gin-gonic/gin"

func GetRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func DeleteRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func UpdateRecord(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
