package unpack

import (
    "errors"
    "strconv"
    "unicode"
)

// ErrInvalidString for returning errors.
var ErrInvalidString = errors.New("invalid string")

// Unpack function to unpack strings.
func Unpack(s string) (string, error) {
    if s == "" {
        return "", nil
    }

    if isAllDigits(s) {
        return "", ErrInvalidString
    }

    var result []rune
    runes := []rune(s)
    length := len(runes)

    for i := 0; i < length; i++ {
        current := runes[i]

    if current == '\\' {
        if i+1 >= length {
            return "", ErrInvalidString
        } 
        escapedChar := runes[i+1]
        i++
        if i+1 < length && unicode.IsDigit(runes[i+1]) {
            count, digits := extractNumber(runes[i+1:])
            if count == 0 {
                return "", ErrInvalidString
            }
            for j := 0; j < count; j++ {
                result = append(result, escapedChar)
            }
            i += digits
        } else {
            result = append(result, escapedChar)
        }
    continue
}

        if !unicode.IsDigit(current) {
            if i+1 < length && unicode.IsDigit(runes[i+1]) {
                count, digits := extractNumber(runes[i+1:])
                if count == 0 {
                    return "", ErrInvalidString
                }
                for j := 0; j < count; j++ {
                    result = append(result, current)
                }
                i += digits
            } else {
                result = append(result, current)
            }
        } else {
            return "", ErrInvalidString
        }
    }

    return string(result), nil
}

func isAllDigits(s string) bool {
    for _, r := range s {
        if !unicode.IsDigit(r) {
            return false
        }
    }
    return s != ""
}

func extractNumber(runes []rune) (int, int) {
    var digits []rune
    
    for i, r := range runes {
        if unicode.IsDigit(r) {
            digits = append(digits, r)
        } else {
            break
        }
        if i >= 10 {
            break
        }
    }
    
    if len(digits) == 0 {
        return 1, 0
    }
    
    num, err := strconv.Atoi(string(digits))
    if err != nil {
        return 1, 0
    }
    
    return num, len(digits)
}