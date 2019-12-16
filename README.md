# ptrtoerr

Detects conversion from pointer to error interface.

## Usage

```sh
go get github.com/makiuchi-d/ptrtoerr/cmd/ptrtoerr
go vet -vettool `which ptrtoerr` <your-source.go>
```

## Example

### example.go

```go
package main

import "fmt"

type MyErr struct {
	error
}

func F() *MyErr {
	return nil // nil pointer, not nil interface
}

func F2() error {
	return F() // nil pointer is not a nil error
}

func main() {
	var err error

	err = F() // nil pointer is not a nil error
	if err != nil {
		fmt.Println("err is not nil")
	}

	err = F2()
	if err != nil {
		fmt.Println("err is not nil")
	}
}
```

### Result

```sh
$ go vet -vettool `which ptrtoerr` example.go
# command-line-arguments
testdata/t.go:14:2: Return pointer as error
testdata/t.go:20:2: Assign pointer to error
```
