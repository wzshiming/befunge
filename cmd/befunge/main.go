package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

		run := befunge.NewRunner(src)
		if *debug {
			buf := bytes.NewBuffer(nil)
			run.SetDebug(*debug)
			run.SetOutput(buf)
			run.SetInput(io.TeeReader(os.Stdin, buf))
			run.SetStep(func() {
				tmp := bytes.NewBuffer(nil)
				tmp.WriteString(cursor.RawClear())
				tmp.WriteString("\n=======Stack=======\n")
				tmp.WriteString(fmt.Sprint(run.Stack()))
				tmp.WriteString("\n=======Debug=======\n")
				tmp.WriteString(run.String())
				tmp.WriteString("\n=======Output=======\n")
				tmp.WriteString(buf.String())
				out := tmp.String()
				time.Sleep(*interval)
				fmt.Print(out)
			})
		}
		err = run.Run()
		if err != nil {
			fmt.Println()
			fmt.Println(err)
			return
		}
	}
	fmt.Println()
}
