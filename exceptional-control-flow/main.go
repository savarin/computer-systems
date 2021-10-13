// TODO:
// - handle eof terminal character

package main

import (
	// "bufio"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
	// "os"
)

var wg sync.WaitGroup

func main() {
	// Rewrote as nested routines.
	//   https://stackoverflow.com/questions/56224836/stop-child-goroutine-when-parent-returns
	fmt.Printf("%d - main start\n", runtime.NumGoroutine())
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go parent(ctx)
	wg.Wait()
	cancel()

	// Stop main from immediate exit.
	time.Sleep(time.Second)
	fmt.Printf("%d - parent complete\n", runtime.NumGoroutine())
}

func parent(ctx context.Context) {
	fmt.Printf("%d - parent start\n", runtime.NumGoroutine())
	c := make(chan int)

	go child(ctx, c)

    select {
    case <-ctx.Done():
        fmt.Printf("%d - parent context expired\n", runtime.NumGoroutine())
    case <-c:
        fmt.Printf("%d - child complete\n", runtime.NumGoroutine())
    case <-time.After(time.Duration(3) * time.Second):
        fmt.Printf("%d - child timeout\n", runtime.NumGoroutine())
    }
    wg.Done()
}

func child(ctx context.Context, c chan int) {
	fmt.Printf("%d - child start\n", runtime.NumGoroutine())
    select {
    case <-ctx.Done():
        fmt.Printf("%d - child context expired\n", runtime.NumGoroutine())
    // Long-running task as a sleep.
    case <-time.After(time.Duration(1) * time.Second):
        fmt.Printf("%d - child done\n", runtime.NumGoroutine())
    }

    fmt.Printf("%d - child end\n", runtime.NumGoroutine())
    c <- 1
}


// 	// NewReader returns a new Reader whose buffer has the default size.
// 	//   https://pkg.go.dev/bufio#NewReader
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("▶️  ")

// 		// ReadString reads until the first occurrence of delim in the input, returning a string
// 		// containing the data up to and including the delimiter.
// 		//   https://pkg.go.dev/bufio#Reader.ReadString
// 		input, _ := reader.ReadString('\n')

// 		if input == "exit\n" {
// 			fmt.Println("Exit shell.")
// 			break
// 		}

// 		fmt.Println(input)
// 	}
// }
