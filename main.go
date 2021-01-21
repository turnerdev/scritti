package main

import (
	"fmt"
	demo "scritti/demopackage"
	"scritti/substringpackage"
)

// Hello test
func Hello(name string) string {
	return "Hello, " + name
}

func main() {
	a := 65
	fmt.Println(demo.Id(1))
	fmt.Println(demo.UseToml())
	fmt.Println(a)
	fmt.Println(Hello("test"))
	fmt.Println(substringpackage.Reverse("world"))
}
