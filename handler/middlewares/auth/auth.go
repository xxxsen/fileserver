package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IAuth interface {
	Name() string
	IsMatchAuthType(ctx *gin.Context) bool
	Auth(ctx *gin.Context, users map[string]string) (string, error)
}

type AuthCreateFunc func() IAuth

var mp = make(map[string]AuthCreateFunc)

func MustCreateByName(name string) IAuth {
	at, err := CreateByName(name)
	if err != nil {
		panic(err)
	}
	return at
}

func CreateByName(name string) (IAuth, error) {
	if v, ok := mp[name]; ok {
		return v(), nil
	}
	return nil, fmt.Errorf("not found")
}

func Regist(name string, fn AuthCreateFunc) {
	mp[name] = fn
}

func AuthList() []string {
	rs := make([]string, 0, len(mp))
	for k := range mp {
		rs = append(rs, k)
	}
	return rs
}
