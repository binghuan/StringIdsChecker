package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*
StringIds ...
*/
type StringIds struct {
	locale string
	ids    map[string]string
}

var stringIdsArray []StringIds

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

	result := true
	isStringFileExisted := false

	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, "Localizable.strings") {
				fmt.Println(path)
				var re = regexp.MustCompile(`.lproj.*`)
				pretemp := re.ReplaceAllString(path, ``)
				var re2 = regexp.MustCompile(`.*/`)
				locale := re2.ReplaceAllString(pretemp, ``)
				fmt.Println("... checking locale \"" + locale + "\"")

				// @BH_Lin ---------------------------------------------------->
				// Read file line by line
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				isStringFileExisted = true

				keymap := make(map[string]string)

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					lineText := strings.TrimSpace(scanner.Text())
					if strings.Contains(lineText, "=") && !strings.HasPrefix(lineText, "//") {
						//fmt.Println(lineText)
						items := strings.Split(lineText, "=")
						//key, value := strings.TrimSpace(items[0]), strings.TrimSpace(items[1])

						//fmt.Println(items)
						key := items[0]
						//vaue := items[1]
						//fmt.Println(key)

						//fmt.Println("length: ", len(items))
						//fmt.Println("key: ", key, "value:", value)
						if value, ok := keymap[key]; ok {
							fmt.Println("<!-- WARNING -->:\n    stringId ", key, "is duplicate! in file \"", path, "\"")
							result = false
						} else {
							keymap[key] = value
						}
					}
				}

				stringIds := StringIds{locale, keymap}
				stringIdsArray = append(stringIdsArray, stringIds)

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
				// @BH_Lin ----------------------------------------------------<

			}

			return nil
		})
	if err != nil {
		log.Println(err)
	} else {

		fmt.Println("Total languages: ", len(stringIdsArray))
		for index, stringIds := range stringIdsArray {
			nextStringIdsIndex := (index + 1) % len(stringIdsArray)
			fmt.Println("index:", index, stringIds.locale, "with next index:", nextStringIdsIndex)

			if len(stringIdsArray) > nextStringIdsIndex ||
				index == len(stringIdsArray)-1 && nextStringIdsIndex == 0 {

				fmt.Println("compare ", stringIdsArray[index].locale, stringIdsArray[nextStringIdsIndex].locale)

				for key := range stringIds.ids {

					if val, ok := stringIdsArray[nextStringIdsIndex].ids[key]; ok {
						//fmt.Println("stringId", key , "was found", val)
					} else {
						fmt.Println("<!-- ERROR   -->\n    stringId \""+key+"\" was not found in locale",
							stringIdsArray[nextStringIdsIndex].locale,
							", but it was found in locale", stringIdsArray[index].locale, val)
						result = false
					}
				}
			}
		}
	}

	if !isStringFileExisted {
		fmt.Println("\nX> There is no strings file.")
	}
	if result == true {
		fmt.Println("\n^_^b OK> All StringIds are uniquie")
	} else {
		fmt.Println("\nX_X! NG> Please check you string resources and try again!")
	}
}
