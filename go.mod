module fileserver

go 1.15

require (
	github.com/didi/gendry v1.7.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/xxxsen/common v0.0.8
	github.com/xxxsen/runner v0.0.1
	github.com/yitter/idgenerator-go v1.3.1
	go.uber.org/zap v1.23.0
	golang.org/x/sync v0.0.0-20220907140024-f12130a52804 // indirect
	google.golang.org/protobuf v1.28.0
)

replace (
	github.com/xxxsen/common => ../common
)
