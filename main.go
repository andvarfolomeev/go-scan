package main

import (
	"andvarfolomeev/go-scan/scanner"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

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

	fmt.Printf("Opened ports: %v\n", scanner.ScanHost(&scanner.ScanHostOptions{
		Host:    *host,
		Workers: *workers,
		To:      *to,
		From:    *from,
		Timeout: 200 * time.Millisecond,
	}))
}
