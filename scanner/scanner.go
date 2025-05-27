package scanner

import (
	"fmt"
	"log/slog"
	"net"
	"sort"
	"sync"
)

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
