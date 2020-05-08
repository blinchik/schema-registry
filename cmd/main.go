package main

import (
	sr "github.com/blinchik-io/go-schema/registry"
)

// Firstschema asdsd
const Firstschema = `{"name":"test_topic10","type":"record", "fields":[{"name":"user","type":"string"},{"name":"password","size":10,"type":"string"},{"name":"number","size":10,"type":["null", "double"], "default": null }]}`

func main() {

	sr.AddSchema(Firstschema, SchemaNames, schemaAddress)

}
