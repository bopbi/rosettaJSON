package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	toJSON(os.Args)
}

func toJSON(args []string) {
	pathSeparator := string(os.PathSeparator)
	var inputFilename = ""
	var outputDir = ""
	if len(args) > 1 {
		fmt.Println("Generating output on current directory")
		inputFilename = args[1]
		if len(args) == 3 {
			outputDir = args[2]
			fmt.Println("Generating output on " + outputDir)
		}
	} else {
		fmt.Println("Please input input-file and optionally output-path")
		os.Exit(1)
	}

	excelFilePath, err := filepath.Abs(inputFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println("Opening From" + excelFilePath)

	xlFile, error := xlsx.OpenFile(excelFilePath)
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}

	if outputDir != "" {
		os.Mkdir(outputDir, 0777)
	}

	/*
		this app work by generating string xml by each language
	*/

	sheet := xlFile.Sheets[0] // only process first sheet

	var languages []string
	// get all available language on the first row
	languagesRow := sheet.Rows[0]
	for cellNumber, cell := range languagesRow.Cells {
		// skip the first cell
		if cellNumber > 0 {
			// insert the language code into the array
      var cellContent, _ = cell.String()
			languages = append(languages, cellContent)
		} else {
			continue
		}

	}

	var stringKey []string

	// save the string key on an array
	for rowNumber, row := range sheet.Rows {
		for cellNumber, cell := range row.Cells {
			// first colomn is for available languages
			var cellContent, _ = cell.String()
			if rowNumber > 0 {
				if (cellNumber == 0) && (cellContent != "") {
					stringKey = append(stringKey, cellContent)
				} else {
					continue
				}
			}
		}
	}

	// now write the xml one by one based on the languages
	for languageIndex, language := range languages {
    fmt.Printf("Working for language [%s] ", language)
		fmt.Println("")
		var stringContent string
    stringContent  = "{\n"
		for rowNumber, row := range sheet.Rows {
			for cellNumber, cell := range row.Cells {
				if rowNumber > 0 {
					var cellContent, _ = cell.String()
					if (cellNumber == languageIndex+1) && (cellContent != "") {
						name := stringKey[rowNumber-1]
            stringContent = stringContent + "\"" + name + "\" : \"" + cellContent + "\""
            // fmt.Println("rowNumber %d length %d", rowNumber, len(stringKey))
            if rowNumber != (len(stringKey)) {
              stringContent = stringContent + ","
            }
						stringContent = stringContent + "\n"
					} else {
						continue
					}
				} else {
					continue
				}

			}
		}
    stringContent  = stringContent + "} \n\n"

		outputFilename := "lang.json"
		if language != "" {
			outputFilename = "lang-" + language + ".json"
		}

		var generatedPath string
		if outputDir == "" {
			generatedPath = strings.Join([]string{outputFilename}, pathSeparator)
		} else {
			generatedPath = strings.Join([]string{outputDir, outputFilename}, pathSeparator)
		}

    file, _ := os.Create(generatedPath)
		n, err := io.WriteString(file, stringContent)
		if err != nil {
			fmt.Println(n, err)
		}
		file.Close()
		fmt.Printf("the json localizable working for language [%s] is generated", language)
		fmt.Println("")
	}
}
