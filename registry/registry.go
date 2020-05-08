package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)



const dockerHost "192.168.99.100"


func ListSchema(host string){


	host = dockerHost

	url := fmt.Sprintf("http://%s:8081/subjects", host)
	resp, err := http.Get(url)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)



	return bodyString


}


func PostSchema(idx int, schemaName string, schemaAddress string, schemas [][]byte, wg *sync.WaitGroup) {

	data := string(json.RawMessage(schemas[idx]))

	body := strings.NewReader(data)

	url := fmt.Sprintf("http://%s/subjects/%s/versions", schemaAddress, schemaName)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s %s %s\n", "[kapput]", "--->", schemaName)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {

		fmt.Printf("%s %s %s\n", "[ok]", "--->", schemaName)

	}

	wg.Done()

}

func AddSchema(schemas [][]byte, schemaNames []string, schemaAddress string) {

	fmt.Println("")
	wg := &sync.WaitGroup{}

	for idx, v := range schemaNames {
		wg.Add(1)

		go postSchema(idx, v, schemaAddress, schemas, wg)

	}

	wg.Wait()

}




func GetSchema(topic string, Type string) string {

	url := fmt.Sprintf("http://%s:%s/subjects/%s-%s/versions/latest", schemaAddress, schemaPort, topic, Type)
	resp, err := http.Get(url)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Header.Values())
	bodyString := string(bodyBytes)

	schemaID := resp.Header.Values("id")
	fmt.Println("+++++++ schema ID ", schemaID)

	return fixSchema(bodyString)
}











func fixSchema(schema string) string {

	schema = strings.Replace(schema, "}\"", "}", -1)
	schema = strings.Replace(schema, "\"{", "{", -1)
	schema = strings.Replace(schema, "\\n", "", -1)
	schema = strings.Replace(schema, "\\", "", -1)
	schema = strings.Replace(schema, " ", "", -1)

	var objmap map[string]interface{}
	err := json.Unmarshal([]byte(schema), &objmap)
	if err != nil {
		fmt.Println(err)
	}

	out, _ := json.Marshal(objmap["schema"])

	return string(out)
}
