package sorttask

import (
    "bufio"
    "io"
    "sort"
    "strings"
	"strconv"
	"fmt"
)

// Options Struct holds options by which one to sort.
type Options struct {
    Column                int
    Numeric               bool
    Reverse               bool
    Unique                bool
    CheckSorted           bool
}

// Sort Sorts .txt file by provided option/flag.
func Sort(input io.Reader, output io.Writer, opts *Options) error {
    lines, err := readLines(input)
    if err != nil {
        return err
    }

    if opts.CheckSorted {
        if isSorted(lines, opts) {
            return nil
        }
        return fmt.Errorf("data is not sorted")
    }

    comparator := createComparator(opts)

    sort.Slice(lines, func(i, j int) bool {
        if opts.Reverse {
            return comparator(lines[j], lines[i])
        }
        return comparator(lines[i], lines[j])
    })

    if opts.Unique {
        lines = removeDuplicates(lines)
    }

    return writeLines(output, lines)
}

func readLines(r io.Reader) ([]string, error) {
    var lines []string
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func writeLines(w io.Writer, lines []string) error {
    for _, line := range lines {
        if _, err := fmt.Fprintln(w, line); err != nil {
            return err
        }
    }
    return nil
}

func isSorted(lines []string, opts *Options) bool {
    comparator := createComparator(opts)
    for i := 1; i < len(lines); i++ {
        if !comparator(lines[i-1], lines[i]) {
            return false
        }
    }
    return true
}

func removeDuplicates(lines []string) []string {
    if len(lines) == 0 {
        return lines
    }
    
    result := make([]string, 0, len(lines))
    result = append(result, lines[0])
    
    for i := 1; i < len(lines); i++ {
        if lines[i] != lines[i-1] {
            result = append(result, lines[i])
        }
    }
    
    return result
}

func createComparator(opts *Options) func(a, b string) bool {
    return func(a, b string) bool {

        if opts.Column > 0 {
            aCol := getColumn(a, opts.Column)
            bCol := getColumn(b, opts.Column)
            if aCol == "" {
                aCol = a
            }
            if bCol == "" {
                bCol = b
            }
            return compareValues(aCol, bCol, opts)
        }

        return compareValues(a, b, opts)
    }
}

func getColumn(line string, column int) string {
    columns := strings.Fields(line) 
    if column > 0 && column <= len(columns) {
        return columns[column-1]
    }
    return ""
}

func compareValues(a, b string, opts *Options) bool {
    switch {
    case opts.Numeric:
        return compareNumeric(a, b)
    default:
        return a < b
    }
}

func compareNumeric(a, b string) bool {
    numA, errA := strconv.ParseFloat(a, 64)
    numB, errB := strconv.ParseFloat(b, 64)

    if errA == nil && errB == nil {
        return numA < numB
    }

    if errA != nil && errB != nil {
        return a < b
    }

    return errA == nil
}