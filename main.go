package main

import (
	filerepository "scritti/filerepository"
)

// Hello test
func Hello(name string) string {
	return "Hello, " + name
}

func main() {
	filerepository.ReadFile()
}
