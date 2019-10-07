# About
This tool allows you to ship a single executable without the need for extra resource files.
Using the provided command line tool you can append a data file to the end of your executable which should leave it intact to be run by the OS.
The library then gives you a function to read the data back in from the end of the new executable.
This way there is no need to pack your files or create an installer, you can just ship a single executable.

# Usage
Install the tool and library by running

    go get -u github.com/gonutz/payload/...

To then combine an executable and a data file, run:

    payload -exe=path/to/exe -data=path/to/data -output=path/for/combined/exe

If you leave out the `output` paramter, the original executable is overwritten.

In the executable file you can add code to read the data back in as a []byte slice, here is an example program that just reads the payload and prints its contents to stdout:

```Go
package main

import (
	"fmt"
	"github.com/gonutz/payload"
)

func main() {
	data, err := payload.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", data)
}
```

# Documentation

[Here is the Godoc for the library](https://godoc.org/github.com/gonutz/payload).

The command line tool `payload` is used as follows:

```
Usage of payload:
  -data string
        Data file to be appended to the executable.
  -exe string
        Executable file to append data to.
  -output string
        Combined output file, defaults to the input exe.
```
