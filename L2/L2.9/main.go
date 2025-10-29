package main

import (
    "fmt"
    "unpack/unpack"
)

func main() {
    testCases := []string{
        "a4bc2d5e",
        "abcd", 
        "45",
        "",
        "qwe\\4\\5",
        "qwe\\45",
        "a1b2c3",
    }

    for _, test := range testCases {
        result, err := unpack.Unpack(test)
        if err != nil {
            fmt.Printf("Input: %q -> Error: %v\n", test, err)
        } else {
            fmt.Printf("Input: %q -> Output: %q\n", test, result)
        }
    }
}