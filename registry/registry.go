package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//SchemaConfig struct defines the schema configuration
type SchemaConfig struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type schema struct {
	Subject string `json:"subject"`
	Version int    `json:"version"`
	ID      int    `json:"id"`
	Schema  string `json:"schema"`
}

//PostSchema will post schema to the Schema Registry
func PostSchema(schema, name string, config SchemaConfig) (respBody []byte, err error) {

	var objmap map[string]interface{}

	err = json.Unmarshal([]byte(schema), &objmap)
	if err != nil {
		return
	}

	objmap["schema"] = schema

	out, err := json.Marshal(objmap)
	if err != nil {
		return
	}

	body := strings.NewReader(string(out))
	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions", config.Protocol, config.Address, config.Port, name)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return respBody, err
}

//GetSchemaLatest get the latest schema from Schema Registry and unmarshl it
func GetSchemaLatest(name string, config SchemaConfig) (sch schema, err error) {

	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions/latest", config.Protocol, config.Address, config.Port, name)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &sch)
	if err != nil {
		return
	}

	schemaStr, err := json.Marshal(sch.Schema)
	if err != nil {
		return
	}

	schemaStrUnquote, err := strconv.Unquote(string(schemaStr))
	if err != nil {
		return
	}

	sch.Schema = schemaStrUnquote

	return sch, err
}
