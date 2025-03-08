package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c, err := New(10000)
	assert.NoError(t, err)
	SetImpl(c)
	ctx := context.Background()
	{ //test single
		key := "aaa"
		value := "bbb"
		cacheValue, ok, err := Load(ctx, key, func(ctx context.Context, c ICache, k string) (interface{}, bool, error) {
			_ = c.Set(ctx, k, value, 1*time.Minute)
			return "bbb", true, nil
		})
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, cacheValue, value)
		//测试get
		cacheValue, ok, err = Get(ctx, key)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, cacheValue, value)
		//del
		_ = Del(ctx, key)
		_, ok, err = Get(ctx, key)
		assert.NoError(t, err)
		assert.False(t, ok)
	}
	{ //test multi
		ks := []string{"k1", "k2", "k3"}
		vs := []string{"v1", "v2", "v3"}
		cacheVs, err := LoadMany(ctx, ks, func(ctx context.Context, c ICache, ks []string) (map[string]interface{}, error) {
			rs := make(map[string]interface{})
			for i := 0; i < len(ks); i++ {
				rs[ks[i]] = vs[i]
				_ = c.Set(ctx, ks[i], vs[i], 1*time.Minute)
			}
			return rs, nil
		})
		assert.NoError(t, err)
		for i, k := range ks {
			{
				cacheV, ok, err := c.Get(ctx, k)
				assert.NoError(t, err)
				assert.True(t, ok)
				assert.Equal(t, vs[i], cacheV)
			}
			{
				cacheV, ok := cacheVs[k]
				assert.True(t, ok)
				assert.Equal(t, vs[i], cacheV)
			}
		}
	}
}
