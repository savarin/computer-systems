// TODO:
// - handle eof terminal character

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func log(msg string) {
	t := time.Now()
	ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Printf("%s thread %d - %s\n", ts, runtime.NumGoroutine(), msg)
}

func main() {
	log("main start")
	// Rewrote as nested routines.
	//   https://stackoverflow.com/questions/56224836/stop-child-goroutine-when-parent-returns
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go parent(ctx)
	wg.Wait()
	cancel()

	// Stop main from immediate exit.
	time.Sleep(time.Second)
	log("parent complete")
	log("main end")
}

func parent(ctx context.Context) {
	log("parent start")
	c := make(chan int)

	// NewReader returns a new Reader whose buffer has the default size.
	//   https://pkg.go.dev/bufio#NewReader
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("▶️  ")

		// ReadString reads until the first occurrence of delim in the input, returning a string
		// containing the data up to and including the delimiter.
		//   https://pkg.go.dev/bufio#Reader.ReadString
		cmd, _ := reader.ReadString('\n')
		cmd = cmd[:len(cmd)-1]

		if cmd == "exit" {
			fmt.Println("Exit shell.")
			break
		}

		if cmd == "pwd" || cmd == "ls" {
			go child(ctx, cmd, c)

			select {
			case <-ctx.Done():
				log("parent content expired")
			case <-c:
				log("child complete")
			case <-time.After(time.Duration(3) * time.Second):
				log("child timeout")
			}

			continue
		}

		fmt.Printf("%s: command not found\n", cmd)
	}

	log("parent end")
	wg.Done()
}

func child(ctx context.Context, cmd string, c chan int) {
	log("child start")

	switch cmd {
	case "pwd":
		fmt.Println("present working directory")
	case "ls":
		fmt.Println("list directory contents")
	}

	log("child end")
	c <- 1
}
