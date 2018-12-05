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

File Syntax:

```env
# line comment
KEY1=VALUE1
KEY2=VALUE2         # inline comment
KEY4="VALUE4 FOO"   # use double quotes to include space(s)
KEY5="VALUE5\nFOO"  #   or line break(s)
KEY6 = VALUE6       # spaces are ignored (key/value are trimmed unless wrapped double quotes)
export KEY3=VALUE3  # "export" before key term is simply ignored
```
