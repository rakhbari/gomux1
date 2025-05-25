package utils

import (
    "fmt"
    "strings"
)

var romanValuesMap = map[rune]int{
    'I': 1,
    'V': 5,
    'X': 10,
    'L': 50,
    'C': 100,
    'D': 500,
    'M': 1000,
}

var subtractersMap = map[rune]string{
    'I': "VX",
    'X': "LC",
    'C': "DM",
}

func RomanToInt(roman string) (int, error) {
    var subtractor rune
    numTotal := 0

    // Validate the length of the passed-in string is between 1 & 15
    if len(roman) < 1 || len(roman) > 15 {
        return -1, fmt.Errorf("ERROR: Given roman numeral string must be between 0 and 15")
    }

    for pos, char := range roman {
        // Validate the char is in romanValuesMap
        if _, exists := romanValuesMap[char]; !exists {
            return -1, fmt.Errorf("ERROR: Given roman numeral string must only include allowed roman numeral characters. Offender: %c", char)
        }

        if subtractor != 0 {
            // Previous char was a subtractor. Reduce the value of char by it and loop.
            numTotal += romanValuesMap[char] - romanValuesMap[subtractor]
            subtractor = 0
            continue
        }

        // Is there another rune after char and is char a subtractor for that rune?
        if pos < len(roman)-1 && strings.ContainsRune(subtractersMap[char], rune(roman[pos+1])) {
            subtractor = char
            continue
        }

        // char is not a subtractor. Add it to total and loop.
        numTotal += romanValuesMap[char]
    }

    return numTotal, nil
}
