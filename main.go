package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"9fans.net/go/acme"
)

var lineF = flag.Bool("l", false, "show in 'line col' format")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("%s [-l] winid\n", os.Args[0])
	}

	winid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	w, err := acme.Open(winid, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer w.CloseFiles()

	// get filename
	tag, err := w.ReadAll("tag")
	if err != nil {
		log.Fatalln(err)
	}

	i := bytes.IndexRune(tag, ' ')
	if i < 0 {
		i = len(tag)
	}
	filename := string(tag[:i])

	// get current dot position
	q0, q1 := 0, 0
	if pos0, err := strconv.Atoi(os.Getenv("acme_pos0")); err == nil {
		q0, q1 = pos0, pos0
	} else {
		_, _, err = w.ReadAddr()
		if err != nil {
			log.Fatalln(err)
		}
		err = w.Ctl("addr=dot")
		if err != nil {
			log.Fatalln(err)
		}
		q0, q1, err = w.ReadAddr()
		if err != nil {
			log.Fatalln(err)
		}
	}

	if *lineF {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanRunes)

		n := 0
		l := 0
		c := 0
		for scanner.Scan() {
			n++
			c++
			if scanner.Text() == "\n" {
				l++
				c = 0
			}
			if n == q0 {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%d %d\n", l, c)
	} else {
		fmt.Printf("%s:#%d,#%d\n", filename, q0, q1)
	}
}

