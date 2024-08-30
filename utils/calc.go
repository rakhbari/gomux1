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

/*
Iterate over runes in the roman numeral string
Do we have a lastChar -> No? add it to the total
Yes?
*/
func RomanToInt(roman string) (int, error) {
    var subtractor rune
    var subtractors string
    numTotal := 0

    // Validation the length of the pass-in string
    if len(roman) < 1 || len(roman) > 15 {
        return -1, fmt.Errorf("ERROR: Given roman numeral string must be between 0 and 15")
    }

    for pos, char := range roman {
        // Validate the char is in romanValuesMap
        if _, exists := romanValuesMap[char]; !exists {
            return -1, fmt.Errorf("ERROR: Given roman numeral string must only include allowed roman numeral characters. Offender: %c", char)
        }

        if subtractor != 0 {
            numTotal += romanValuesMap[char] - romanValuesMap[subtractor]
            subtractor = 0
        } else {
            subtractors = subtractersMap[char]
            // char is not a subtractor
            if subtractors == "" {
                numTotal += romanValuesMap[char]
            } else {
                // Is there another rune after char and is char a subtractor for that rune?
                if pos < len(roman)-1 && strings.ContainsRune(subtractors, rune(roman[pos+1])) {
                    subtractor = char
                } else {
                    numTotal += romanValuesMap[char]
                    subtractor = 0
                }
            }
        }
    }

    return numTotal, nil
}
