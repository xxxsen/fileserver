package bot

import (
	"context"
	"fileserver/core"
	"fileserver/utils"

	"github.com/xxxsen/common/errs"
)

type MultiBot struct {
	*TGBot
	download map[uint32]*TGBot
}

func NewMultiBot(bots ...*TGBot) (*MultiBot, error) {
	if len(bots) == 0 {
		return nil, errs.New(errs.ErrParam, "no bot found")
	}
	m := &MultiBot{
		TGBot:    bots[0],
		download: make(map[uint32]*TGBot),
	}
	for _, bt := range bots {
		if _, ok := m.download[bt.GetBotHash()]; ok {
			return nil, errs.New(errs.ErrParam, "bot conflict, same bot hash:%d, chatid:%d, token:%s",
				bt.GetBotHash(), bt.GetChatId(), bt.GetToken())
		}
		m.download[bt.GetBotHash()] = bt
	}
	return m, nil
}

func (m *MultiBot) FileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (*core.FileDownloadResponse, error) {
	extra, err := utils.DecodeBotFileExtra(fctx.Extra)
	if err != nil {
		return nil, err
	}
	client, ok := m.download[extra.GetBotHash()]
	if !ok {
		return nil, errs.New(errs.ErrServiceInternal, "not found any bot with bothash:%d", extra.GetBotHash())
	}
	return client.FileDownload(ctx, fctx)
}
