package main

import "fmt"

func main() {

	inp := make(chan interface{},100)

	out := make(chan interface{},100)


	//go CombineResults(inp, out)
	//
	//inp <- "29568666068035183841425683795340791879727309630931025356555"
	//inp <- "4958044192186797981418233587017209679042592862002427381542"
	//close(inp)
	//output := <- out
	//
	//fmt.Println(output)

	go SingleHash(inp, out)

	inp <- 0
	inp <- 1
	inp <- 2
	close(inp)

	for output := range out{
		fmt.Println(output)
	}

}
