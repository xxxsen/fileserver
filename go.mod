module fileserver

go 1.15

require (
	github.com/didi/gendry v1.7.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/uuid v1.3.0
	github.com/hnlq715/golang-lru v0.3.0
	github.com/stretchr/testify v1.8.0
	github.com/xxxsen/common v0.0.8
	go.uber.org/zap v1.23.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/xxxsen/common => ../common
