package llmutilio

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (l *llmutilio) conn() error {
	url := l.cfg.LlmUtilURL
	if url[len(url)-1] != '/' {
		url = url + "/"
	}
	l.url = url
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Timeout: 4 * time.Minute, Transport: tr}
	l.client = client

	return l.ping()
}

func (l *llmutilio) ping() error {
	ctx := context.Background()
	url := fmt.Sprintf("%sping", l.url)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("cannot create ping request: %w", err)
		return err
	}
	request.Header.Set("Content-Type", "text/plain")

	resp, err := l.client.Do(request)
	if err != nil {
		err = fmt.Errorf("cannot get pong: %w", err)
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("cannot read pong body: %w", err)
		return err
	}

	pong := string(respBytes)

	if !strings.HasPrefix(pong, "Pong!") {
		return fmt.Errorf("Wrong pong reply: %s", pong)
	}

	return nil
}
