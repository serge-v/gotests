package main

import (
	"regexp"
)

func main() {
	re := regexp.MustCompile("https?://[^ ]+")
	
	text := "http://localhost:6060/pkg/mime/multipart/#Part.Close"
	s := re.ReplaceAllString(text, "<a href=\"$0\">$0</a>")
	println(s)
}
