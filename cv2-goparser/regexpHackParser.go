package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

// This is a POC horrible hack, I am moving towards a proper parser.
// Uses regexp and all sort of nastiness
func hackParse(file string) (map[string]interface{}, error) {
	cv2 := make(map[string]interface{})

	aliasRegexp, err := regexp.Compile(`\s*(.+?)\s+AS\s+(.+?)\s*\n`)
	if err != nil {
		return cv2, err
	}
	tagsRegexp, err := regexp.Compile(`\s*(.+?)\s*:\s*\[(.+?)\]\s*\n`)
	if err != nil {
		return cv2, err
	}
	rangeRegexp, err := regexp.Compile(`(\d+)-(\d+)`)
	if err != nil {
		return cv2, err
	}
	sectionRegexp, err := regexp.Compile(`\s*\\\\\s*(.+?)\s*\n`)
	if err != nil {
		return cv2, err
	}
	//dottedDateLineRegexp, err := regexp.Compile(`\s*(.+?)\.(.+?)\s*:\s*(\d\d\.\d\d\.\d\d\d\d)\s*\n`)
	dottedDateLineRegexp, err := regexp.Compile(`\s*(.+?)\.(.+?)\s*:\s*(\d\d\.\d\d\.\d\d\d\d)\s*\n`)
	if err != nil {
		return cv2, err
	}
	dateLineRegexp, err := regexp.Compile(`\s*(.+?)\s*:\s*(\d\d\.\d\d\.\d\d\d\d)\s*\n`)
	if err != nil {
		return cv2, err
	}
	dottedRegexp, err := regexp.Compile(`\s*(.+?)\.(.+?)\s*:\s*(.+)\s*\n`)
	if err != nil {
		return cv2, err
	}
	otherRegexp, err := regexp.Compile(`\s*(.+?)\s*:\s*(.+)\s*\n`)
	if err != nil {
		return cv2, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cv2, err
	}
	s := string(b)
	inComment := false
	inEscape := false
	inString := false
	pendingString := ""
	pendingComment := ""
	pendingUnknown := ""
	pendingEscape := ""
	currentSection := ""
	lastString := ""
	lastEscape := ""
	for index, runeValue := range s {
		if runeValue == '/' && s[index+1] == '*' {
			inComment = true
			log.Println("inComment: ", inComment)
		}
		// TODO: will break if index+1>len
		if runeValue == '*' && s[index+1] == '/' && inComment {
			inComment = false
			log.Println("inComment: ", inComment)
			//log.Println("comment: ", pendingComment)
			pendingComment = ""
		}
		if inComment {
			pendingComment += string(runeValue)
			continue
		}

		if runeValue == '\\' && s[index+1] == '{' && !inEscape {
			inEscape = true
			log.Println("inEscape: ", inEscape)
		}
		// TODO: will break if index+1>len
		if runeValue == '}' && inEscape {
			inEscape = false
			log.Println("inEscape: ", inEscape)
			lastEscape = pendingEscape[2:]
			log.Println("escape: ", lastEscape)
			pendingEscape = ""
		}
		if inEscape {
			pendingEscape += string(runeValue)
			continue
		}

		if runeValue == '"' {
			if inString {
				inString = false
				// TODO: will break on 0 length strings
				lastString = pendingString[1:]
				log.Println("lastString:", lastString)
				pendingString = ""
			} else {
				inString = true
			}
			log.Println("inString:", inString)
		}

		if inString {
			pendingString += string(runeValue)
		} else if inEscape {
			pendingEscape += string(runeValue)
		} else {
			pendingUnknown += string(runeValue)
		}

		foundSection := sectionRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundSection) > 0 {
			currentSection = strings.Trim(foundSection[0][1], "\n ")
			log.Printf("foundSection %s", currentSection)
			cv2[currentSection] = make(map[string]interface{})
			pendingUnknown = ""
		}

		if currentSection == "" {
			continue
		}

		foundAs := aliasRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundAs) > 0 {
			value := strings.Replace(strings.Trim(foundAs[0][2], " "), `"`, lastString, 1)
			value = strings.Replace(value, `}`, `\{`+lastEscape+`}`, 1)
			log.Printf("foundAs %s AS %s", foundAs[0][1], value)
			pendingUnknown = ""

			ref := cv2[currentSection].(map[string]interface{})
			if _, ok := ref["cv2:aliases"]; !ok {
				ref["cv2:aliases"] = make(map[string]string)
			}

			ref2 := ref["cv2:aliases"].(map[string]string)
			ref2[foundAs[0][1]] = value
		}

		foundTag := tagsRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundTag) > 0 {
			value := strings.Replace(strings.Trim(foundTag[0][2], " "), `"`, lastString, 1)
			value = strings.Replace(value, `}`, `\{`+lastEscape+`}`, 1)
			log.Printf("foundTag %s AS %s", foundTag[0][1], value)
			pendingUnknown = ""
			foundRange := rangeRegexp.FindAllStringSubmatch(value, -1)
			if len(foundRange) > 0 {
				log.Printf("foundRange %s to %s", foundRange[0][1], foundRange[0][2])
			} else {
				enums := strings.Split(value, ",")
				log.Println("Enums are ")
				for _, enum := range enums {
					log.Println("\t", strings.Trim(enum, " "))
				}
			}
		}

		foundDottedDateLine := dottedDateLineRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundDottedDateLine) > 0 {
			t, err := time.Parse("02.01.2006", foundDottedDateLine[0][3])
			if err != nil {
				return cv2, err
			}
			log.Printf("foundDottedDateLine %s.%s %s", foundDottedDateLine[0][1], foundDottedDateLine[0][2], t.Format("01/02/2006"))
			pendingUnknown = ""

			ref := cv2[currentSection].(map[string]interface{})
			if _, ok := ref["cv2:values"]; !ok {
				ref["cv2:values"] = make(map[string]interface{})
			}

			ref2 := ref["cv2:values"].(map[string]interface{})
			if _, ok := ref2[foundDottedDateLine[0][1]]; !ok {
				ref2[foundDottedDateLine[0][1]] = make(map[string]string)
			}

			ref3 := ref2[foundDottedDateLine[0][1]].(map[string]string)
			ref3[foundDottedDateLine[0][2]] = t.Format("01/02/2006")
		}

		foundDateLine := dateLineRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundDateLine) > 0 {
			t, err := time.Parse("02.01.2006", foundDateLine[0][2])
			if err != nil {
				return cv2, err
			}
			log.Printf("foundDateLine %s %s", foundDateLine[0][1], t.Format("01/02/2006"))
			pendingUnknown = ""

			ref := cv2[currentSection].(map[string]interface{})
			if _, ok := ref["cv2:values"]; !ok {
				ref["cv2:values"] = make(map[string]string)
			}

			ref2 := ref["cv2:values"].(map[string]string)
			ref2[foundDateLine[0][1]] = t.Format("01/02/2006")
		}

		foundDotted := dottedRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundDotted) > 0 {
			value := strings.Replace(strings.Trim(foundDotted[0][3], " "), `"`, lastString, 1)
			value = strings.Replace(value, `}`, `\{`+lastEscape+`}`, 1)
			log.Printf("foundDotted %s.%s: %s", foundDotted[0][1], foundDotted[0][2], value)
			pendingUnknown = ""

			ref := cv2[currentSection].(map[string]interface{})
			if _, ok := ref["cv2:values"]; !ok {
				ref["cv2:values"] = make(map[string]interface{})
			}

			ref2 := ref["cv2:values"].(map[string]interface{})
			if _, ok := ref2[foundDotted[0][1]]; !ok {
				ref2[foundDotted[0][1]] = make(map[string]string)
			}

			ref3 := ref2[foundDotted[0][1]].(map[string]string)
			ref3[foundDotted[0][2]] = value
		}

		foundOther := otherRegexp.FindAllStringSubmatch(pendingUnknown, -1)
		if len(foundOther) > 0 {
			value := strings.Replace(strings.Trim(foundOther[0][2], " "), `"`, lastString, 1)
			value = strings.Replace(value, `}`, `\{`+lastEscape+`}`, 1)
			log.Printf("foundOther %s: %s", foundOther[0][1], value)
			pendingUnknown = ""

			ref := cv2[currentSection].(map[string]interface{})
			if _, ok := ref["cv2:values"]; !ok {
				ref["cv2:values"] = make(map[string]string)
			}

			ref2 := ref["cv2:values"].(map[string]string)
			ref2[foundOther[0][1]] = value
		}
		//log.Println("pendingUnknown", pendingUnknown)
	}

	return cv2, nil
}
