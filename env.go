package env

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	envFileTokens   = regexp.MustCompile(`^\s*(?:export\s+)?\s*([^\s"=]+)\s*=\s*(?:(?:"([^"]*)")|([^\s"]*))\s*$`)
	envEntryComment = regexp.MustCompile(`#.+$`)
)

// ParseKeyValue extracts a key-value pair from a given input s.
func ParseKeyValue(s string) (key string, value string, ok bool) {
	// remove comments if any
	s = envEntryComment.ReplaceAllString(s, "")

	tokens := envFileTokens.FindAllStringSubmatch(s, -1)
	if len(tokens) == 1 && len(tokens[0]) == 4 {
		key = tokens[0][1]

		if tokens[0][2] != "" {
			// Value is wrapped in double quotes.
			// In this case, we replace "\n" with real line break character.
			value = strings.Replace(tokens[0][2], "\n", `\n`, -1)
		} else {
			value = tokens[0][3]
		}

		ok = true
	}

	return
}

// Load reads the file and sets environment variables.
func Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		k, v, ok := ParseKeyValue(scanner.Text())
		if ok {
			os.Setenv(k, v)
		}
	}

	return nil
}