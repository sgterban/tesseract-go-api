package main

import (
	"fmt"
	"github.com/otiai10/gosseract"
	"os"
)

func main() {
	file := os.Args[1]
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(file)
	text, _ := client.Text()
	fmt.Println(text)
}
