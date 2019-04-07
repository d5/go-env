package env_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/d5/go-env"
)

func TestParseKeyValue(t *testing.T) {
	// empty input: fail
	testParseKeyValue(t, "", "", "", false)
	testParseKeyValue(t, " ", "", "", false)
	testParseKeyValue(t, "   ", "", "", false)

	// key=value
	testParseKeyValue(t, "foo=bar", "foo", "bar", true)
	testParseKeyValue(t, " foo=bar", "foo", "bar", true)
	testParseKeyValue(t, "foo=bar ", "foo", "bar", true)
	testParseKeyValue(t, "foo= bar", "foo", "", true)

	// illegal whitespaces
	testParseKeyValue(t, "foo = bar", "", "", false)
	testParseKeyValue(t, "foo =bar", "", "", false)
	testParseKeyValue(t, " foo = bar ", "", "", false)

	// 'export' ignored
	testParseKeyValue(t, "export foo=bar", "foo", "bar", true)
	testParseKeyValue(t, " export foo=bar", "foo", "bar", true)
	testParseKeyValue(t, "export  foo=bar", "foo", "bar", true)
	testParseKeyValue(t, "exportfoo=bar", "exportfoo", "bar", true)

	// unquoted multi-words (only the first word taken)
	testParseKeyValue(t, "foo=word1 word2", "foo", "word1", true)
	testParseKeyValue(t, "foo=word1 word2 word3", "foo", "word1", true)

	// invalid key
	testParseKeyValue(t, `"foo"=bar`, "", "", false)
	testParseKeyValue(t, `foo foo=bar`, "", "", false)
	testParseKeyValue(t, `foo/foo=bar`, "", "", false)
	testParseKeyValue(t, `ÏØØ=bar`, "", "", false) // non-ascii key is invalid

	// quoted value
	testParseKeyValue(t, `foo="bar"`, "foo", "bar", true)
	testParseKeyValue(t, `foo=" bar "`, "foo", " bar ", true)
	testParseKeyValue(t, `foo="word1 word2 word3"`, "foo", "word1 word2 word3", true)
	testParseKeyValue(t, `foo="bar\nbar"`, "foo", `bar\nbar`, true)
	testParseKeyValue(t, `foo='bar'`, "foo", "bar", true)
	testParseKeyValue(t, `foo=' bar '`, "foo", " bar ", true)
	testParseKeyValue(t, `foo='word1 word2 word3'`, "foo", "word1 word2 word3", true)
	testParseKeyValue(t, `foo='bar\nbar'`, "foo", `bar\nbar`, true)

	// unicode values
	testParseKeyValue(t, `foo=∫å®`, "foo", "∫å®", true)
	testParseKeyValue(t, `foo=가나다`, "foo", "가나다", true)
	testParseKeyValue(t, `foo="가 나 다"`, "foo", "가 나 다", true)
	testParseKeyValue(t, `foo='가 나 다'`, "foo", "가 나 다", true)

	// comments
	testParseKeyValue(t, "# foo=bar", "", "", false)
	testParseKeyValue(t, " #foo=bar", "", "", false)
	testParseKeyValue(t, "foo=#bar", "foo", "#bar", true)             // '#' without preceding space
	testParseKeyValue(t, `foo="bar # bar"`, "foo", "bar # bar", true) // '#' inside quoted value
	testParseKeyValue(t, `foo='bar # bar'`, "foo", "bar # bar", true) // '#' inside quoted value

	// empty value
	testParseKeyValue(t, "foo=", "foo", "", true)
	testParseKeyValue(t, "foo= ", "foo", "", true)
}

func TestLoad(t *testing.T) {
	// file not found
	notExistingFile := fmt.Sprintf("/tmp/__test_file_%d__", rand.Uint64())
	_ = os.Remove(notExistingFile)
	err := env.Load(notExistingFile)
	if err == nil {
		t.Error("Should not be nil")
	}

	// empty file: nothing should be set
	curEnvs := os.Environ()
	_ = env.Load(writeToFile(""))
	diff := compareAndUnsetEnvs(curEnvs)
	if len(diff) != 0 {
		t.Error("Should be empty")
	}

	// some entries
	curEnvs = os.Environ()
	_ = env.Load(writeToFile(`
# line comments using '#' character
 KEY1=VALUE1     # ok: preceding and trailing whitespaces are ignored
KEY2 = VALUE2    # invalid: cannot have whitespaces between key and value
KEY3=FOO#BAR     # ok: set to "FOO#BAR" ('#' without preceding space)
KEY4=FOO BAR     # ok: only the first word will be taken (KEY3=FOO)
KEY5= FOO        # ok: but set to empty (first word is empty)
KEY6='FOO BAR'   # ok: quote values for multi-words values   
KEY7="FOO # BAR" # ok: you can include '#' in quoted value	`))
	diff = compareAndUnsetEnvs(curEnvs)
	if len(diff) != 6 {
		t.Error("Should have 6 elements")
	}
	assertEqualString(t, "VALUE1", diff["KEY1"])
	assertEqualString(t, "FOO#BAR", diff["KEY3"])
	assertEqualString(t, "FOO", diff["KEY4"])
	assertEqualString(t, "", diff["KEY5"])
	assertEqualString(t, "FOO BAR", diff["KEY6"])
	assertEqualString(t, "FOO # BAR", diff["KEY7"])
}

func compareAndUnsetEnvs(oldEnvs []string) map[string]string {
	oldEnvMap := envsToMap(oldEnvs)
	curEnvMap := envsToMap(os.Environ())

	diff := make(map[string]string)
	for ek, ev := range curEnvMap {
		if _, ok := oldEnvMap[ek]; !ok {
			diff[ek] = ev
			_ = os.Unsetenv(ek)
		}
	}

	return diff
}

func envsToMap(envs []string) map[string]string {
	envMap := make(map[string]string)
	for _, e := range envs {
		tokens := strings.SplitN(e, "=", 2)
		envMap[tokens[0]] = tokens[1]
	}
	return envMap
}

func writeToFile(s string) string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	if _, err := file.WriteString(s); err != nil {
		panic(err)
	}

	return file.Name()
}

func assertEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}

func testParseKeyValue(t *testing.T, input, expectedKey, expectedValue string, expectedOK bool) {
	t.Helper()

	k, v, ok := env.ParseKeyValue(input)
	if ok != expectedOK {
		t.Logf("%s\nExpected: %v, actual: %v", input, expectedOK, ok)
		t.Fail()
		return
	}
	if k != expectedKey {
		t.Logf("%s\nExpected key: %v, actual: %v", input, expectedKey, k)
		t.Fail()
		return
	}
	if v != expectedValue {
		t.Logf("%s\nExpected value: %v, actual: %v", input, expectedValue, v)
		t.Fail()
	}

}
