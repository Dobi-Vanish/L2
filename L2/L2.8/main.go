package main

import (
    "fmt"
    "os"
    "ntp/ntp"
)

func main() {
    time, err := ntp.GetCurrentTime()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Current time by NTP: %s\n", time)
}