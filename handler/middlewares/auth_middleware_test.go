package middlewares

import (
	"fileserver/utils"
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

	utils.CreateCodeAuthRequest(r, ak, sk, uint64(now))

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
