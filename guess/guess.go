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
	ans := r.Intn(73)
	fmt.Println(ans)
	fmt.Print("enter guess: ")
	var input string
	fmt.Scanln(&input)
	guess, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(guess)
	if guess == ans {
		fmt.Println("GREAT JOB")
	} else {
		fmt.Println("YOU SUCK AT THIS")
	}
}
