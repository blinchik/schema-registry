package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	av "github.com/blinchik/schema-registry/avro"
	sr "github.com/blinchik/schema-registry/registry"
)

type record struct {
	BB float64 `json:"b"`
	B  float64 `json:"B"`
	AA float64 `json:"a"`
	A  float64 `json:"A"`
	T  int64   `json:"t"`
}

func main() {

	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)

	var schemaConfig sr.SchemaConfig

	schemaConfig.Address = "localhost"
	schemaConfig.Port = "8081"
	schemaConfig.Protocol = "http"

	sch, err := sr.GetSchemaLatest("test-value", schemaConfig)
	if err != nil {
		log.Fatal(err)
	}

	msg := record{BB: 1.0, B: 1.0, AA: 1.0, A: 1.0, T: time.Now().UnixMilli()}

	msgByte, err := json.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	confluentSchemaData, err := av.SerializeMessageConfluentAvro(sch.Schema, sch.ID, msgByte)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SerializeMessageConfluentAvro succeded", string(confluentSchemaData))

	textual, err := av.DeserializeMessageConfluentAvro(sch.Schema, sch.ID, confluentSchemaData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DeserializeMessageConfluentAvro succeded", string(textual))

}
