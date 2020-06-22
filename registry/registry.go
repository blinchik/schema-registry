package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

func FetchLatest(host, schemaPort, schemaName string, protocol string) string {

	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions/latest",protocol, host, schemaPort, schemaName)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	return FixSchema(string(json.RawMessage(bodyBytes)))

}

func DeleteSchemaSpecfic(host, schemaPort, schemaName string, protocol string) {

	url := fmt.Sprintf("%s://%s:%s/subjects/%s",protocol, host, schemaPort, schemaName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("\nDelete: ", schemaName, "\n")

}

func DeleteSchemaList(host, schemaPort string, schemaList []string, protocol string) {

	wg := &sync.WaitGroup{}

	for _, v := range schemaList {
		wg.Add(1)

		go deleteSchema(host, schemaPort, v, wg, protocol)

	}

	wg.Wait()

}

func deleteSchema(host, schemaPort, schemaName string, wg *sync.WaitGroup,protocol string) {

	url := fmt.Sprintf("%s://%s:%s/subjects/%s",protocol, host, schemaPort, schemaName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf("%s %s %s\n", "[Delete]", "---> ", schemaName)

	wg.Done()

}

func ListSchema(host, schemaPort, regex string,protocol string) []string {

	var stringList []string
	var schemaList []string

	url := fmt.Sprintf("%s://%s:%s/subjects",protocol, host, schemaPort)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bodyBytes, &stringList)

	re := regexp.MustCompile(regex)

	for index := range stringList {

		if len(re.FindString(stringList[index])) != 0 {

			schemaList = append(schemaList, stringList[index])

		}
	}

	return schemaList

}

func postSchema(idx int, schemaName, host, schemaPort string, schemas [][]byte, wg *sync.WaitGroup,protocol string) {

	data := string(json.RawMessage(schemas[idx]))

	body := strings.NewReader(data)

	url := fmt.Sprintf("%s://%s:%s/subjects/%s/versions",protocol, host, schemaPort, schemaName)
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
		fmt.Printf("%s %s %s\n", "[error]", "--->  ", schemaName)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {

		fmt.Printf("%s %s %s\n", "[ok]", "--->  ", schemaName)

	}

	wg.Done()

}

func AddSchema(schemas [][]byte, schemaNames []string, host, schemaPort string, protocol string) {

	fmt.Println("")
	wg := &sync.WaitGroup{}

	for idx, v := range schemaNames {
		wg.Add(1)

		go postSchema(idx, v, host, schemaPort, schemas, wg, protocol)

	}

	wg.Wait()

}

func FixSchema(schema string) string {

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

	// out, _ := json.Marshal(objmap["schema"])

	out, _ := json.Marshal(objmap)

	return string(out)
}

func StrToSchema(schema string) string {

	return fmt.Sprintf(`{"schema": "%s"}`, strings.ReplaceAll(schema, `"`, `\"`))

}
