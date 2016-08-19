package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var inFile = flag.String("in", "", "input file (defaults to stdin)")
var outFile = flag.String("out", "", "output file (defaults to stdout)")
var replace = flag.Bool("replace", false, "update file inplace")

func main() {
	flag.BoolVar(replace, "i", false, "update file inplace")
	flag.Parse()

	args := flag.Args()

	if *outFile != "" && *replace {
		log.Fatal("Cannot use -out and -replace/-i ")
	}

	if *inFile == "" && len(args) == 1 {
		inFile = &args[0]
	}

	var err error
	input := os.Stdin
	output := os.Stdout

	if *inFile != "" {
		input, err = os.Open(*inFile)
		if err != nil {
			log.Fatalf("Could not open input file: %s", err)
		}
		defer input.Close()
	} else {
		info, err := os.Stdin.Stat()
		if err == nil {
			mode := info.Mode()
			if mode&os.ModeCharDevice != 0 {
				fmt.Fprintf(os.Stderr, "Reading from stdin\n")
			}
		}
	}

	inBody, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("Error reading input: %s", err)
	}

	if len(inBody) == 0 {
		return
	}

	var msg json.RawMessage
	err = json.Unmarshal(inBody, &msg)
	if err == io.EOF {
		return
	} else if err != nil {
		log.Fatalf("Error parsing json: %s", err)
	}

	if *replace {
		outFile = inFile
	}

	if *outFile != "" {
		output, err = os.Create(*outFile)
		if err != nil {
			log.Fatalf("Error creating output file: %s", err)
		}
		defer output.Close()
	}

	out, err := json.MarshalIndent(&msg, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling json: %s", err)
	}

	if _, err := output.Write(out); err != nil {
		log.Fatalf("Error writing result: %s", err)
	}

	output.Write([]byte("\n"))
}
