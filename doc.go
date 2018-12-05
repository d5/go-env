// A simple library to load environment variables from the file.
//
// Example:
//  import github.com/d5/go-env
//  ...
//  env.Load(".env")
//
// File syntax:
//   # line comment
//   KEY1=VALUE1
//   KEY2=VALUE2         # inline comment
//   KEY4="VALUE4 FOO"   # use double quotes to include space(s)
//   KEY5="VALUE5\nFOO"  #   or line break(s)
//   KEY6 = VALUE6       # spaces are ignored (key/value are trimmed unless wrapped double quotes)
//   export KEY3=VALUE3  # "export" before key term is simply ignored
//
package env
