package main

import (
    "flag"
    "fmt"
    "io"
    "os"

    "grep/grep"
)

func main() {
    var (
        after      int
        before     int
        context    int
        count      bool
        ignoreCase bool
        invert     bool
        fixed      bool
        lineNum    bool
    )

    flag.IntVar(&after, "A", 0, "print N lines after match")
    flag.IntVar(&before, "B", 0, "print N lines before match")
    flag.IntVar(&context, "C", 0, "print N lines around match")
    flag.BoolVar(&count, "c", false, "print only count of matching lines")
    flag.BoolVar(&ignoreCase, "i", false, "ignore case")
    flag.BoolVar(&invert, "v", false, "select non-matching lines")
    flag.BoolVar(&fixed, "F", false, "interpret pattern as fixed string")
    flag.BoolVar(&lineNum, "n", false, "print line numbers")

    flag.Parse()

    if context > 0 {
        after = context
        before = context
    }

    if flag.NArg() == 0 {
        fmt.Fprintln(os.Stderr, "Error: pattern is required")
        os.Exit(1)
    }

    pattern := flag.Arg(0)

    opts := &grep.Options{
        After:      after,
        Before:     before,
        Count:      count,
        IgnoreCase: ignoreCase,
        Invert:     invert,
        Fixed:      fixed,
        LineNum:    lineNum,
        Pattern:    pattern,
    }

    var input io.Reader
    if flag.NArg() > 1 {
        filename := flag.Arg(1)
        file, err := os.Open(filename)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
            os.Exit(1)
        }
        defer file.Close()
        input = file
    } else {
        input = os.Stdin
    }

    result, err := grep.Search(input, opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    for _, line := range result {
        fmt.Println(line)
    }
}