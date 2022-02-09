package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB opens a connection to the database
func main() {
	channelexample()
}

func channelexample() {
	ch := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go func() { ch <- 15 }()
	}

	var num1 int
	for i := 0; i < 10; i++ {
		num1 = <-ch
	}

	fmt.Println(num1)
}
