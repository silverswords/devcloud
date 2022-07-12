package rest

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()

	// todo(abserari): add group for each resource
	router.GET("/records/:record_name/:id", GetRecord)
	router.GET("/records/:record_name", ListRecord)
	router.POST("/records/:record_name", CreateRecord)
	router.DELETE("/records/:record_name/:id", DeleteRecord)
	router.PUT("/records/:record_name/:id", UpdateRecord)

	router.GET("/collections", ListDocument)
	router.GET("/collections/:id", GetDocument)
	router.POST("/collections", CreateDocument)
	router.DELETE("/collections/:id", DeleteDocument)
	router.PUT("/collections/:id", UpdateDocument)

	router.Run("0.0.0.0:8080")
}
