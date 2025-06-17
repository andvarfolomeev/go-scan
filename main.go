package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ScanOptions struct {
	Host    string
	Workers int
	To      int
	From    int
	Timeout time.Duration
}

func main() {
	var host = flag.String("host", "google.com", "target hostname or IP address")
	var from = flag.Int("from", 1, "starting port number to scan (inclusive)")
	var to = flag.Int("to", 65535, "ending port number to scan (inclusive)")
	var workers = flag.Int("workers", 10, "number of concurrent scanning workers")
	var timeout = flag.Int("timeout", 200, "timeout per port scan in milliseconds")

	flag.Parse()
	opts := &ScanOptions{*host, *workers, *to, *from, time.Duration(*timeout) * time.Millisecond}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nReceived interrupt. Shutting down...")
		cancel()
	}()

	portCh := generatePorts(ctx, *from, *to)
	openedPortCh := scanPorts(ctx, opts, portCh)

	for {
		select {
		case port, ok := <-openedPortCh:
			if !ok {
				return
			}
			fmt.Printf("%d: opened\n", port)
		case <-ctx.Done():
			return
		}
	}
}

func generatePorts(ctx context.Context, from, to int) chan int {
	outCh := make(chan int)

	go func() {
		defer close(outCh)

		for port := from; port < to; port++ {
			select {
			case outCh <- port:
			case <-ctx.Done():
				return
			}
		}
	}()

	return outCh
}

func scanPorts(ctx context.Context, opts *ScanOptions, inCh chan int) chan int {
	var wg sync.WaitGroup
	outCh := make(chan int)

	wg.Add(opts.Workers)
	for range opts.Workers {
		go scanPortsWorker(ctx, opts, inCh, outCh, &wg)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func scanPortsWorker(ctx context.Context, opts *ScanOptions, inCh chan int, outCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range inCh {
		address := fmt.Sprintf("%s:%d", opts.Host, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}

		_ = conn.Close()

		select {
		case outCh <- port:
		case <-ctx.Done():
			return
		}
	}
}
