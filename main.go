package main

import (
	"log"
	"time"
)

func gen(nums ...int) <-chan int {
	input := make(chan int)
	go func() {
		for _, num := range nums {
			log.Printf("gen: sending %v", num)
			input <- num
		}
		log.Println("gen: closing input chan")
		close(input)
	}()
	return input
}

func sq(in <-chan int) <-chan int {
	results := make(chan int)
	go func() {
		for num := range in {
			log.Printf("sq: processing %v", num)
			results <- num * num
		}
		log.Println("sq: closing results chan")
		close(results)
	}()
	return results
}

func add(in <-chan int) <-chan int {
	results := make(chan int)
	go func() {
		for num := range in {
			log.Printf("add: processing %v", num)
			results <- num + num
		}
		log.Println("add: closing results chan")
		close(results)
	}()
	return results
}

func main() {
	log.Println("Hello Gophers!")

	nums := []int{1, 2, 3}
	log.Println("input: ", nums)

	results := add(sq(gen(nums...)))

	for i := 0; i < len(nums); i++ {
		log.Printf("results[%v]: %v", i, <-results)
	}

	time.Sleep(2 * time.Second)
}
