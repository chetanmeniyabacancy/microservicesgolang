package generallib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
 Three ways of taking input
   1. fmt.Scanln(&input)
   2. reader.ReadString()
   3. scanner.Scan()

   Here we recommend using bufio.NewScanner
*/

func sort() {
	// To create dynamic array
	arr := make(map[int]string)
	var sortedarr []int
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < 6; i++ {
		fmt.Print("Enter Number " + strconv.Itoa(i) + ": ")
		// Scans a line from Stdin(Console)
		scanner.Scan()

		text := scanner.Text()
		if len(text) != 0 {
			arr[i] = text
		} else {
			i--
		}
	}

	for i := 0; i < 6; i++ {
		largenumber, largenumberposition := getlargetstnumberfromarray(arr)
		sortedarr = append(sortedarr, largenumber)
		delete(arr, largenumberposition)
	}
	fmt.Println("Sorted Array Is: ")
	fmt.Println(sortedarr)
}

func getlargetstnumberfromarray(arr map[int]string) (int, int) {
	var largenumber int
	var largenumberposition int
	var arrayfirst int64
	largenumberposition = 0
	arrayfirst, _ = strconv.ParseInt(arr[0], 10, 64)
	largenumber = int(arrayfirst)

	var arrayvvalueint int
	var arrayvvalue int64
	for i, value := range arr {

		arrayvvalue, _ = strconv.ParseInt(value, 10, 64)
		arrayvvalueint = int(arrayvvalue)

		if largenumber < arrayvvalueint {
			largenumber = arrayvvalueint
			largenumberposition = i
		}
	}
	return largenumber, largenumberposition
}
