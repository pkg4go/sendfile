
[![Build status][travis-img]][travis-url]
[![License][license-img]][license-url]
[![GoDoc][doc-img]][doc-url]

### sendfile

* It does the following:
  - Check if a file exists
  - Set content-length, content-type, and last-modified headers
  - 304 based on last-modified
  - Handle HEAD requests

* It does not:
  - Cache control
  - OPTIONS method

### Install

```bash
go get github.com/pkg4go/sendfile
```

### Example

```go
import "github.com/pkg4go/sendfile"
import "net/http"

func main() {
  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    sendfile.Send(res, req, "/path/to/file")
  })

  http.ListenAndServe(":3000", nil)
}
```

### License
MIT

[doc-img]: http://img.shields.io/badge/GoDoc-reference-green.svg?style=flat-square
[doc-url]: http://godoc.org/github.com/pkg4go/sendfile
[travis-img]: https://img.shields.io/travis/pkg4go/sendfile.svg?style=flat-square
[travis-url]: https://travis-ci.org/pkg4go/sendfile
[license-img]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[license-url]: http://opensource.org/licenses/MIT
