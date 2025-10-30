package grep

import (
    "bufio"
    "io"
    "strconv"
)

// Options used to hold search settings.
type Options struct {
    After      int
    Before     int
    Count      bool
    IgnoreCase bool
    Invert     bool
    Fixed      bool
    LineNum    bool
    Pattern    string
}

// Search used to search for the provided word in providedd file or string.
func Search(input io.Reader, opts *Options) ([]string, error) {
    lines, err := readLines(input)
    if err != nil {
        return nil, err
    }

    matcher, err := newMatcher(opts)
    if err != nil {
        return nil, err
    }

    if opts.Count {
        count := countMatches(lines, matcher, opts)
        return []string{strconv.Itoa(count)}, nil
    }

    return searchWithContext(lines, matcher, opts)
}

func readLines(r io.Reader) ([]string, error) {
    var lines []string
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func countMatches(lines []string, matcher func(string) bool, opts *Options) int {
    count := 0
    for _, line := range lines {
        if matcher(line) != opts.Invert {
            count++
        }
    }
    return count
}

func searchWithContext(lines []string, matcher func(string) bool, opts *Options) ([]string, error) {
    var result []string
    matched := make([]bool, len(lines))

    for i, line := range lines {
        matched[i] = matcher(line) != opts.Invert
    }

    printed := make([]bool, len(lines))
    for i := 0; i < len(lines); i++ {
        if matched[i] && !printed[i] {
            start := max(0, i-opts.Before)
            for j := start; j < i; j++ {
                if !printed[j] {
                    result = appendLine(result, lines[j], j+1, opts.LineNum)
                    printed[j] = true
                }
            }

            result = appendLine(result, lines[i], i+1, opts.LineNum)
            printed[i] = true

            end := min(len(lines)-1, i+opts.After)
            for j := i + 1; j <= end; j++ {
                if !printed[j] {
                    result = appendLine(result, lines[j], j+1, opts.LineNum)
                    printed[j] = true
                }
            }

            if opts.After > 0 || opts.Before > 0 {
                result = append(result, "--")
            }
        }
    }

    if len(result) > 0 && result[len(result)-1] == "--" {
        result = result[:len(result)-1]
    }
    
    return result, nil
}

func appendLine(result []string, line string, lineNum int, showLineNum bool) []string {
    if showLineNum {
        return append(result, strconv.Itoa(lineNum)+":"+line)
    }
    return append(result, line)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}