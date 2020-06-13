# A Golang and Command-Line Interface to Archive.org

This package installs a command-line tool named archive.org for using Archive.org from the command-line. It also installs the Golang package for programmatic snapshot webpage to archive.org. Please report all bugs and issues on [Github](https://github.com/wabarc/archive.org/issues).

## Installation

```sh
$ go get github.com/wabarc/archive.org
```

## Usage

#### Command-line

```sh
$ archive.org https://www.google.com https://www.bbc.com

Output:
version: 0.0.1
date: unknown

1.06s  250013 https://www.bbc.com
5.94s   15303 https://www.google.com
7.01s elapsed

https://web.archive.org/web/20200613094506/https://www.bbc.com
https://web.archive.org/web/20200613094506/https://www.google.com
```

#### Go package interfaces

```go
package main

import (
	"fmt"
	"github.com/wabarc/archive.org/pkg"
	"strings"
)

func main() {
	links := []string{"https://www.google.com", "https://www.bbc.com"}
	r := ia.Wayback(links)
	fmt.Println(strings.Join(r, "\n"))
}

// Output:
// 0.96s  250013 https://www.bbc.com
// 6.06s  213507 https://www.google.com
// 7.02s elapsed
//
// https://web.archive.org/web/20200613094640/https://www.bbc.com
// https://web.archive.org/web/20200613094640/https://www.google.com
```

## License

Permissive GPL 3.0 license, see the [LICENSE](https://github.com/wabarc/archive.org/blob/master/LICENSE) file for details.

