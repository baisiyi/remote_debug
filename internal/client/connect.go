package client

import (
	"context"
	"encoding/json"
	"github.com/siyibai/remote_debug/internal/model"
)

type APIs interface {
	RunCommand(ctx context.Context, command string) error
	UploadFile(filePath string) error
}

type SerApi struct {
	client *Client
}

func NewSerApi() *SerApi {
	return &SerApi{
		client: NewClient(),
	}
}

func (s *SerApi) RunCommand(ctx context.Context, req *model.RunCommandReq) (err error) {
	reqByte, _ := json.Marshal(req)
	_, err = s.client.Post(ctx, "command", reqByte)
	if err != nil {
		return err
	}
	return
}

func (s *SerApi) UploadFile(ctx context.Context, filePath string) (err error) {
	_, err = s.client.UploadFile(ctx, "upload_file", filePath)
	if err != nil {
		return err
	}
	return
}
