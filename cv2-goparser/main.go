package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	cv2 := hackParse(os.ExpandEnv("$GOPATH/src/github.com/cv2me/cv2-tools/cv2-goparser/example.cv2"))
	b, err := json.MarshalIndent(cv2, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
