package unpack

import (
    "testing"
)

func TestUnpack(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        hasError bool
    }{
        {
            name:     "simple unpack",
            input:    "a4bc2d5e",
            expected: "aaaabccddddde",  
            hasError: false,
        },
        {
            name:     "no digits",
            input:    "abcd",
            expected: "abcd",
            hasError: false,
        },
        {
            name:     "only digits",
            input:    "45",
            expected: "",
            hasError: true,
        },
        {
            name:     "empty string",
            input:    "",
            expected: "",
            hasError: false,
        },
        {
            name:     "single digit after char",
            input:    "a1",
            expected: "a",
            hasError: false,
        },
        {
            name:     "multiple digits after char",
            input:    "a10",
            expected: "aaaaaaaaaa",
            hasError: false,
        },
        {
            name:     "escaped digits",
            input:    "qwe\\4\\5",
            expected: "qwe45",
            hasError: false,
        },
        {
            name:     "escaped digit with multiplier",
            input:    "qwe\\45",
            expected: "qwe44444",
            hasError: false,
        },
        {
            name:     "escaped backslash",
            input:    "qwe\\\\5",
            expected: "qwe\\\\\\\\\\",
            hasError: false,
        },
        {
            name:     "escaped character",
            input:    "qwe\\a3",
            expected: "qweaaa",
            hasError: false,
        },
        
        // Ошибочные случаи
        {
            name:     "digit at start",
            input:    "1abc",
            expected: "",
            hasError: true,
        },
        {
            name:     "multiple digits at start",
            input:    "12abc",
            expected: "",
            hasError: true,
        },
        {
            name:     "escape at end",
            input:    "abc\\",
            expected: "",
            hasError: true,
        },
        {
            name:     "zero repeat",
            input:    "a0",
            expected: "",
            hasError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Unpack(tt.input)
            
            if tt.hasError {
                if err == nil {
                    t.Errorf("Unpack(%q) expected error, but got none", tt.input)
                }
                if err != ErrInvalidString {
                    t.Errorf("Unpack(%q) expected ErrInvalidString, got %v", tt.input, err)
                }
            } else {
                if err != nil {
                    t.Errorf("Unpack(%q) unexpected error: %v", tt.input, err)
                }
                if result != tt.expected {
                    t.Errorf("Unpack(%q) = %q, want %q", tt.input, result, tt.expected)
                }
            }
        })
    }
}

// Бенчмарк для проверки производительности
func BenchmarkUnpack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Unpack("a4bc2d5e\\4\\5qwe\\45")
    }
}