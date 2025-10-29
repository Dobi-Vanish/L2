package ntp

import (
    "fmt"
    "time"

    "github.com/beevik/ntp"
)

//GetCurrentTime returns time by ntp.
func GetCurrentTime() (time.Time, error) {
    response, err := ntp.Query("pool.ntp.org")
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to query NTP server: %w", err)
    }

    currentTime := time.Now().Add(response.ClockOffset)
    
    return currentTime, nil
}