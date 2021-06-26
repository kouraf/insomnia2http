package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// grab the file path from the args
	jsonFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	// import the json file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()

	var collection InsomniaCollection

	// Unmarshal it to a type
	json.Unmarshal(byteValue, &collection)
	var groups []resource

	for _, resource := range collection.Resources {
		if resource.Type == "request_group" {
			groups = append(groups, resource)
		}
	}

	// create an out string
	for _, group := range groups {
		outString := ""
		for idx, resource := range collection.Resources {
			if resource.ParentId == group.Id {

				outString += resource.Method + " " + resource.Url + "\n"
				for _, header := range resource.Headers {
					outString += header.Name + ": " + header.Value + "\n"
				}
				if resource.Body.Data != "" {
					outString += "\n" + resource.Body.Data + "\n"
				}

				// TODO: Counter to not do this on the last one
				if idx+1 != len(collection.Resources) {
					outString += "\n###\n\n"
				}
			}
		}
		// save to the disk
		// f, err := os.Create(collection.Name + ".http")
		f, err := os.Create(group.Name + ".http")
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = f.WriteString(outString)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

}

type InsomniaCollection struct {
	Name      string     `json:"__export_source"`
	Resources []resource `json:"resources"`
}

type resource struct {
	Id          string   `json:"_id"`
	ParentId    string   `json:"parentId"`
	Url         string   `json:"url"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Method      string   `json:"method"`
	Body        Body     `json:"body"`
	Headers     []Header `json:"headers"`
	Type        string   `json:"_type"`
}

type Body struct {
	Data string `json:"text"`
}
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
