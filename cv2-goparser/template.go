package main

import (
	"fmt"
	"html/template"
	"os"
)

func renderTemplate(file string, cv2 map[string]interface{}) error {
	tmpl, err := template.New("template.txt").ParseFiles(file)
	if err != nil {
		return err
	}
	fmt.Println("\n\n\n")
	err = tmpl.Execute(os.Stdout, cv2)
	if err != nil {
		return err
	}

	return nil
}
