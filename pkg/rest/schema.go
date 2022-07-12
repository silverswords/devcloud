// generate restful API for model with nosql database.
package rest

type Schema struct {
	fields []*SchemaFeild // ordered field
}

type SchemaFeild struct {
	Name string
	Type string
}
