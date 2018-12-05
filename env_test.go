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
	k, v, ok := env.ParseKeyValue("")
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(" ")
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue("   ")
	assertFalse(t, ok)

	// key = value
	k, v, ok = env.ParseKeyValue("foo=bar")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(" foo = bar ")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue("  foo  =  bar  ")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)

	// multi-words value without double quotes: fail
	k, v, ok = env.ParseKeyValue(`"key" = bar bar`)
	assertFalse(t, ok)

	// key = "value"
	k, v, ok = env.ParseKeyValue(`foo="bar"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(` foo = "bar" `)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`  foo  =  "bar"  `)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)

	// spaces wrapped in double quotes
	k, v, ok = env.ParseKeyValue(`  foo  =  " bar "`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, " bar ", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`  foo  =  " bar bar "`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, " bar bar ", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`  foo  =  " bar bar bar   "`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, " bar bar bar   ", v)
	assertTrue(t, ok)

	// line break wrapped in double quotes
	k, v, ok = env.ParseKeyValue(`foo="bar\n"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, `bar\n`, v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo="\nbar"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, `\nbar`, v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo="bar\nbar"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, `bar\nbar`, v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo="  \nbar  \nbar   \n"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, `  \nbar  \nbar   \n`, v)
	assertTrue(t, ok)

	// preceding "export" is simply ignore
	k, v, ok = env.ParseKeyValue("export foo=bar")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(" export   foo=bar")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`export foo="bar"`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`exportfoo=bar`) // "export" must be separate for sure
	assertEqualString(t, "exportfoo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)

	// incorrect double quotes
	k, v, ok = env.ParseKeyValue(`"foo" = bar`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = "bar`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = bar"`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = "bar"bar"`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = "bar " bar"`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = "bar" "bar"`)
	assertFalse(t, ok)

	// comments ignore
	k, v, ok = env.ParseKeyValue(`# full line comments`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`  # full line comments`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = bar # inline comments`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = "bar" # inline comments`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = bar   # inline comments`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "bar", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`#foo = bar`)
	assertFalse(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = #bar`) // #bar is comment
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo = #"bar"`) // #"bar" is comment
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)

	// empty value is ok
	k, v, ok = env.ParseKeyValue("foo=")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue("foo =  ")
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo=""`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo =  ""`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "", v)
	assertTrue(t, ok)
	k, v, ok = env.ParseKeyValue(`foo =  "  "`)
	assertEqualString(t, "foo", k)
	assertEqualString(t, "  ", v)
	assertTrue(t, ok)
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
	env.Load(writeToFile(""))
	diff := compareAndUnsetEnvs(curEnvs)
	if len(diff) != 0 {
		t.Error("Should be empty")
	}

	// some entries
	curEnvs = os.Environ()
	env.Load(writeToFile(`
		# line comment
		KEY1=VALUE1
		KEY2=VALUE2 		# inline comment (and an empty line below)
		export KEY3=VALUE3  # export allowed (ignored)
		KEY4="VALUE4 FOO"   # use double quotes to include space(s)
		KEY5="VALUE5\nFOO"  #   or line break(s)
		KEY6 = VALUE6       # spaces are ignored (key/value are trimmed unless wrapped double quotes)
	`))
	diff = compareAndUnsetEnvs(curEnvs)
	if len(diff) != 6 {
		t.Error("Should have 6 elements")
	}
	assertEqualString(t, "VALUE1", diff["KEY1"])
	assertEqualString(t, "VALUE2", diff["KEY2"])
	assertEqualString(t, "VALUE3", diff["KEY3"])
	assertEqualString(t, "VALUE4 FOO", diff["KEY4"])
	assertEqualString(t, `VALUE5\nFOO`, diff["KEY5"])
	assertEqualString(t, "VALUE6", diff["KEY6"])
}

func compareAndUnsetEnvs(oldEnvs []string) map[string]string {
	oldEnvMap := envsToMap(oldEnvs)
	curEnvMap := envsToMap(os.Environ())

	diff := make(map[string]string)
	for ek, ev := range curEnvMap {
		if _, ok := oldEnvMap[ek]; !ok {
			diff[ek] = ev
			os.Unsetenv(ek)
		}
	}

	return diff
}

func envsToMap(envs []string) map[string]string {
	envMap := make(map[string]string)
	for _, env := range envs {
		tokens := strings.SplitN(env, "=", 2)
		envMap[tokens[0]] = tokens[1]
	}
	return envMap
}

func writeToFile(s string) string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(s); err != nil {
		panic(err)
	}

	return file.Name()
}

func assertTrue(t *testing.T, value bool) {
	if !value {
		t.Errorf("Should be true")
	}
}

func assertFalse(t *testing.T, value bool) {
	if value {
		t.Errorf("Should be false")
	}
}

func assertEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected: %q, Actual: %q", expected, actual)
	}
}