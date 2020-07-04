# A Golang and Command-Line Interface to Archive.org

This package is a command-line tool named `archive.org` saving webpage to [Internet Archive](https://archive.org), it also supports imports as a Golang package for a programmatic. Please report all bugs and issues on [Github](https://github.com/wabarc/archive.org/issues).

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

https://www.bbc.com => https://web.archive.org/web/20200613094506/https://www.bbc.com
https://www.google.com => https://web.archive.org/web/20200613094506/https://www.google.com
```

#### Go package interfaces

```go
package main

package ia

import (
        "fmt"

        "github.com/wabarc/archive.org/pkg"
)

func main() {
        wbrc := &ia.Archiver{}
        saved, _ := wbrc.Wayback(args)
        for orig, dest := range saved {
                fmt.Println(orig, "=>", dest)
        }
}

// Output:
// https://www.bbc.com => https://web.archive.org/web/20200613094640/https://www.bbc.com
// https://www.google.com => https://web.archive.org/web/20200613094640/https://www.google.com
```

## License

Permissive GPL 3.0 license, see the [LICENSE](https://github.com/wabarc/archive.org/blob/master/LICENSE) file for details.

