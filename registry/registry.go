package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//SchemaConfig struct defines the schema configuration
type SchemaConfig struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

// Schema struct defines the structure of returned schema from Schema Registry
type Schema struct {
	Subject string `json:"subject"`
	Version int    `json:"version"`
	ID      int    `json:"id"`
	Schema  string `json:"schema"`
}

//PostSchema will post schema to the Schema Registry
func PostSchema(schema string, config SchemaConfig) []byte {

	var objmap map[string]interface{}
	err := json.Unmarshal([]byte(schema), &objmap)
	if err != nil {
		log.Fatal(err)
	}

	objmap["schema"] = schema

	out, err := json.Marshal(objmap)

	if err != nil {
		log.Fatal(err)
	}

	body := strings.NewReader(string(out))
	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions", config.Protocol, config.Address, config.Port, config.Name)

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes
}

//GetSchemaLatest get the latest schema from Schema Registry and unmarshl it
func GetSchemaLatest(schema Schema, config SchemaConfig) (string, int) {

	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions/latest", config.Protocol, config.Address, config.Port, config.Name)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, &schema)
	if err != nil {
		log.Fatal(err)
	}

	schemaStr, err := json.Marshal(schema.Schema)
	if err != nil {
		log.Fatal(err)
	}

	schemaStrUnquote, err := strconv.Unquote(string(schemaStr))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	id := schema.ID

	return schemaStrUnquote, id
}
