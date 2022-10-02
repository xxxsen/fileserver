package bot

import (
	"context"
	"fileserver/core"
	"fileserver/utils"
	"sync/atomic"

	"github.com/xxxsen/common/errs"
	"google.golang.org/protobuf/proto"
)

type MultiBot struct {
	botList []*TGBot
	bots    map[uint32]*TGBot
	idx     uint32
}

func NewMultiBot(bots ...*TGBot) (*MultiBot, error) {
	if len(bots) == 0 {
		return nil, errs.New(errs.ErrParam, "no bot found")
	}
	m := &MultiBot{
		botList: bots,
		bots:    make(map[uint32]*TGBot),
	}
	for _, bt := range bots {
		if _, ok := m.bots[bt.GetBotHash()]; ok {
			return nil, errs.New(errs.ErrParam, "bot conflict, same bot hash:%d, chatid:%d, token:%s",
				bt.GetBotHash(), bt.GetChatId(), bt.GetToken())
		}
		m.bots[bt.GetBotHash()] = bt
	}
	return m, nil
}

func (m *MultiBot) FileDownload(ctx context.Context, fctx *core.FileDownloadRequest) (*core.FileDownloadResponse, error) {
	extra, err := utils.DecodeBotFileExtra(fctx.Extra)
	if err != nil {
		return nil, err
	}
	client, ok := m.bots[extra.GetBotHash()]
	if !ok {
		return nil, errs.New(errs.ErrServiceInternal, "not found any bot with bothash:%d", extra.GetBotHash())
	}
	return client.FileDownload(ctx, fctx)
}

func (m *MultiBot) StType() uint8 {
	return m.botList[0].StType()
}

func (m *MultiBot) BlockSize() int64 {
	return m.botList[0].BlockSize()
}

func (m *MultiBot) MaxFileSize() int64 {
	return m.botList[0].MaxFileSize()
}

func (m *MultiBot) chooseBot() *TGBot {
	v := atomic.AddUint32(&m.idx, 1) % uint32(len(m.botList))
	return m.botList[v]
}

func (m *MultiBot) FileUpload(ctx context.Context, uctx *core.FileUploadRequest) (*core.FileUploadResponse, error) {
	return m.chooseBot().FileUpload(ctx, uctx)
}

func (m *MultiBot) BeginFileUpload(ctx context.Context, fctx *core.BeginFileUploadRequest) (*core.BeginFileUploadResponse, error) {
	bt := m.chooseBot()
	rsp, err := bt.BeginFileUpload(ctx, fctx)
	if err != nil {
		return nil, err
	}
	uctx, err := utils.DecodeUploadID(rsp.UploadID)
	if err != nil {
		return nil, errs.Wrap(errs.ErrUnmarshal, "decode upload id to append bot hash fail", err)
	}
	uctx.BotHash = proto.Uint32(bt.GetBotHash())
	//将bothash打包到uploadid中, 方便后续进行寻址
	upid, err := utils.EncodeUploadID(uctx)
	if err != nil {
		return nil, errs.Wrap(errs.ErrMarshal, "encode upload id with bot hash fail", err)
	}
	rsp.UploadID = upid
	return rsp, nil
}

func (m *MultiBot) PartFileUpload(ctx context.Context, pctx *core.PartFileUploadRequest) (*core.PartFileUploadResponse, error) {
	uctx, err := utils.DecodeUploadID(pctx.UploadId)
	if err != nil {
		return nil, errs.Wrap(errs.ErrUnmarshal, "decode upload id fail", err)
	}
	bt, ok := m.bots[uctx.GetBotHash()]
	if !ok {
		return nil, errs.New(errs.ErrNotFound, "bot with spec hash not found in server, bothash:%d", uctx.GetBotHash())
	}
	return bt.PartFileUpload(ctx, pctx)
}

func (m *MultiBot) FinishFileUpload(ctx context.Context, fctx *core.FinishFileUploadRequest) (*core.FinishFileUploadResponse, error) {
	uctx, err := utils.DecodeUploadID(fctx.UploadId)
	if err != nil {
		return nil, errs.Wrap(errs.ErrUnmarshal, "decode upload id fail", err)
	}
	bt, ok := m.bots[uctx.GetBotHash()]
	if !ok {
		return nil, errs.New(errs.ErrNotFound, "bot with spec hash not found in server, bothash:%d", uctx.GetBotHash())
	}
	return bt.FinishFileUpload(ctx, fctx)
}
