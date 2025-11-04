package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Config holds strings configs
type Config struct {
	fields    string
	delimiter string
	separated bool
}

func parseFields(fieldsStr string) (map[int]bool, error) {
	fields := make(map[int]bool)

	if fieldsStr == "" {
		return fields, nil
	}

	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("Invalid radius: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("Invalid field number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("Invalid field number: %s", rangeParts[1])
			}

			if start < 1 || end < 1 {
				return nil, fmt.Errorf("Field numbers must be >0: %s", part)
			}

			if start > end {
				return nil, fmt.Errorf("Radius can't start from more that the end: %s", part)
			}

			for i := start; i <= end; i++ {
				fields[i] = true
			}
		} else {
			fieldNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("Invalid field number: %s", part)
			}
			if fieldNum < 1 {
				return nil, fmt.Errorf("Field numbers must be >0: %d", fieldNum)
			}
			fields[fieldNum] = true
		}
	}

	return fields, nil
}

func processLine(line string, fields map[int]bool, delimiter string, separated bool) string {
	containsDelimiter := strings.Contains(line, delimiter)

	if separated && !containsDelimiter {
		return ""
	}

	if !containsDelimiter {
		return line
	}

	parts := strings.Split(line, delimiter)

	if len(fields) == 0 {
		return line
	}

	var result []string
	for i, part := range parts {
		fieldNum := i + 1
		if fields[fieldNum] {
			result = append(result, part)
		}
	}

	if len(result) == 0 {
		return ""
	}

	return strings.Join(result, delimiter)
}

func processInput(input io.Reader, config *Config) error {
	var fields map[int]bool
	var err error

	if config.fields != "" {
		fields, err = parseFields(config.fields)
		if err != nil {
			return err
		}
	} else {
		fields = make(map[int]bool)
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		processed := processLine(line, fields, config.delimiter, config.separated)
		if processed != "" {
			fmt.Println(processed)
		}
	}

	return scanner.Err()
}

func main() {
	config := &Config{}

	flag.StringVar(&config.fields, "f", "", "fields number for output")
	flag.StringVar(&config.delimiter, "d", "\t", "delimeter")
	flag.BoolVar(&config.separated, "s", false, "only strings with delimeter")

	flag.Parse()
	err := processInput(os.Stdin, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
