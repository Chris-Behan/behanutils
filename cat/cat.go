package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	// First element in args is the executable itself, which we don't want to include.
	args := os.Args[1:]
	options, files := parseOptionsAndFiles(args)
	for _, f := range files {
		if !fileExists(f) {
			fmt.Printf("File '%v' does not exist", f)
			return
		}
	}

	for _, o := range options {
		if !isValidOption(o) {
			fmt.Printf("%v is not a valid option.", o)
			return
		}
		if o == "-h" || o == "--help" {
			printHelpMessage()
			return
		}
	}

	for _, f := range files {
		printFileContents(f, options)
	}
}

func printHelpMessage() {
	helpMessage := `Prints files to stdout.
Usage: cat [OPTION].. [FILE].. 

Options:
-n, --number                Prefix each line in the output with its line number.
-b, --number-nonblank       Prefix each nonempty line in the output with its line number. Overrides -n
`
	fmt.Println(helpMessage)
}

func validOptions() []string {
	return []string{
		"-b",
		"--number-nonblank",
		"-n",
		"--number",
		"-h",
		"--help",
	}
}

func parseOptionsAndFiles(args []string) ([]string, []string) {
	options := []string{}
	files := []string{}
	for _, arg := range args {
		if rune(arg[0]) == '-' {
			options = append(options, arg)
		} else {
			files = append(files, arg)
		}
	}
	return options, files
}

func validateOptions(options []string) error {
	return nil
}

func validateFiles(files []string) error {
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isValidOption(option string) bool {
	for _, validOption := range validOptions() {
		if option == validOption {
			return true
		}
	}
	return false
}

func printFileContents(path string, options []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// defer schedules the code to be executed after the all other code in the surrounding
	// braces has finished executing. In this case it would be after this function finishes.
	// It is done in case there is an error encountered in the following code that prevents
	// us from closing the file properly. defer in this context is like the try-except-finally
	// block in python.
	defer file.Close()

	numberOption := false
	numberNonBlankOption := false
	for _, o := range options {
		if o == "-n" || o == "--number" {
			numberOption = true
		} else if o == "-b" || o == "--number-nonblank" {
			numberNonBlankOption = true
		}
	}
	// Create in memory buffered writer that we will use to build up our output, then eventually
	// flush to standard out. Saves us having to perform an IO operation to standard out on each line.
	writer := bufio.NewWriter(os.Stdout)
	// Default split function for the scanner is ReadLine, meaning each call to scan will return
	// the the bytes that compose the next line of text.
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		// number non blank option overrides number option so have this conditional first
		if numberNonBlankOption {
			if line != "" {
				line = strconv.Itoa(lineNumber) + " " + line
				// Only increment the line number for non blank lines if number non blank option is selected
				lineNumber++
			}
		} else if numberOption {
			line = strconv.Itoa(lineNumber) + " " + line
			lineNumber++
		}
		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Writes buffer data to stdout
	writer.Flush()
}
