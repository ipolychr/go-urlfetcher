# Go URL Fetcher - Project Report

## Overview
This is a Go-based URL fetcher application that efficiently downloads content from multiple URLs concurrently using a worker pool pattern. The application is designed to process a list of URLs from a text file and generate a JSON report with the results.

## Project Structure
```
go-urlfetcher/
├── main.go              # Main application entry point
├── fetcher/
│   ├── fetcher.go       # Core worker pool implementation
│   └── fetcher_test.go  # Unit tests
├── urld.txt             # Sample URLs file
├── go.mod              # Go module definition
└── bin/
    └── urlfetcher      # Compiled binary
```

## Core Functionality

### 1. Concurrent URL Fetching
The application uses a worker pool pattern to fetch multiple URLs concurrently:
- **Configurable Workers**: Default 5 concurrent workers (configurable via `-workers` flag)
- **Timeout Handling**: 10-second timeout per HTTP request
- **Context Cancellation**: Supports graceful shutdown via SIGINT/SIGTERM signals

### 2. Input Processing
- **File-based Input**: Reads URLs from a text file (default: `urls.txt`)
- **Line-by-line Processing**: Each line in the file is treated as a separate URL
- **Error Handling**: Continues processing even if individual URLs fail

### 3. Output Generation
- **JSON Report**: Generates `results.json` with detailed results
- **Comprehensive Data**: Each result includes:
  - URL
  - HTTP status code
  - Response body length
  - Error messages (if any)

## Key Components

### Main Application (`main.go`)
- **Command-line Interface**: Accepts `-file` and `-workers` flags
- **Signal Handling**: Graceful shutdown on SIGINT/SIGTERM
- **Result Aggregation**: Collects all results and outputs to JSON

### Worker Pool (`fetcher/fetcher.go`)
- **Concurrent Processing**: Uses goroutines for parallel URL fetching
- **Context Awareness**: Respects cancellation signals
- **Error Resilience**: Continues processing even when individual requests fail
- **Resource Management**: Properly closes HTTP response bodies

### Result Structure
```go
type Result struct {
    URL    string `json:"url"`
    Status int    `json:"status"`
    Length int    `json:"length"`
    Err    string `json:"error,omitempty"`
}
```

## Usage Examples

### Basic Usage
```bash
# Run with default settings (5 workers, urls.txt file)
./bin/urlfetcher

# Specify custom file and worker count
./bin/urlfetcher -file urld.txt -workers 10
```

### Command-line Options
- `-file`: Path to file containing URLs (default: "urls.txt")
- `-workers`: Number of concurrent workers (default: 5)

## Sample Input/Output

### Input File (`urld.txt`)
```
https://example.com
https://httpbin.org/get
```

### Output (`results.json`)
```json
[
  {
    "url": "https://example.com",
    "status": 200,
    "length": 1256
  },
  {
    "url": "https://httpbin.org/get",
    "status": 200,
    "length": 292
  }
]
```

## Technical Features

### 1. Concurrency Model
- **Worker Pool Pattern**: Fixed number of goroutines processing jobs
- **Channel-based Communication**: Jobs and results passed via channels
- **WaitGroup Synchronization**: Ensures all workers complete before closing results

### 2. Error Handling
- **Network Errors**: Captures and reports HTTP client errors
- **Read Errors**: Handles body reading failures
- **Request Errors**: Manages malformed URLs and request creation failures

### 3. Resource Management
- **HTTP Client**: Reused across requests with 10-second timeout
- **Response Bodies**: Properly closed to prevent memory leaks
- **Context Cancellation**: Supports graceful shutdown

## Testing
The project includes unit tests in `fetcher_test.go`:
- **TestWorkerPoolSimple**: Tests basic worker pool functionality
- **HTTP Test Server**: Uses httptest for isolated testing
- **Status Verification**: Ensures correct HTTP status code handling

## Build and Run

### Compilation
```bash
go build -o bin/urlfetcher
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestWorkerPoolSimple ./fetcher
```

## Use Cases
1. **Website Monitoring**: Check status of multiple URLs
2. **Content Analysis**: Measure response sizes across different endpoints
3. **Load Testing**: Concurrent requests to test server capacity
4. **Data Collection**: Gather information from multiple sources
5. **Health Checks**: Verify availability of various services

## Performance Characteristics
- **Concurrent Processing**: Scales with number of workers
- **Memory Efficient**: Streams response bodies without storing full content
- **Timeout Protection**: Prevents hanging on slow/unresponsive URLs
- **Graceful Shutdown**: Handles interruption signals properly

## Dependencies
- **Standard Library Only**: No external dependencies
- **Go 1.25.1**: Requires Go 1.25.1 or later
- **Cross-platform**: Works on Windows, Linux, and macOS

This project demonstrates Go's strengths in concurrent programming and provides a solid foundation for URL processing tasks that require high throughput and reliability.
