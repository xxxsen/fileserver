package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fileserver/utils"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/xxxsen/common/errs"
)

type Client struct {
	c *config
}

func New(opts ...Option) (*Client, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.addr) == 0 || len(c.ak) == 0 || len(c.sk) == 0 {
		return nil, errs.New(errs.ErrParam, "invalid param")
	}
	return &Client{c: c}, nil
}

func (c *Client) buildAPI(api string) string {
	return fmt.Sprintf("%s%s", c.c.addr, api)
}

func (c *Client) attachAuth(req *http.Request) {
	ts := uint64(time.Now().Add(60 * time.Second).Unix())
	utils.CreateCodeAuthRequest(req, c.c.ak, c.c.sk, ts)
}

func (c *Client) formUpload(ctx context.Context, api string, kv map[string]string, name string, file io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range kv {
		writer.WriteField(k, v)
	}
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return fmt.Errorf("create form file fail, err:%w", err)
	}
	io.Copy(part, file)
	writer.Close()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildAPI(api), body)
	if err != nil {
		return fmt.Errorf("build request fail, err:%w", err)
	}
	httpReq.Header.Add("Content-Type", writer.FormDataContentType())
	c.attachAuth(httpReq)
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call http request fail, err:%w", err)
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%d", httpRsp.StatusCode)
	}
	return nil
}

func (c *Client) jsonPost(ctx context.Context, api string, req interface{}, rsp interface{}) error {
	raw, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("encode json fail, err:%w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildAPI(api), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("build request fail, err:%w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	c.attachAuth(httpReq)
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call http request fail, err:%w", err)
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("http code:%d", httpRsp.StatusCode)
	}

	body, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("read data fail, err:%w", err)
	}
	frame := &JsonMessageFrame{
		Data: rsp,
	}
	if err := json.Unmarshal(body, frame); err != nil {
		return fmt.Errorf("decode data fail, err:%w", err)
	}
	if frame.Code != 0 {
		return fmt.Errorf("code:%d, msg:%s", frame.Code, frame.Message)
	}
	return nil
}

func (c *Client) BeginUpload(ctx context.Context, req *BeginUploadRequest) (*BeginUploadResponse, error) {
	rsp := &BeginUploadResponse{}
	if err := c.jsonPost(ctx, apiBeginUpload, req, rsp); err != nil {
		return nil, fmt.Errorf("call begin upload fail, err:%w", err)
	}
	return rsp, nil
}

func (c *Client) PartUpload(ctx context.Context, req *PartUploadRequest) (*PartUploadResponse, error) {
	spartid := fmt.Sprintf("%d", req.PartID)
	name := "part_" + spartid
	m := map[string]string{
		"upload_ctx": req.UploadCtx,
		"md5":        req.PartMD5,
		"part_id":    spartid,
	}
	if err := c.formUpload(ctx, apiPartUpload, m, name, req.Reader); err != nil {
		return nil, fmt.Errorf("call part upload fail, err:%w", err)
	}
	return &PartUploadResponse{}, nil
}

func (c *Client) EndUpload(ctx context.Context, req *EndUploadRequest) (*EndUploadResponse, error) {
	rsp := &EndUploadResponse{}
	if err := c.jsonPost(ctx, apiEndUpload, req, rsp); err != nil {
		return nil, fmt.Errorf("call end upload fail, err:%w", err)
	}
	return rsp, nil
}

func (c *Client) UploadFile(ctx context.Context, f string) (string, error) {
	st, err := os.Stat(f)
	if err != nil {
		return "", fmt.Errorf("stat file fail, err:%w", err)
	}
	sz := st.Size()
	if sz == 0 {
		return "", errs.New(errs.ErrParam, "zero size file")
	}
	file, err := os.Open(f)
	if err != nil {
		return "", fmt.Errorf("open file fail, err:%w", err)
	}
	defer file.Close()
	//begin
	beginRsp, err := c.BeginUpload(ctx, &BeginUploadRequest{
		FileSize: sz,
	})
	if err != nil {
		return "", fmt.Errorf("begin upload fail, err:%w", err)
	}
	//part
	count := utils.CalcFileBlockCount(uint64(sz), uint64(beginRsp.BlockSize))
	for i := 0; i < count; i++ {
		partid := i + 1
		r := io.LimitReader(file, int64(beginRsp.BlockSize))
		md5v, err := utils.ReaderMd5(r)
		if err != nil {
			return "", fmt.Errorf("calc md5 fail, err:%w", err)
		}
		if _, err := file.Seek(int64(i*int(beginRsp.BlockSize)), io.SeekStart); err != nil {
			return "", fmt.Errorf("seek fail, err:%w", err)
		}
		if _, err := c.PartUpload(ctx, &PartUploadRequest{
			PartID:    uint32(partid),
			PartMD5:   md5v,
			UploadCtx: beginRsp.UploadCtx,
			Reader:    io.LimitReader(file, int64(beginRsp.BlockSize)),
		}); err != nil {
			return "", fmt.Errorf("part upload fail, err:%w", err)
		}
	}
	//end
	endRsp, err := c.EndUpload(ctx, &EndUploadRequest{
		UploadCtx: beginRsp.UploadCtx,
		FileName:  filepath.Base(f),
	})
	if err != nil {
		return "", fmt.Errorf("end upload fail, err:%w", err)
	}
	return endRsp.DownKey, nil
}
