# go-scan

A fast, concurrent TCP port scanner written in Go.

## Features

- Scans TCP ports on a target host
- Configurable port range
- Adjustable concurrency level
- Timeout control for scan operations
- Graceful termination on interrupt signals

## Usage

```
./go-scan [options]
```

### Options

- `-host string`: Target hostname or IP address (default "google.com")
- `-from int`: Starting port number to scan, inclusive (default 1)
- `-to int`: Ending port number to scan, inclusive (default 65535)
- `-workers int`: Number of concurrent scanning workers (default 10)
- `-timeout int`: Timeout per port scan in milliseconds (default 200)

### Examples

Scan all ports on localhost:
```
./go-scan -host localhost
```

Scan a specific port range with higher concurrency:
```
./go-scan -host example.com -from 20 -to 1000 -workers 50
```

Scan with a longer timeout for slower connections:
```
./go-scan -host 192.168.1.1 -timeout 500
```

## Building

Clone the repository and build using Go:

```
git clone https://github.com/yourusername/go-scan.git
cd go-scan
go build
```

## License

MIT