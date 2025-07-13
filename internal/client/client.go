package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/siyibai/remote_debug/config"
)

type Client struct {
	Cli     *http.Client
	CliURL  string
	Timeout time.Duration
}

var (
	client *Client
	once   sync.Once
)

func NewClient() *Client {
	once.Do(func() {
		cfg, err := config.GetConfig()
		if err != nil {
			return
		}
		url := fmt.Sprintf("http://%s:%s", cfg.RemoteAddress.RemoteIP, cfg.RemoteAddress.RemotePort)
		timeout := cfg.RemoteAddress.Timeout
		cli := &http.Client{}
		client = &Client{Cli: cli, CliURL: url, Timeout: timeout}
	})
	return client
}

func (c *Client) Do(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	headers map[string]string,
) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.CliURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置默认 headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置超时头信息，让 server 端知道客户端的超时设置
	timeoutSeconds := int(c.Timeout.Seconds())
	req.Header.Set("X-Timeout", fmt.Sprintf("%d", timeoutSeconds))

	return c.doRequest(ctx, req)
}

func (c *Client) Post(ctx context.Context, path string, body []byte) ([]byte, error) {
	return c.Do(ctx, "POST", path, bytes.NewBuffer(body), map[string]string{
		"Content-Type": "application/json",
	})
}

func (c *Client) Get(ctx context.Context, path string, headers map[string]string) ([]byte, error) {
	return c.Do(ctx, "GET", path, nil, headers)
}

func (c *Client) UploadFile(ctx context.Context, path string, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	writer.Close()

	return c.Do(ctx, "POST", path, body, map[string]string{
		"Content-Type": writer.FormDataContentType(),
	})
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	resp, err := c.Cli.Do(req)
	if err != nil {
		// 网络错误、超时等
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 统一处理非2xx响应
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
