package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sr "github.com/blinchik-io/go-schema/registry"
)

// Firstschema asdsd
const Firstschema = `{"name":"test_topic11","type":"record", "fields":[{"name":"user","type":"string","size":12},{"name":"password","size":12,"type":"string"},{"name":"number","size":10,"type":["null","double"], "default": null }]}`
const dockerHost = "192.168.99.100"
const schemaPort = "8081"

var schemaList [][]byte

func main() {

	list := flag.Bool("l", false, "list")
	des := flag.Bool("des", false, "describe")
	delete := flag.Bool("delete", false, "delete")
	re := flag.Bool("r", false, "regex")
	add := flag.Bool("add", false, "add")

	flag.Parse()

	if *add {

		var schemaList [][]byte
		var schemaName []string

		files, err := ioutil.ReadDir(os.Args[2])

		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			schemaName = append(schemaName, strings.Split(f.Name(), ".")[0])

			dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", os.Args[2], f.Name()))
			if err != nil {
				log.Fatal(err)
			}

			schemaList = append(schemaList, dat)
		}

		sr.AddSchema(schemaList, schemaName, dockerHost, schemaPort)

		return

	}

	if *delete {

		if *re {

			subjects := sr.ListSchema(dockerHost, schemaPort, os.Args[3])
			sr.DeleteSchemaList(dockerHost, schemaPort, subjects)

			return
		} else {
			sr.DeleteSchemaSpecfic(dockerHost, schemaPort, os.Args[2])

			return
		}

	}

	if *list {

		subjects := sr.ListSchema(dockerHost, schemaPort, os.Args[2])

		for _, v := range subjects {

			fmt.Println(v)
		}

		return

	}

	if *des {

		subjectLatest := sr.FetcbLatest(dockerHost, schemaPort, os.Args[2])

		fmt.Println(subjectLatest)

		return

	}

}
