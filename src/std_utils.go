// Responsible for some additional std functions.
package main

import (
	"log"
	"strings"
)

// Same as Crystal's String#rjust.
// It adjusts a string from the left side.
// ("1", 3, "0") => "001"
// ("test", 3, "0") => "test" (because test > 3)
func rjust(str string, char_no int, filling string) string {
	var safe_char_no int = char_no

	if char_no > len(str) {
		safe_char_no = char_no - len(str)
	}

	return strings.Repeat(filling, safe_char_no) + str
}

// It checks if err is not nil
// and log.Fatals it.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// [Unused] Checks if item exists in array.
// (["1", "2", "3"], "1") => true
// (["1", "2", "3"], "5") => false
// func includes(arr []string, item string) bool {
// 	for _, arr_item := range arr {
// 		if arr_item == item {
// 			return true
// 		}
// 	}
// 	return false
// }
