package main

import (
	"fmt"
	"log"

	sr "github.com/blinchik/schema-registry/registry"
)

func main() {
	var schemaConfig sr.SchemaConfig

	schemaConfig.Address = "localhost"
	schemaConfig.Port = "8081"
	schemaConfig.Protocol = "http"

	const bookTickerSchema = `
	{
	"name":"test",
	"type":"record", 
		"fields":[
			{"name":"b","type":"double"},
			{"name":"B","type":"double"},
			{"name":"a","type":"double"},
			{"name":"A","type":"double"},
			{"name": "t", "type": "long", "logicalType": "timestamp-millis"}
		]
	}`

	res, err := sr.PostSchema(bookTickerSchema, "test-value", schemaConfig)

	bodyString := string(res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("POST Schema succeded: ", bodyString)

	sch, err := sr.GetSchemaLatest("test-value", schemaConfig)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GET Latet Schema succeded: ", sch)

}
