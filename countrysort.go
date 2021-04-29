package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func searchFile(file []byte, path string, countries []map[string]string, regionAbbv string) []map[string]string {

	regionalList := make([]map[string]string, 0)
	for _, c := range countries {

		param := []byte(c["English"])
		if bytes.Contains(file, param) {
			// fmt.Printf("Found %s in %s. Translation: %s\n", string(param), string(path), c["Translated"])

			entry := map[string]string{
				"continent": regionAbbv,
				// "countryCode": "code",
				"countryName": c["Translated"],
			}

			regionalList = append(regionalList, entry)
		}
	}
	return regionalList
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
	outputDir := args[2]

	files, err := os.ReadDir(dirPath)
	check(err)

	e := os.Mkdir(outputDir, 0744)
	check(e)

	for _, f := range files {
		fmt.Printf("File path: %s/%s\n", dirPath, f.Name())
		filepath := dirPath + "/" + f.Name()
		file, err := os.ReadFile(filepath)
		check(err)
		pathStrs := strings.Split(filepath, "/")
		region := strings.TrimSuffix(pathStrs[len(pathStrs)-1], ".json")

		regionMap := searchFile(file, filepath, countries, region)

		fp := path.Join(outputDir, region+".json")
		fmt.Println(fp)
		output, err := os.Create(fp)
		check(err)
		err = json.NewEncoder(output).Encode(regionMap)
		check(err)
	}
	// fmt.Printf("Files: %v", files)
}
