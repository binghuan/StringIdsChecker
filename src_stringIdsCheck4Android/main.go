package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type StringsRes struct {
	XMLName       xml.Name    `xml:"resources"`
	StringDefList []StringDef `xml:"string"`
	Locale        string
	StringIdsMap  map[string]string
}

type StringDef struct {
	Value string `xml:",chardata"`
	Name  string `xml:"name,attr"`
}

func main() {

	//folderPath := "../resource_for_testing/"

	if len(os.Args) < 2 {
		fmt.Println("Please give me the project's folder path you want to check")
		return
	}

	folderPath := os.Args[1]
	fmt.Println("check", folderPath)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		fmt.Println("<!-- ERROR -->: Folder was not found!")
	}

	//stringsResInLocales := make(map[string]StringsRes)
	var localeStringsResArray []StringsRes

	// apps/bnb-android/module-convert/convert-internal/src/binance/res/values/strings.xml
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//fmt.Println("Check:", path)

			var re = regexp.MustCompile(`values.*strings\.xml`)
			isL10nFileFound := re.MatchString(path)
			if isL10nFileFound {
				//fmt.Println("O>", path)
				dir := filepath.Dir(path)
				parent := filepath.Base(dir)
				//fmt.Println("folder=", parent)
				locale := strings.ReplaceAll(parent, "values-", "")
				if locale == "values" {
					locale = "en"
				}

				// Open xmlFile
				xmlFile, err := os.Open(path)

				// if we os.Open returns an error then handle it
				if err != nil {
					fmt.Println(err)
				}
				// defer the closing of our xmlFile so that we can parse it later on
				defer xmlFile.Close()

				// read our opened xmlFile as a byte array.
				byteValue, _ := ioutil.ReadAll(xmlFile)
				// we unmarshal our byteArray which contains our
				var stringsRes StringsRes
				xml.Unmarshal(byteValue, &stringsRes)
				stringsRes.Locale = locale
				stringsRes.StringIdsMap = make(map[string]string)
				for i := 0; i < len(stringsRes.StringDefList); i++ {
					stringName := stringsRes.StringDefList[i].Name
					stringValue := stringsRes.StringDefList[i].Value
					stringsRes.StringIdsMap[stringName] = stringValue
				}

				localeStringsResArray = append(localeStringsResArray, stringsRes)
			}

			return nil
		})

	if err != nil {
		log.Println(err)
	}

	result := true
	fmt.Println("Total languages: ", len(localeStringsResArray))
	for index, stringIds := range localeStringsResArray {
		nextStringIdsIndex := (index + 1) % len(localeStringsResArray)
		//fmt.Println("index:", index, stringIds.Locale, "with next index:", nextStringIdsIndex)

		if len(localeStringsResArray) > nextStringIdsIndex ||
			index == len(localeStringsResArray)-1 && nextStringIdsIndex == 0 {
			//fmt.Println("compare ", localeStringsResArray[index].Locale, localeStringsResArray[nextStringIdsIndex].Locale)
			for _, stringDef := range stringIds.StringDefList {

				if val, ok := localeStringsResArray[nextStringIdsIndex].StringIdsMap[stringDef.Name]; ok {
					//fmt.Println("stringId", key , "was found", val)
				} else {
					fmt.Println("<!-- ERROR   -->\n    stringId \""+stringDef.Name+"\" was not found in locale",
						localeStringsResArray[nextStringIdsIndex].Locale,
						", but it was found in locale", localeStringsResArray[index].Locale, val)
					result = false
				}
			}
		}
	}

	if result {
		fmt.Println("\n^_^b OK> All StringIds are uniquie")
	} else {
		fmt.Println("\nX_X! NG> Please check you string resources and try again!")
	}
}
