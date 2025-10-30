package main

import (
    "sort"
    "strings"
	"fmt"
)

func findAnagrams(words []string) map[string][]string {
    if len(words) == 0 {
        return make(map[string][]string)
    }

    anagramGroups := make(map[string][]string)
    
    for _, word := range words {
        lowerWord := strings.ToLower(word)
        
        key := sortString(lowerWord)
        
        anagramGroups[key] = append(anagramGroups[key], lowerWord)
    }

    result := make(map[string][]string)
    
    for _, group := range anagramGroups {
        if len(group) <= 1 {
            continue
        }
        
        uniqueGroup := removeDuplicates(group)
        sort.Strings(uniqueGroup)
        
        firstWord := uniqueGroup[0]
        result[firstWord] = uniqueGroup
    }

    return result
}

func sortString(s string) string {
    runes := []rune(s)
    
    sort.Slice(runes, func(i, j int) bool {
        return runes[i] < runes[j]
    })
    
    return string(runes)
}

func removeDuplicates(words []string) []string {
    if len(words) == 0 {
        return words
    }
    
    seen := make(map[string]bool)
    result := make([]string, 0, len(words))
    
    for _, word := range words {
        if !seen[word] {
            seen[word] = true
            result = append(result, word)
        }
    }
    
    return result
}

func main() {
    words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
    
    anagrams := FindAnagrams(words)
    
    for key, group := range anagrams {
        fmt.Printf("%s: %v\n", key, group)
    }
}