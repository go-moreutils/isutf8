package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

var (
	quiet = flag.Bool("q", false, "Don't print messages telling which files are invalid UTF-8, merely indicate it with the exit status.")
)

func check(name string, br *bufio.Reader) {
	l, c, o := 1, 1, 0
	for {
		r, size, err := br.ReadRune()
		o += size
		if r == 0xfffd || !utf8.ValidRune(r) {
			err = errors.New("invalid UTF-8 code")
		}
		if err != nil {
			if err != io.EOF {
				if !*quiet {
					fmt.Fprintf(os.Stderr, "%s: line %d, char %d, byte offset %d: %s\n", name, l, c, o, err)
				}
				os.Exit(1)
			}
			return
		}
		if r == '\n' {
			l++
			c = 1
		}
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		check("stdin", bufio.NewReader(os.Stdin))
		os.Exit(0)
	}
	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n", os.Args[0], arg, err)
			os.Exit(1)
		}
		check(arg, bufio.NewReader(f))
		f.Close()
	}
}
