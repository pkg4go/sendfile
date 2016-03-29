package sendfile

import . "github.com/pkg4go/assert"
import "net/http/httptest"
import "path/filepath"
import "io/ioutil"
import "net/http"
import "testing"
import "os"

const expectText = `* {
  display: flex;
}
`

func TestBasic(t *testing.T) {
	a := A{t}

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cur, _ := os.Getwd()
		p := filepath.Join(cur, "fixture/a.css")
		err := Send(res, req, p)
		a.Nil(err)
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	a.Nil(err)

	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()

	a.Nil(err)

	a.Equal(res.Header.Get("Content-Type"), "text/css; charset=utf-8")
	a.Equal(res.Header.Get("Content-Length"), "23")
	a.Equal(getBody(res, a), expectText)
}

// test util
func getBody(res *http.Response, a A) string {
	body, err := ioutil.ReadAll(res.Body)
	a.Nil(err)
	text := string(body[:])
	return text
}
