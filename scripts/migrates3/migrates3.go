package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/xxxsen/common/errs"
)

var src = flag.String("src", "", "file to upload")
var uploadhost = flag.String("upload_host", "127.0.0.1:9901", "upload host")
var user = flag.String("user", "abc", "user")
var pwd = flag.String("pwd", "123456", "pwd")

func readFullLine(f string) ([]string, error) {
	raw, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	linedata := string(raw)
	lines := strings.Split(linedata, "\n")
	return lines, nil
}

func downloadStream(client *http.Client, uri string) (string, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", nil, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return "", nil, errs.New(errs.ErrUnknown, "status code:%d", rsp.StatusCode)
	}
	raw, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", nil, err
	}
	return req.URL.Path, raw, nil
}

func uploadStream(client *http.Client, path string, raw []byte, host string, user string, pwd string) error {
	pts := strings.SplitN(strings.Trim(path, "/"), "/", 2)
	if len(pts) < 2 {
		return errs.New(errs.ErrParam, "path should contains bucket and obj name, path:%s", path)
	}
	fulluri := fmt.Sprintf("http://%s/s3%s", host, path)
	req, err := http.NewRequest(http.MethodPut, fulluri, bytes.NewReader(raw))
	if err != nil {
		return errs.Wrap(errs.ErrParam, "build upload request fail", err)
	}
	rsp, err := client.Do(req)
	if err != nil {
		return errs.Wrap(errs.ErrIO, "do http request fail", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return errs.New(errs.ErrParam, "upload fail, statuscode:%d, fulluri:%s", rsp.StatusCode, fulluri)
	}
	return nil
}

func uploadOne(uri string, host string, user string, pwd string) error {
	client := &http.Client{}
	path, raw, err := downloadStream(client, uri)
	if err != nil {
		return errs.Wrap(errs.ErrIO, "download stream fail", err)
	}
	if err := uploadStream(client, path, raw, host, user, pwd); err != nil {
		return errs.Wrap(errs.ErrIO, "upload stream fail", err)
	}
	log.Printf("upload uri:%s succ", uri)
	return nil
}

func upload(lines []string, host string, user string, pwd string) error {
	for _, line := range lines {
		if len(line) == 0 {
			break
		}
		if err := uploadOne(line, *uploadhost, user, pwd); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	lines, err := readFullLine(*src)
	if err != nil {
		panic(err)
	}
	if err := upload(lines, *uploadhost, *user, *pwd); err != nil {
		panic(err)
	}
	log.Printf("upload finish")
}
