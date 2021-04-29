package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func searchFile(file []byte, path string, countries []map[string]string) {
	for _, c := range countries {
		fmt.Printf("%v\n", c)
		json.NewEncoder(os.Stdout).Encode(c)
		// _, ok := c["English"]
		// fmt.Println(ok)
	}
	param := []byte("Aruba")
	if bytes.Contains(file, param) {
		// fmt.Printf("Found %s in %s \n", string(param), string(path))
	}
}

func getTrans(filename string) []map[string]string {

	raw, err := os.Open(filename)
	check(err)
	reader := csv.NewReader(raw)

	data, e := reader.ReadAll()
	check(e)

	processedData := make([]map[string]string, len(data)-1)
	for i, row := range data {
		if i > 0 {
			newMap := make(map[string]string)
			for k, val := range row {
				newMap[data[0][k]] = val
			}
			processedData[i-1] = newMap
		}
	}

	return processedData
}

func main() {
	flag.Parse()
	args := flag.Args()
	dirPath := args[0]
	transCSV := args[1]
	countries := getTrans(transCSV)

	files, err := os.ReadDir(dirPath)
	check(err)

	for _, f := range files {
		fmt.Printf("File path: %s/%s\n", dirPath, f.Name())
		filepath := dirPath + "/" + f.Name()
		file, err := os.ReadFile(filepath)
		check(err)

		searchFile(file, filepath, countries)
	}
	// fmt.Printf("Files: %v", files)
}
