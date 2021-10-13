// TODO:
// - handle eof terminal character

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// NewReader returns a new Reader whose buffer has the default size.
	//   https://pkg.go.dev/bufio#NewReader
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("▶️  ")

		// ReadString reads until the first occurrence of delim in the input, returning a string
		// containing the data up to and including the delimiter.
		//   https://pkg.go.dev/bufio#Reader.ReadString
		input, _ := reader.ReadString('\n')

		if input == "exit\n" {
			fmt.Println("Exit shell.")
			break
		}

		fmt.Println(input)
	}
}
