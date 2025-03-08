package proxyutil

import "context"

type UserInfo struct {
	AuthType string
	Username string
}

type userInfoType struct{}

var (
	defaultUserInfoKey = userInfoType{}
)

func GetUserInfo(ctx context.Context) (*UserInfo, bool) {
	c, ok := ctx.Value(defaultUserInfoKey).(*UserInfo)
	if !ok {
		return nil, false
	}
	return c, true
}

func SetUserInfo(ctx context.Context, info *UserInfo) context.Context {
	return context.WithValue(ctx, defaultUserInfoKey, info)
}
