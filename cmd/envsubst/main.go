// envsubst command line tool
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/a8m/envsubst"
	"io"
	"os"
)

var (
	input  = flag.String("i", "", "")
	output = flag.String("o", "", "")
)

var usage = `Usage: envsubst [options...] <input>
Options:
  -i  Specify file input, otherwise use last argument as input file. 
      If no input file is specified, read from stdin.
  -o  Specify file output. If none is specified, write to stdout.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	// Reader
	var reader *bufio.Reader
	if *input != "" {
		file, err := os.Open(*input)
		if err != nil {
			usageAndExit(fmt.Sprintf("Error to open file input: %s.", *input))
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	} else {
		if stat, err := os.Stdin.Stat(); err == nil && stat.Size() > 0 {
			reader = bufio.NewReader(os.Stdin)
		} else {
			usageAndExit("")
		}
	}
	// Collect data
	var data string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			usageAndExit("Failed to read input.")
		}
		data += line
	}
	// Writer
	var file *os.File
	var err error
	if *output != "" {
		file, err = os.Create(*output)
		if err != nil {
			usageAndExit("Error to create the wanted output file.")
		}
	} else {
		file = os.Stdout
	}
	// Parse input string
	result, err := envsubst.String(data)
	if err != nil {
		usageAndExit(fmt.Sprintf("envsubst error: %s.", err))
	}
	if _, err := file.WriteString(result); err != nil {
		filename := *output
		if filename == "" {
			filename = "STDOUT"
		}
		usageAndExit(fmt.Sprintf("Error writing output to: %s.", filename))
	}
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
