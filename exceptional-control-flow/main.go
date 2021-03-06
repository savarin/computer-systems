package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
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

var debug bool
var timeout = 60

func status(msg string) {
	if debug {
		t := time.Now()
		ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Printf("%s thread %d - %s\n", ts, runtime.NumGoroutine(), msg)
	}
}

func main() {
	// Bind cli argument to variable.
	//   https://pkg.go.dev/flag#BoolVar
	flag.BoolVar(&debug, "debug", false, "Set up debug mode.")
	flag.Parse()

	status("main start")
	defer func() {
		status("main end")
	}()

	// Re-implement basic loop as nested routines.
	//   https://stackoverflow.com/questions/56224836/stop-child-goroutine-when-parent-returns
	wg.Add(1)

	// Package context defines the Context type, which carries deadlines, cancellation signals, and
	// other request-scoped values across API boundaries and between processes.
	//   https://pkg.go.dev/context
	ctx := context.Background()
	// WithCancel returns a copy of parent with a new Done channel.
	//   https://pkg.go.dev/context#WithCancel
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	go parent(ctx, cancel)
	defer func() {
		status("parent complete")
	}()

	wg.Wait()
}

func parent(ctx context.Context, cancel context.CancelFunc) {
	status("parent start")
	defer func() {
		status("parent end")
		wg.Done()
	}()

	fmt.Println("✨ enter shell ✨")
	defer func() {
		fmt.Println("✨ exit shell ✨")
	}()

	// Set up channel for child goroutines.
	c := make(chan int)

	// Package signal implements access to incoming signals.
	//   https://pkg.go.dev/os/signal
	s := make(chan os.Signal, 1)
	// Set up channel to catch interrupt signals.
	//   https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
	signal.Notify(s, os.Interrupt)
	status("signal start")
	defer func() {
		signal.Stop(s)
		close(s)
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
		input, err := reader.ReadString('\n')

		// Allow exit via eof terminal control character or "exit".
		if err == io.EOF || input == "exit\n" {
			break
		}

		// Remove newline character and split string by space delimiter.
		args := strings.Split(input[:len(input)-1], " ")
		cmd, args := args[0], args[1:]

		if cmd == "pwd" || cmd == "ls" || cmd == "echo" || cmd == "sleep" {
			go child(ctx, cancel, cmd, args, c, s)

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
}

func child(ctx context.Context, cancel context.CancelFunc, cmd string, args []string, c chan int, s chan os.Signal) {
	status("child start")
	defer func() {
		status("child end")
		c <- 1
	}()

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

		fmt.Printf("sleep for %d seconds...\n", duration)

		select {
		case <-s:
			fmt.Println("sleep interrupted!")
			cancel()
		case <-time.After(time.Duration(duration) * time.Second):
			fmt.Println("sleep complete!")
		}
	}
}
