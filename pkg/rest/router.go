package rest

import "github.com/gin-gonic/gin"

func Run() {
	router := gin.Default()
	app := NewRestfulAPI()
	// todo(abserari): add group for each resource
	records := router.Group("/record")
	{
		records.Use(getRecord())
		records.GET("/:record_name/:id", app.GetRecord)
		records.GET("/:record_name", app.ListRecord)
		records.POST("/:record_name", app.CreateRecord)
		records.DELETE("/:record_name/:id", app.DeleteRecord)
		records.PUT("/:record_name/:id", app.UpdateRecord)

	}

	documents := router.Group("/document")
	{
		documents.GET("/", app.ListDocument)
		documents.GET("/:id", app.GetDocument)
		documents.POST("/", app.CreateDocument)
		documents.DELETE("/:id", app.DeleteDocument)
		documents.PUT("/:id", app.UpdateDocument)
	}

	router.Run("0.0.0.0:8080")
}
