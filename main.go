package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("logfile.log")
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	errorLines := make(chan string)
	normalLines := make(chan string)
	done := make(chan bool)

	go writeNormalLines(normalLines, done)
	go writeErrorLines(errorLines, done)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if containsError(line) {
			errorLines <- line
		} else {
			normalLines <- line
		}
	}
	close(errorLines)
	close(normalLines)

	<-done
	<-done

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning log file: %v", err)
	}
}

func writeErrorLines(lines <-chan string, done chan<- bool) {
	errorFile, err := os.Create("errorFile.log")
	if err != nil {
		log.Fatalf("Error creating error log file: %v", err)
	}
	defer errorFile.Close()

	for line := range lines {
		if _, err := errorFile.WriteString(line + "\n"); err != nil {
			log.Fatalf("Error writing to error log file: %v", err)
		}
	}
	done <- true
}


func writeNormalLines(lines <-chan string, done chan<- bool) {
	normalFile, err := os.Create("normalFile.log")
	if err != nil {
		log.Fatalf("Error creating normal log file: %v", err)
	}
	defer normalFile.Close()

	for line := range lines {
		if _, err := normalFile.WriteString(line + "\n"); err != nil {
			log.Fatalf("Error writing to normal log file: %v", err)
		}
	}
	done <- true
}

func containsError(line string) bool {
	return strings.Contains(line, "error")
}
