package main

import (
    "flag"
    "fmt"
    "io"
    "os"

    "sortTask/sorttask"
)

func main() {
    var (
        column               int
        numeric              bool
        reverse              bool
        unique               bool
        checkSorted          bool
    )

    flag.IntVar(&column, "k", 0, "sort by column number (1-based)")
    flag.BoolVar(&numeric, "n", false, "sort by numeric value")
    flag.BoolVar(&reverse, "r", false, "sort in reverse order")
    flag.BoolVar(&unique, "u", false, "output only unique lines")
    flag.BoolVar(&checkSorted, "c", false, "check if data is sorted")
    
    flag.Parse()

    opts := &sorttask.Options{
        Column:                column,
        Numeric:               numeric,
        Reverse:               reverse,
        Unique:                unique,
        CheckSorted:           checkSorted,
    }

    var input io.Reader
    if flag.NArg() > 0 {
        file, err := os.Open(flag.Arg(0))
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
            os.Exit(1)
        }
        defer file.Close()
        input = file
    } else {
        input = os.Stdin
    }

    err := sorttask.Sort(input, os.Stdout, opts)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}