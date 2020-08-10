package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var inFile = flag.String("in", "", "input file (defaults to stdin)")
var outFile = flag.String("out", "", "output file (defaults to stdout)")
var ugly = flag.Bool("ugly", false, "format compactly")
var replace = flag.Bool("replace", false, "update file inplace")
var stream = flag.Bool("stream", false, "read streaming input")

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

	if *replace {
		dir := filepath.Dir(*inFile)
		output, err = ioutil.TempFile(dir, "ppjson")
		if err != nil {
			log.Fatalf("Error creating tmp output file: %s", err)
		}
	}

	if *outFile != "" {
		output, err = os.Create(*outFile)
		if err != nil {
			log.Fatalf("Error creating output file: %s", err)
		}
	}

	if *stream {
		streamDecode(input, output)
	} else {
		singleDecode(input, output)
	}

	output.Close()
	if *replace {
		err = os.Rename(output.Name(), *inFile)
		if err != nil {
			log.Fatalf("Error renaming tmp file %s to %s: %s", output.Name(), *inFile, err)
		}
	}
}

func streamDecode(input io.Reader, output io.Writer) {
	dec := json.NewDecoder(input)
	enc := json.NewEncoder(output)
	if !*ugly {
		enc.SetIndent("", "  ")
	}
	var msg json.RawMessage
	for {
		err := dec.Decode(&msg)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error reading input: %s", err)
		}

		err = enc.Encode(msg)
		if err != nil {
			log.Fatalf("Error encoding message: %s", err)
		}
	}
}

func singleDecode(input io.Reader, output io.Writer) {
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

	var out []byte
	if *ugly {
		out, err = json.Marshal(&msg)
	} else {
		out, err = json.MarshalIndent(&msg, "", "  ")
	}
	if err != nil {
		log.Fatalf("Error marshaling json: %s", err)
	}

	if _, err := output.Write(out); err != nil {
		log.Fatalf("Error writing result: %s", err)
	}

	output.Write([]byte("\n"))
}
