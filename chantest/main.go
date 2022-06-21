package main

import "fmt"

func main() {
	chanlist := make([]chan interface{}, 10)
	for _, x := range chanlist {

		if x == chanlist[0] {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	}

}
