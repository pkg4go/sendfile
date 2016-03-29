package main

import "../../sendfile"
import "path/filepath"
import "net/http"
import "fmt"
import "os"

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL.Path)
		cur, _ := os.Getwd()
		sendfile.Send(res, req, filepath.Join(cur, "fixture/a.css"))
	})
	http.ListenAndServe(":3000", nil)
}
