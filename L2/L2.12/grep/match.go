package grep

import (
    "regexp"
    "strings"
)

func newMatcher(opts *Options) (func(string) bool, error) {
    pattern := opts.Pattern

    if opts.IgnoreCase {
        pattern = "(?i)" + pattern
    }

    if opts.Fixed {
        pattern = regexp.QuoteMeta(pattern)
    }

    re, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }

    return func(line string) bool {
        if opts.IgnoreCase && !opts.Fixed {
            return re.MatchString(strings.ToLower(line))
        }
        return re.MatchString(line)
    }, nil
}