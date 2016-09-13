package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	cv2, err := hackParse(os.ExpandEnv("$GOPATH/src/github.com/cv2me/cv2-tools/cv2-goparser/example.cv2"))
	if err != nil {
		log.Fatal(err)
	}

	// print as json for sanity check
	b, err := json.MarshalIndent(cv2, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(b))

	err = renderTemplate(os.ExpandEnv("$GOPATH/src/github.com/cv2me/cv2-tools/cv2-goparser/template.txt"), cv2)
	if err != nil {
		log.Fatal(err)
	}
}
