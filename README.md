# Shared Memory Data Concurrency

This project demonstrates concurrent data processing using shared memory in the Go programming language. The main goal is to process student records concurrently using multiple goroutines, filter out students based on a performance threshold, and write the sorted results to a file.

## Features

- Concurrent data processing using goroutines and channels.
- Threshold-based filtering: only students with calculated scores above 50.0 are retained.
- Custom sorting of the resulting student list based on names.
- File I/O operations to read input data and write filtered and sorted results.

## File Structure

- `main.go`: The main application file that sets up concurrency, processes data, and writes the output.
- Input data is expected from: `dat.txt`
- Output written to: `rez.txt`

## How it Works

1. The `main` function reads student data from a file.
2. Data is fed into multiple goroutines via channels for processing.
3. Each worker checks if the student's grade is above threshold (50.0 after computation).
4. Valid records are sent to a result processor which maintains a sorted list.
5. The sorted list is written to an output file.

## Technologies Used

- Go (Golang)
- Goroutines and Channels (Concurrency Primitives)

## How to Run

1. Ensure Go is installed on your system.
2. Place the input file (`dat.txt`) in the same directory.
3. Run the application:

```bash
go run main.go
```

4. Check the output in `rez.txt`.
