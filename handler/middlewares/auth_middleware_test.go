package middlewares

import (
	"fileserver/utils"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCodeAuth(t *testing.T) {
	at := &codeAuth{}
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1/test", nil)
	assert.NoError(t, err)

	ak := "abc"
	sk := "123456"
	now := time.Now().Unix()

	r.Header.Set("x-fs-ak", ak)
	r.Header.Set("x-fs-ts", fmt.Sprintf("%d", now))
	code := utils.GetMd5([]byte(fmt.Sprintf("%s:%s:%d", ak, sk, now)))
	r.Header.Set("x-fs-code", code)

	{
		users := map[string]string{
			"test": "123456",
			"abc":  "123456",
		}
		ckak, pass, err := at.Auth(&gin.Context{
			Request: r,
		}, users)
		assert.NoError(t, err)
		assert.True(t, pass)
		assert.Equal(t, ak, ckak)
	}
	{
		users := map[string]string{
			"test": "123456",
			"abc":  "1234567",
		}
		ckak, pass, err := at.Auth(&gin.Context{
			Request: r,
		}, users)
		assert.NoError(t, err)
		assert.False(t, pass)
		assert.NotEqual(t, ak, ckak)
	}
}

func TestBasicAuth(t *testing.T) {
	at := &basicAuth{}
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1/test", nil)
	assert.NoError(t, err)
	ak := "abc"
	sk := "123456"
	r.SetBasicAuth(ak, sk)
	{
		users := map[string]string{
			"test": "123456",
			"abc":  "123456",
		}
		ckak, pass, err := at.Auth(&gin.Context{
			Request: r,
		}, users)
		assert.NoError(t, err)
		assert.True(t, pass)
		assert.Equal(t, ak, ckak)
	}
	{
		users := map[string]string{
			"test": "123456",
			"abc":  "1234567",
		}
		ckak, pass, err := at.Auth(&gin.Context{
			Request: r,
		}, users)
		assert.NoError(t, err)
		assert.False(t, pass)
		assert.NotEqual(t, ak, ckak)
	}
}
