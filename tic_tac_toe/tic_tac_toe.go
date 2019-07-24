package main

import (
	"fmt"
	"strings"
)

type bored [3][3]string

func medallion(plywood bored) {
	for i := 0; i < 3; i++ {
		if i > 0 {
			fmt.Println("-----------")
		}
		fmt.Println(strings.Join(plywood[i][:], " | "))
	}
}

func main() {
	var plywood bored = bored{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}

	medallion(plywood)
}
