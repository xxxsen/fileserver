package getter

import (
	"context"
	"fileserver/constants"
	"fileserver/core"
	"fmt"

	"github.com/xxxsen/common/naivesvr"
)

func MustGetFsClient(ctx context.Context) core.IFsCore {
	key := constants.KeyStorageClient
	iclient, exist := naivesvr.GetAttachKey(ctx, key)
	if !exist {
		panic(fmt.Errorf("key:%s not found", key))
	}
	return iclient.(core.IFsCore)
}
