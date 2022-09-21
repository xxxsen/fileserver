package core

import "fmt"

const (
	StTypeS3    = 1
	StTypeTGBot = 2
)

func NameToType(st string) int {
	for typ, name := range typeNameMap {
		if name == st {
			return typ
		}
	}
	panic(fmt.Errorf("unknown type:%s", st))
}

func TypeToName(st int) string {
	if v, ok := typeNameMap[st]; ok {
		return v
	}
	panic(fmt.Errorf("unknown type:%d", st))
}

var typeNameMap = map[int]string{
	StTypeS3:    "s3",
	StTypeTGBot: "tgbot",
}
