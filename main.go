package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type ScanHostOptions struct {
	Host    string
	Workers int
	To      int
	From    int
	Timeout time.Duration
}

func main() {
	var debugFlag = flag.Bool("debug", false, "debug")
	var host = flag.String("host", "google.com", "host")
	var from = flag.Int("from", 1, "from port")
	var to = flag.Int("to", 65535, "to port")
	var workers = flag.Int("workers", 10, "workers count")
	var _ = flag.Int("timeout", 200, "timeout in milliseconds")

	flag.Parse()

	logLevel := slog.LevelInfo
	if *debugFlag {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	fmt.Printf("Opened ports: %v\n", ScanHost(&ScanHostOptions{
		Host:    *host,
		Workers: *workers,
		To:      *to,
		From:    *from,
		Timeout: 200 * time.Millisecond,
	}))
}

func scanHostWorker(opt *ScanHostOptions, id int, portsCh chan int, resultsCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range portsCh {
		address := fmt.Sprintf("%s:%d", opt.Host, p)
		slog.Debug("scanning", "workerId", id, "port", p)
		conn, err := net.DialTimeout("tcp", address, opt.Timeout)
		if err == nil {
			slog.Debug("close", "workerId", id, "port", p)
			conn.Close()
			resultsCh <- p
		}
	}
}

func ScanHost(opt *ScanHostOptions) []int {
	portsCh := make(chan int, opt.Workers)
	resultsCh := make(chan int)

	var wg sync.WaitGroup

	for id := range opt.Workers {
		wg.Add(1)
		go scanHostWorker(opt, id, portsCh, resultsCh, &wg)
	}

	go func() {
		for i := opt.From; i <= opt.To; i++ {
			portsCh <- i
		}
		close(portsCh)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var openedPorts []int

	for p := range resultsCh {
		openedPorts = append(openedPorts, p)
	}

	sort.Ints(openedPorts)
	return openedPorts
}
