package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	cv2 := hackParse(os.ExpandEnv("$GOPATH/src/github.com/cv2me/cv2-tools/cv2-goparser/example.cv2"))
	b, err := json.MarshalIndent(cv2, "", " ")
	if err != nil {
		// Try my local path
		cv2 = hackParse(os.ExpandEnv("$GOPATH/src/github.com/temblortenor/cv2-tools/cv2-goparser/example.cv2"))
		if err != nil {
			log.Fatal("error:", err)
		}
	}
	fmt.Print(string(b))
}
