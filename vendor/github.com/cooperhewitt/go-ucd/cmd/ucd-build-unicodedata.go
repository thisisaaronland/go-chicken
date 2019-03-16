package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {

	unicode_data := flag.String("data", "https://unicode.org/Public/UCD/latest/ucd/UnicodeData.txt", "The URL of your Unicode source data")

	flag.Parse()

	rsp, err := http.Get(*unicode_data)
	defer rsp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	delimiter, _ := utf8.DecodeRuneInString(";")

	reader := csv.NewReader(rsp.Body)
	reader.Comma = delimiter

	ts := time.Now()

	fmt.Printf("%s\n\n", "package unicodedata")

	fmt.Printf("/* %s */\n", *unicode_data)
	fmt.Printf("/* This file was generated by robots at %s */\n\n", ts.UTC())

	fmt.Printf("%s\n", "var UCD = map[string]string{")

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			continue
		}

		if strings.HasPrefix(record[1], "<") {
			fmt.Printf("\t\t/* \"%s\":\t\"%s\", */\n", record[0], record[1])
		} else if strings.HasSuffix(record[1], ">") {
			fmt.Printf("\t\t/* \"%s\":\t\"%s\", */\n", record[0], record[1])
		} else {
			fmt.Printf("\t\t\"%s\":\t\"%s\",\n", record[0], record[1])
		}
	}

	fmt.Printf("%s\n", "}")
}
