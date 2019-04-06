# go-env

[![GoDoc](https://godoc.org/github.com/d5/go-env?status.svg)](https://godoc.org/github.com/d5/go-env)

Simple Go library to load environment variables from files.

```golang
package main

import github.com/d5/go-env

func main() {
  // read ".env" file and set environment variables
  env.Load(".env")

  // ...
}
```

## Environment File Syntax 

Basic format per line is `key=value`.

```
KEY1=VALUE1
KEY2=VALUE2
```

You can include the comment using `#` character:

```
# line commenets
FOO=bar          # trailing whitespaces are ignored: key "FOO" value "bar"
```

Preceding and trailing whitespaces are ignored and not included when setting the environment variables:

```
  FOO=bar        # key "FOO" value "bar"   
```

But, the whitespaces between key and value are illegal.

```
FOO = bar        # illegal
```

Use double/single quotes to set multi-words values.

```
FOO="bar bar"    # key "FOO" value "bar bar"
FOO=bar bar      # if not quoted, only the first words will be taken
                 # key "FOO" value "bar"
FOO='bar # bar'  # key "FOO" value "bar # bar"   
FOO= bar         # empty value: key "FOO" value ""              
```

For the compatibility with Bash, preceding 'export' is ignored.

```
export FOO=bar  # key "FOO" value "bar
```

Additional notes:

- Variable names can contain alphabetical letters, numbers, underlines (`_`), and, dash (`-`).
- Comment lines or invalid/illegal lines are ignored by `env.Load()` function. 
