package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	low := 0
	high := 50
	ans := r.Intn(high-low+1) + low
	i := 1
	for i <= 7 {
		fmt.Printf("enter guess between %d and %d: ", low, high)
		var input string
		fmt.Scanln(&input)
		guess, err := strconv.Atoi(input)
		switch {
		case err != nil:
			fmt.Println(err)
			continue
		case guess > high:
			fmt.Println("too high")
			continue
		case guess < low:
			fmt.Println("too low")
			continue
		case guess == ans:
			fmt.Println("GREAT JOB")
			return
		case guess < ans:
			fmt.Println("higher")
			low = guess + 1
		case guess > ans:
			fmt.Println("lower")
			high = guess - 1
		}
		i = i + 1
	}
	fmt.Println("YOU SUCK AT THIS")
}
