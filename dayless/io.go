package dayless

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/**
Read a file and returning the lines as array (without newlines)
 */
func ReadFileToArray(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	r := bufio.NewReader(f)
	var lines []string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return nil, err
		}
		lines = append(lines, line[0:len(line)-1]) // remove newline (last char)
	}
	return lines, nil
}

/**
Read a file and returns its content as one string
 */
func ReadFileToString(path string) (*string, error) {
	if lines, err := ReadFileToArray(path); err != nil {
		return nil, err
	} else {
		result := strings.Join(lines, "\n")
		return &result, nil
	}
}
