package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"9fans.net/go/acme"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("%s winid\n", os.Args[0])
	}

	winid, err := strconv.Atoi(os.Args[1])
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

	fmt.Printf("%s:#%d,#%d\n", filename, q0, q1)
}
