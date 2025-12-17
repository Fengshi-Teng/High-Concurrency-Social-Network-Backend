package main

import (
	"encoding/json"
	"os"
	"proj2/server"
)

var Usage = `Usage: twitter <number of consumers>
    <number of consumers> = the number of goroutines (i.e., consumers) to be part of the parallel version.`

func main() {
	args := os.Args[1:]

	config := server.Config{
		Encoder: json.NewEncoder(os.Stdout),
		Decoder: json.NewDecoder(os.Stdin),
	}

	// If no argument → sequential mode
	if len(args) == 0 {
		config.Mode = "s"
	} else {
		config.Mode = "p"
		// convert the first argument to int (no error checking needed)
		config.ConsumersCount = int(args[0][0] - '0')
	}

	// Run the server
	// fmt.Fprintln(os.Stderr, ">>> Starting server...")
	server.Run(config)
}
