package scanner

import "time"

type ScanHostOptions struct {
	Host    string
	Workers int
	To      int
	From    int
	Timeout time.Duration
}
