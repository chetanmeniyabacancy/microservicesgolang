package main

import "fmt"

// "bufio"

func main() {
	// var grades []int
	// grades = append(grades, 1, 2)
	// grades[1] = 88888

	// fmt.Printf("%v, %T", grades, grades)
	// fmt.Printf("%v, %v", len(grades), cap(grades))

	var grades string
	grades = "𩸽"
	gradebytes := []byte(grades)

	fmt.Printf("%v, %T", gradebytes, gradebytes)

}
