# Timer

[![Go Reference](https://pkg.go.dev/badge/github.com/gmodx/timer.svg)](https://pkg.go.dev/github.com/gmodx/timer)

The `timer` package provides a simple way to schedule and execute functions at specified intervals.

## Installation

To use this package in your Go project, you can install it using `go get`:

```sh
go get github.com/gmodx/timer
```

## Usage
Here's a basic example of how to use the timer package to schedule and execute functions at specified intervals:

### Tick
The Tick function schedules a job function to be executed repeatedly with a specified delay and interval. It takes the following parameters:

* delay: Initial delay before starting the job execution.
* interval: Time interval between consecutive executions of the job function.
* jobFunc: The function to be executed as a job.
* jobErrCallback: A callback function that handles errors encountered during job execution.
* params: Optional parameters to be passed to the job function.


``` go
func main() {
    err := timer.Tick(
        2 * time.Second, // Initial delay of 2 seconds
        1 * time.Second, // Execute the job every 1 second
        myJobFunction,   // The job function to be executed
        errorHandler,    // Function to handle job execution errors
        "param1",        // Optional parameters for the job function
        42,
    )
    if err != nil {
        fmt.Println("Error:", err)
    }

    // Keep the program running to allow the scheduled jobs to execute.
    select {}
}

func myJobFunction(param1 string, param2 int) {
    // Your job logic here
}

func errorHandler(err error) {
    fmt.Println("Job error:", err)
}
```

## Documentation
For more information and usage examples, please refer to the GoDoc documentation.

