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

func (s *SerApi) RunCommand(ctx context.Context, req *model.CommandRequest) (
	rsp *model.CommandResponse, err error) {
	reqByte, _ := json.Marshal(req)
	rspByte, err := s.client.Post(ctx, "command", reqByte)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(rspByte, &rsp)
	return
}

func (s *SerApi) UploadFile(ctx context.Context, filePath string, destPath string) (err error) {
	_, err = s.client.UploadFile(ctx, "upload", filePath, destPath)
	return
}
