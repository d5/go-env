package env

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var (
	exportRE = regexp.MustCompile(`^\s*(export\s+)?`)
	keyRE    = regexp.MustCompile(`^[\w-]+$`)
)

// ParseKeyValue extracts a key-value pair from a given input s.
func ParseKeyValue(s string) (key string, value string, ok bool) {
	// remove preceding whitespaces and 'export'
	s = exportRE.ReplaceAllString(s, "")

	// split into key = value
	tokens := strings.SplitN(s, "=", 2)
	if len(tokens) != 2 {
		return "", "", false
	}

	// key
	key = tokens[0]
	if !keyRE.MatchString(key) {
		return "", "", false
	}

	// value
	value = tokens[1]
	if len(value) == 0 {
		return key, "", true // empty value
	}
	rvalue := []rune(value)
	if rvalue[0] == '"' || rvalue[0] == '\'' { // quoted value
		for i := 1; i < len(rvalue); i++ {
			if rvalue[i] == rvalue[0] {
				return key, string(rvalue[1:i]), true
			}
		}
		return "", "", false
	} else {
		// unquoted value: take the first word
		for i, c := range rvalue {
			if unicode.IsSpace(c) {
				return key, string(rvalue[:i]), true
			}
		}

		return key, value, true
	}
}

// Load reads the file and sets environment variables.
//
// Load will return an error if it fails to read the input file, but,
// invalid/illegal syntax in the input file will be simply ignored, and,
// will not cause this function to return an error.
func Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		k, v, ok := ParseKeyValue(scanner.Text())
		if ok {
			_ = os.Setenv(k, v)
		}
	}

	return nil
}
