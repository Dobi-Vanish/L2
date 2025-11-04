package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[int]bool
		wantErr  bool
	}{
		{
			name:     "single field",
			input:    "1",
			expected: map[int]bool{1: true},
			wantErr:  false,
		},
		{
			name:     "multiple fields",
			input:    "1,3,5",
			expected: map[int]bool{1: true, 3: true, 5: true},
			wantErr:  false,
		},
		{
			name:     "range",
			input:    "1-3",
			expected: map[int]bool{1: true, 2: true, 3: true},
			wantErr:  false,
		},
		{
			name:     "mixed fields and ranges",
			input:    "1,3-5,7",
			expected: map[int]bool{1: true, 3: true, 4: true, 5: true, 7: true},
			wantErr:  false,
		},
		{
			name:     "with spaces",
			input:    "1, 3-5 , 7",
			expected: map[int]bool{1: true, 3: true, 4: true, 5: true, 7: true},
			wantErr:  false,
		},
		{
			name:    "invalid field",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid range",
			input:   "1-3-5",
			wantErr: true,
		},
		{
			name:    "negative field",
			input:   "-1",
			wantErr: true,
		},
		{
			name:    "zero field",
			input:   "0",
			wantErr: true,
		},
		{
			name:    "reverse range",
			input:   "5-3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFields(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFields(%q) wanted error, but got: %v", tt.input, result)
				}
				return
			}

			if err != nil {
				t.Errorf("parseFields(%q) error: %v", tt.input, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("parseFields(%q) = %v, expected %v", tt.input, result, tt.expected)
				return
			}

			for k := range tt.expected {
				if !result[k] {
					t.Errorf("parseFields(%q) missing field %d in %v", tt.input, k, result)
				}
			}
		})
	}
}

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		fields    map[int]bool
		delimiter string
		separated bool
		expected  string
	}{
		{
			name:      "single field",
			line:      "a\tb\tc",
			fields:    map[int]bool{1: true},
			delimiter: "\t",
			separated: false,
			expected:  "a",
		},
		{
			name:      "multiple fields",
			line:      "a\tb\tc\td",
			fields:    map[int]bool{1: true, 3: true},
			delimiter: "\t",
			separated: false,
			expected:  "a\tc",
		},
		{
			name:      "range fields",
			line:      "a\tb\tc\td\te",
			fields:    map[int]bool{2: true, 3: true, 4: true},
			delimiter: "\t",
			separated: false,
			expected:  "b\tc\td",
		},
		{
			name:      "comma delimiter",
			line:      "a,b,c",
			fields:    map[int]bool{2: true},
			delimiter: ",",
			separated: false,
			expected:  "b",
		},
		{
			name:      "field out of range",
			line:      "a\tb\tc",
			fields:    map[int]bool{5: true},
			delimiter: "\t",
			separated: false,
			expected:  "",
		},
		{
			name:      "separated flag with delimiter",
			line:      "a\tb\tc",
			fields:    map[int]bool{1: true},
			delimiter: "\t",
			separated: true,
			expected:  "a",
		},
		{
			name:      "separated flag without delimiter",
			line:      "a b c",
			fields:    map[int]bool{1: true},
			delimiter: "\t",
			separated: true,
			expected:  "",
		},
		{
			name:      "no fields specified",
			line:      "a\tb\tc",
			fields:    map[int]bool{},
			delimiter: "\t",
			separated: false,
			expected:  "a\tb\tc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processLine(tt.line, tt.fields, tt.delimiter, tt.separated)
			if result != tt.expected {
				t.Errorf("processLine(%q, %v, %q, %t) = %q, expected %q",
					tt.line, tt.fields, tt.delimiter, tt.separated, result, tt.expected)
			}
		})
	}
}

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   *Config
		expected string
		wantErr  bool
	}{
		{
			name:  "basic functionality",
			input: "a\tb\tc\nd\te\tf\n",
			config: &Config{
				fields:    "1,3",
				delimiter: "\t",
				separated: false,
			},
			expected: "a\tc\nd\tf\n",
			wantErr:  false,
		},
		{
			name:  "with separated flag",
			input: "a\tb\tc\nd e f\ng\th\ti\n",
			config: &Config{
				fields:    "1",
				delimiter: "\t",
				separated: true,
			},
			expected: "a\ng\n",
			wantErr:  false,
		},
		{
			name:  "comma delimiter",
			input: "a,b,c\nd,e,f\n",
			config: &Config{
				fields:    "2",
				delimiter: ",",
				separated: false,
			},
			expected: "b\ne\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := processInput(input, tt.config)

			w.Close()
			os.Stdout = oldStdout

			io.Copy(output, r)
			result := output.String()

			if tt.wantErr {
				if err == nil {
					t.Errorf("processInput() expected error, but didnt get one")
				}
				return
			}

			if err != nil {
				t.Errorf("processInput() error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("processInput() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
