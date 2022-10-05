package utils

import (
	"fmt"
	"net/http"
)

func CreateCodeAuth(ak, sk string, ts uint64) string {
	code := GetMd5([]byte(fmt.Sprintf("%s:%s:%d", ak, sk, ts)))
	return code
}

func CreateCodeAuthRequest(r *http.Request, ak, sk string, ts uint64) {
	code := CreateCodeAuth(ak, sk, ts)
	r.Header.Set("x-fs-ak", ak)
	r.Header.Set("x-fs-ts", fmt.Sprintf("%d", ts))
	r.Header.Set("x-fs-code", code)
}
