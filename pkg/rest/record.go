package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gocloud.dev/docstore"
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
	name := c.Param("record_name")
	document := &Document{
		Name: name,
	}

	if err := app.Coll.Get(context.Background(), document); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	url := Mongo(document.GetURL())

	var err error
	Collection, err = docstore.OpenCollection(context.Background(), url)
	if err != nil {
		logrus.Fatal(err)
	}
	defer Collection.Close()

	json := make(map[string]interface{})
	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	record := document.Record(json)
	if err := Collection.Create(context.Background(), record); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": record})
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
