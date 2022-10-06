// Responsible for handling dist.torproject.org mirrors.
package main

import "strings"

type Mirror string

const (
	TOR_PROJECT Mirror = "https://dist.torproject.org/"
	EFF         Mirror = "https://tor.eff.org/dist/"
	CALYX       Mirror = "https://tor.calyxinstitute.org/dist/"
)

var mirrors = map[string]Mirror{
	"TOR_PROJECT": TOR_PROJECT,
	"EFF":         EFF,
	"CALYX":       CALYX,
}

// Returns a mirror from str.
// If it doesn't exist, it
// returns TOR_PROJECT.
func get_mirror(str string) Mirror {
	mirror := mirrors[strings.ToUpper(str)]

	if mirror == "" {
		return TOR_PROJECT
	}

	return mirror
}

// Returns a string of all
// mirrors seperated by a
// comma.
func get_mirrors() string {
	res := ""
	i := 0

	for k := range mirrors {
		if i != 0 {
			res = res + ", "
		}
		res = res + k

		i++
	}

	return res
}

// Uses a mirror on a url.
func use_mirror(mirror Mirror, url string) string {
	if mirror == TOR_PROJECT {
		return url
	}

	return strings.Replace(strings.ToLower(url), string(TOR_PROJECT), string(mirror), 1)
}
