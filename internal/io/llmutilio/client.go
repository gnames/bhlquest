package llmutilio

import (
	"bytes"
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
		IdleConnTimeout: 300 * time.Second,
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

func (l *llmutilio) embed(texts []string) ([][]float32, error) {
	ctx := context.Background()
	url := fmt.Sprintf("%sembed", l.url)
	bs, err := l.enc.Encode(embedPayload{Texts: texts})
	if err != nil {
		return nil, err
	}
	pld := bytes.NewReader(bs)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pld)
	if err != nil {
		err = fmt.Errorf("cannot create embed request: %w", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(request)
	if err != nil {
		err = fmt.Errorf("cannot get embed response: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	bs, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("cannot read embed body: %w", err)
		return nil, err
	}
	var res [][]float32
	err = l.enc.Decode(bs, &res)
	if err != nil {
		err = fmt.Errorf("cannot decode embed body: %w", err)
		return nil, err
	}

	return res, nil
}
