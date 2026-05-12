# go-slicest

A collection of generic slice helpers.

## Add to your project

```sh
go get github.com/jannes-dailidow/go-slicest
```

## Project Structure

### Files

The functions are grouped into files based on their general purpose.

### Function name

The function names start with their purpose (the filename) and can have additional suffixes that indicate different variants:
| Suffix | Meaning                                                     |
| ------ | ----------------------------------------------------------- |
| X      | Callback and function return an `error`                     |
| I      | Callback is provided with index                             |
| C      | Callback and function are provided with a `context.Context` |
| Value  | Used when a single value is used instead of a callback      |

Suffixes can be combined in the order shown above to make the perfect tool for any use case.

### Implementation
- Variables referencing a generic should use the same name in lowercase.
- Variables used for collecting the result should be named `result`.
- Preallocate slice results when the final size is known, but never make guesses!
- Order the functions from simple to complex in each file.
- Reuse complex versions of the same function to avoid code duplication.
- Do not chain multiple reuses. It bloats the callstack and makes debugging even more painful. In most cases, implementing the most complex version of a function and just reusing that will suffice.
