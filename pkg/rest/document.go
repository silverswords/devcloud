package rest

import "github.com/gin-gonic/gin"

func GetDocument(c *gin.Context) {
	_ = c.Param("record_name")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ListDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func DeleteDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func UpdateDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
