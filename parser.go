package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type fileType struct {
	fileName    string
	topic       string // всегда одна?
	tag         []string
	link        string
	bindingFile []string
	data        time.Time
}

// перевести файлы на .csv
/*
todo
link
день недели убрать из даты
*/

var base []fileType

func fileRead(fileName string) {
	bs, err := os.ReadFile(fileName)

	if err != nil {
		return
	}
	str := string(bs)
	sl := strings.Split(str, "\n")

	for i, line := range sl {
		fmt.Printf("%d: %s\n", i, string(line))
	}
	///---------------

	var val fileType

	val.fileName = fileName

	for _, line := range sl {
		if strings.Contains(line, "topic:") {
			val.topic = strings.TrimPrefix(line, "topic: ")
		}

		if strings.Contains(line, "#") {
			slice := strings.Split(line, "#")
			for _, s := range slice {
				if s != "" && s != "\r" {
					s = strings.TrimSuffix(s, " ")
					val.tag = append(val.tag, s)
				}
			}
		}

		if strings.Contains(line, "link:") {
			s := strings.TrimSuffix(line, "\r")
			val.link = strings.TrimPrefix(s, "link: ")
		}

		if strings.Contains(line, "[") { // todo парс строк
			slice := strings.Split(line, "[")
			for _, s := range slice {
				s = strings.TrimSuffix(s, "\r")
				s = strings.TrimSuffix(s, "]")
				if s != "" {
					val.bindingFile = append(val.bindingFile, s)
				}
			}
		}

		if strings.Contains(line, "data:") {
			s := strings.TrimPrefix(line, "data: ")
			s = strings.TrimSuffix(s, "\r")
			d, err := time.Parse("2006.01.02 15:04", s)
			fmt.Println("v", d.Weekday())
			fmt.Println(d, err)
			val.data = d
		}
	}

	base = append(base, val)
}
