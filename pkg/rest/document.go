package rest

import (
	"context"
	"io"

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

	IDField string
	Schema  map[string]interface{}
}

func (d *Document) Record(data map[string]interface{}) map[string]interface{} {
	record := make(map[string]interface{})

	// todo(abserari): add type constraint.
	for name, _ := range d.Schema {
		record[name] = data[name]
	}

	record["DocstoreRevision"] = nil

	return record
}

func (d *Document) GetURL() URL {
	var url URL = URL{
		DB:         d.Name,
		Collection: d.Name,
		IDField:    d.IDField,
	}
	return url
}

func (app *RestfulAPI) GetDocument(c *gin.Context) {
	name := c.Param("id")

	var document = &Document{
		Name: name,
	}
	if err := app.Coll.Get(context.Background(), document); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": document})
}

func (app *RestfulAPI) ListDocument(c *gin.Context) {
	ctx := context.TODO()
	iter := app.Coll.Query().Get(ctx)

	defer iter.Stop()

	var documents []*Document
	// Query.Get returns an iterator. Call Next on it until io.EOF.
	for {
		var d Document
		err := iter.Next(ctx, &d)
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Error(err)
			c.JSON(400, gin.H{"err": err.Error()})
			return
		} else {
			documents = append(documents, &d)
		}
	}

	c.JSON(200, gin.H{"data": documents})
}

func (app *RestfulAPI) CreateDocument(c *gin.Context) {
	json := make(map[string]interface{})
	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	var document = &Document{
		Name:   json["name"].(string),
		Schema: json["schema"].(map[string]interface{}),
	}

	if err := app.Coll.Create(context.Background(), document); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": document,
	})
}

func (app *RestfulAPI) DeleteDocument(c *gin.Context) {
	name := c.Param("id")

	var document = &Document{
		Name: name,
	}
	if err := app.Coll.Delete(context.Background(), document); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}

func (app *RestfulAPI) UpdateDocument(c *gin.Context) {
	name := c.Param("id")
	json := make(map[string]interface{})

	err := c.BindJSON(&json)
	logrus.Info(json, err)

	var document = &Document{
		Name:   name,
		Schema: json["schema"].(map[string]interface{}),
	}

	if err := app.Coll.Actions().Put(document).Get(document).Do(context.Background()); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, gin.H{})
}
