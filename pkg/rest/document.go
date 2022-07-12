package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BaseModel struct {
	// ReserveID        interface{} `docstore:"_id"`
	DocstoreRevision interface{}
}

type Document struct {
	BaseModel

	Name string

	Schema
}

func (d *Document) Fields() []string {
	var fields []string = make([]string, len(d.Schema))
	for _, field := range d.Schema {
		fields = append(fields, field.Name)
	}

	return fields
}

func (app *RestfulAPI) GetDocument(c *gin.Context) {
	_ = c.Param("id")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) ListDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) CreateDocument(c *gin.Context) {
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	logrus.Info(json, err)

	var document = &Document{
		Name: json["name"].(string),
		Schema: Schema{SchemaFeild{
			Name: "name",
			Type: "string",
		},
		},
	}

	if err := app.Coll.Create(context.Background(), document); err != nil {
		logrus.Info(err)
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) DeleteDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (app *RestfulAPI) UpdateDocument(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
