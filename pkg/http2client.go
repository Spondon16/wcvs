package pkg

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/http2"
)

// http2Client adapts WCVS's fasthttp.Request/fasthttp.Response objects onto a
// real HTTP/2 connection. fasthttp has no HTTP/2 client support at all, so
// this is a second, minimal transport used only when -http2/-h2 is passed.
// It implements the same Do(req, resp) shape as *fasthttp.Client, so every
// existing caller (firstRequest, deception.go, recon.go) works unchanged.
//
// ponytail: TLS+ALPN h2 only — no h2c (cleartext HTTP/2) and no proxy
// support. Request smuggling never uses this (it speaks raw HTTP/1.1 over
// its own sockets), and header-oversize/meta-character/CRLF-splice tricks
// simply error out here instead of silently degrading to HTTP/1.1, since
// HTTP/2's HPACK framing has no text-based header line for them to corrupt.
type http2Client struct {
	httpClient *http.Client
}

func newHTTP2Client() *http2Client {
	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"h2"}},
	}
	return &http2Client{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   time.Duration(Config.TimeOut) * time.Second,
			// A cache poisoning/deception check needs to see the response of
			// this exact hop, not one an automatic redirect chain lands on.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *http2Client) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	scheme := string(req.URI().Scheme())
	if !strings.EqualFold(scheme, "https") {
		return fmt.Errorf("http2: %s is not supported over -http2 (cleartext h2c isn't implemented); scan without -http2 or use an https target", req.URI().String())
	}

	httpReq, err := http.NewRequest(string(req.Header.Method()), req.URI().String(), bytes.NewReader(req.Body()))
	if err != nil {
		return fmt.Errorf("http2: building request: %w", err)
	}

	for k, v := range req.Header.All() {
		key := string(k)
		if strings.EqualFold(key, "Host") || strings.EqualFold(key, "Content-Length") {
			continue // set from req.Host()/body below, not copied as a header
		}
		httpReq.Header.Add(key, string(v))
	}
	if host := req.Host(); len(host) > 0 {
		httpReq.Host = string(host)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http2: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("http2: reading response body: %w", err)
	}

	resp.Reset()
	resp.SetStatusCode(httpResp.StatusCode)
	resp.SetBody(body)
	for key, values := range httpResp.Header {
		for _, v := range values {
			resp.Header.Add(key, v)
		}
	}

	return nil
}
