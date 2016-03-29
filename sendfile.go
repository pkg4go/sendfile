// Note: port something from go/src/net/http/fs.go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style

package sendfile

import "path/filepath"
import "net/http"
import "strconv"
import "strings"
import "errors"
import "time"
import "mime"
import "os"
import "io"

var unixEpochTime = time.Unix(0, 0)

func Send(res http.ResponseWriter, req *http.Request, path string) error {
	if containsDotDot(path) {
		res.WriteHeader(http.StatusForbidden)
		return errors.New("malicious path")
	}
	if filepath.Separator != '/' && strings.ContainsRune(path, filepath.Separator) ||
		strings.Contains(path, "\x00") {
		res.WriteHeader(http.StatusForbidden)
		return errors.New("invalid character in file path")
	}
	if !filepath.IsAbs(path) {
		res.WriteHeader(http.StatusForbidden)
		return errors.New("path must be absolute")
	}

	f, err := open(path)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return err
	}
	defer f.Close()

	stat, err := f.Stat()

	if checkLastModified(res, req, stat.ModTime()) {
		return nil
	}

	_, haveType := res.Header()["Content-Type"]
	if !haveType {
		// set content-type
		t := mime.TypeByExtension(filepath.Ext(path))
		if t != "" {
			res.Header().Set("Content-Type", t)
		}
	}

	if stat.Size() > 0 {
		res.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	}

	if req.Method != "HEAD" {
		io.Copy(res, f)
	}

	return nil
}

func open(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if d.IsDir() {
		return nil, errors.New("can't send dir")
	}

	return f, nil
}

func checkLastModified(res http.ResponseWriter, req *http.Request, modtime time.Time) bool {
	if modtime.IsZero() || modtime.Equal(unixEpochTime) {
		return false
	}

	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := res.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		res.WriteHeader(http.StatusNotModified)
		return true
	}
	res.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return false
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}

	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool {
	return r == '/' || r == '\\'
}
