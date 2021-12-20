package main

import (
	"log"
	"sync"
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

func sq(workerId int, in <-chan int) <-chan int {
	results := make(chan int)
	go func() {
		for num := range in {
			log.Printf("sq(%v): processing %v", workerId, num)
			results <- num * num
		}
		log.Printf("sq(%v): closing results chan", workerId)
		close(results)
	}()
	return results
}

func mergeOutputs(outChans ...<-chan int) <-chan int {
	wg := new(sync.WaitGroup)
	results := make(chan int)

	// anonymous function to read from each output chan in arguments
	consume := func(out <-chan int) {
		for num := range out {
			results <- num
		}
		wg.Done()
	}

	// wait to read from all the channels
	wg.Add(len(outChans))
	for _, outChan := range outChans {
		// fan-in all the outputs!
		go consume(outChan)
	}

	// Request a gopher to close results chan after all the output goroutines are
	// done
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func main() {
	log.Println("Hello Gophers!")

	nums := []int{1, 2, 3}
	log.Println("input: ", nums)

	in := gen(nums...)
	// fan-out: request (be nice to them) 2 gophers to process the input
	firstGopherOutputs := sq(1, in)
	secondGopherOutputs := sq(2, in)

	// fan-in: combine outputs from the 2 gophers onto 1 output channel
	// 'r' can be the result from either the 1st or 2nd gopher. Order is not
	// guaranteed
	for r := range mergeOutputs(firstGopherOutputs, secondGopherOutputs) {
		log.Printf("result: %v", r)
	}

	log.Printf("Done with loops")
}
