package sorttask

import (
    "bytes"
    "strings"
    "testing"
)

func TestSort(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        opts     *Options
    }{
        {
            name:     "basic sort",
            input:    "c\nb\na\n",
            expected: "a\nb\nc\n",
            opts:     &Options{},
        },
        {
            name:     "numeric sort",
            input:    "10\n2\n1\n",
            expected: "1\n2\n10\n",
            opts:     &Options{Numeric: true},
        },
        {
            name:     "reverse sort",
            input:    "a\nb\nc\n",
            expected: "c\nb\na\n",
            opts:     &Options{Reverse: true},
        },
        {
            name:     "unique sort",
            input:    "a\nb\na\nc\nb\n",
            expected: "a\nb\nc\n",
            opts:     &Options{Unique: true},
        },
        {
            name:     "column sort",
            input:    "1\tc\n2\ta\n3\tb\n",
            expected: "2\ta\n3\tb\n1\tc\n",
            opts:     &Options{Column: 2},
        },
        {
            name:     "numeric reverse sort",
            input:    "1\n3\n2\n",
            expected: "3\n2\n1\n",
            opts:     &Options{Numeric: true, Reverse: true},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            input := strings.NewReader(tt.input)
            var output bytes.Buffer
            
            err := Sort(input, &output, tt.opts)
            if err != nil {
                t.Fatalf("Sort failed: %v", err)
            }
            
            if output.String() != tt.expected {
                t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, output.String())
            }
        })
    }
}

func TestCheckSorted(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        isSorted    bool
        shouldError bool
        opts        *Options
    }{
        {
            name:        "sorted data",
            input:       "a\nb\nc\n",
            isSorted:    true,
            shouldError: false,
            opts:        &Options{CheckSorted: true},
        },
        {
            name:        "unsorted data",
            input:       "c\nb\na\n",
            isSorted:    false,
            shouldError: true,
            opts:        &Options{CheckSorted: true},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            input := strings.NewReader(tt.input)
            var output bytes.Buffer
            
            err := Sort(input, &output, tt.opts)
            
            if tt.shouldError && err == nil {
                t.Error("Expected error but got none")
            }
            if !tt.shouldError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
        })
    }
}