package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/wzshiming/befunge"
	"github.com/wzshiming/cursor"
)

var debug = flag.Bool("d", false, "debug")
var interval = flag.Duration("i", time.Second/100, "debug interval")

func init() {
	flag.Usage = func() {
		usage := os.Args[0]
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [option] [file]:\n", usage)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	for _, path := range args {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		var input io.Reader = os.Stdin
		var output io.Writer = os.Stdout

		run := befunge.NewRunner(src)
		if *debug {
			buf := bytes.NewBuffer(nil)
			output = buf
			input = io.TeeReader(input, output)
			run.SetDebug(*debug)
			run.SetStep(func() {
				tmp := bytes.NewBuffer(nil)
				tmp.WriteString(cursor.RawClear())
				tmp.WriteString("\n=======Stack=======\n")
				for i, v := range run.Stack() {
					if i%5 == 0 {
						tmp.WriteString("\n")
					} else {
						tmp.WriteString(" ")
					}
					tmp.WriteString(fmt.Sprintf("%7s%4d,", strconv.QuoteRuneToASCII(rune(v)), v))
				}

				tmp.WriteString("\n=======Debug=======\n")
				tmp.WriteString(run.String())

				errs := run.Errors()
				if len(errs) != 0 {
					tmp.WriteString("\n=======Warning=======\n")
					for i, err := range errs {
						tmp.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
					}
				}
				tmp.WriteString("\n=======Output=======\n")
				tmp.WriteString(buf.String())
				out := tmp.String()
				time.Sleep(*interval)
				fmt.Print(out)
			})
		}
		run.SetOutput(output)
		run.SetInput(input)
		err = run.Run()
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			os.Exit(1)
			return
		}
		fmt.Println()
	}
}
