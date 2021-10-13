// TODO:
// - handle eof terminal character

package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

var timeout = 60

func status(msg string) {
	t := time.Now()
	ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fmt.Printf("%s thread %d - %s\n", ts, runtime.NumGoroutine(), msg)
}

func main() {
	status("main start")
	// Rewrote as nested routines.
	//   https://stackoverflow.com/questions/56224836/stop-child-goroutine-when-parent-returns
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go parent(ctx)
	wg.Wait()
	cancel()

	// Stop main from immediate exit.
	time.Sleep(time.Second)
	status("parent complete")
	status("main end")
}

func parent(ctx context.Context) {
	status("parent start")
	c := make(chan int)

	// NewReader returns a new Reader whose buffer has the default size.
	//   https://pkg.go.dev/bufio#NewReader
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("▶️  ")

		// ReadString reads until the first occurrence of delim in the input, returning a string
		// containing the data up to and including the delimiter.
		//   https://pkg.go.dev/bufio#Reader.ReadString
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		args := strings.Split(input, " ")
		cmd, args := args[0], args[1:]

		if cmd == "exit" {
			fmt.Println("✨ exit shell ✨")
			break
		}

		if cmd == "pwd" || cmd == "ls" || cmd == "echo" || cmd == "sleep" {
			go child(ctx, cmd, args, c)

			select {
			case <-ctx.Done():
				status("parent content expired")
			case <-c:
				status("child complete")
			case <-time.After(time.Duration(timeout) * time.Second):
				status("child timeout")
			}

			continue
		}

		fmt.Printf("%s: command not found\n", cmd)
	}

	status("parent end")
	wg.Done()
}

func child(ctx context.Context, cmd string, args []string, c chan int) {
	status("child start")

	switch cmd {

	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)

	case "ls":
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}

	case "echo":
		fmt.Println(strings.Join(args, " "))

	case "sleep":
		duration, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(duration) * time.Second)
	}

	status("child end")
	c <- 1
}
