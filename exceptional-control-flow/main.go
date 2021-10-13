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
	"os/signal"
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

	// WithCancel returns a copy of parent with a new Done channel.
	//   https://pkg.go.dev/context#WithCancel
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go parent(ctx, cancel)
	wg.Wait()

	// Stop main from immediate exit.
	time.Sleep(time.Second)
	status("parent complete")
	status("main end")
}

func parent(ctx context.Context, cancel context.CancelFunc) {
	status("parent start")
	fmt.Println("✨ enter shell ✨")

	c := make(chan int)

	// Set up channel to catch interrupt signal.
	//   https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	status("signal start")
	defer func() {
		signal.Stop(s)
		status("signal end")
	}()

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
			case <-s:
				cancel()
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
		// Getwd returns a rooted path name corresponding to the current directory.
		//   https://pkg.go.dev/os#Getwd
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)

	case "echo":
		fmt.Println(strings.Join(args, " "))

	case "ls":
		path := "."
		if len(args) != 0 {
			path = args[0]
		}

		// ReadDir reads the directory named by dirname and returns a list of fs.FileInfo for the
		// directory's contents, sorted by filename
		//   https://pkg.go.dev/io/ioutil#ReadDir
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Printf("%s: please specify a valid path\n", args[0])
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}

	case "sleep":
		// ParseInt interprets a string s in the given base (0, 2 to 36) and bit size (0 to 64) and
		// returns the corresponding value i.
		//   https://pkg.go.dev/strconv#ParseInt
		duration, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("%s: please specify an integer\n", args[0])
		}
		time.Sleep(time.Duration(duration) * time.Second)
		fmt.Println("SLEEP COMPLETE")
	}

	status("child end")
	c <- 1
}
