package library

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

const MaxJSONBytes = 4 << 20

type remote struct {
	base   string
	key    string
	client *http.Client
}

func ValidateURL(value string) (string, error) {
	u, err := url.Parse(strings.TrimRight(strings.TrimSpace(value), "/"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" ||
		u.User != nil || u.RawQuery != "" || u.Fragment != "" || strings.Contains(u.Host, "\\") {
		return "", errors.New("请输入不含凭据、查询参数的 HTTP(S) 服务根地址")
	}
	if ip, err := netip.ParseAddr(u.Hostname()); err == nil && !allowedAddress(ip) {
		return "", errors.New("不允许访问元数据或链路本地地址")
	}
	return u.String(), nil
}

func allowedAddress(ip netip.Addr) bool {
	ip = ip.Unmap()
	return ip.IsValid() && !ip.IsUnspecified() && !ip.IsLinkLocalUnicast() && !ip.IsMulticast()
}

func newRemote(base, key string) (*remote, error) {
	base, err := ValidateURL(base)
	if err != nil {
		return nil, err
	}
	return &remote{base: base, key: key, client: &http.Client{
		Transport: safeTransport, Timeout: 2 * time.Minute,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}, nil
}

var safeTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Resolve and validate at dial time, then dial the validated IP itself.
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, ErrUpstream
		}
		ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
		if err != nil || len(ips) == 0 {
			return nil, ErrUpstream
		}
		for _, ip := range ips {
			if !allowedAddress(ip) {
				return nil, ErrUpstream
			}
		}
		dialer := net.Dialer{Timeout: 10 * time.Second}
		for _, ip := range ips {
			conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
			if err == nil {
				return conn, nil
			}
		}
		return nil, ErrUpstream
	}
	return transport
}()

func (r *remote) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.base+path, nil)
	if err != nil {
		return ErrUpstream
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", r.key)
	response, err := r.client.Do(req)
	if err != nil {
		return ErrUpstream
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return ErrUpstream
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, MaxJSONBytes+1))
	if len(data) > MaxJSONBytes {
		return ErrTooLarge
	}
	if err != nil || json.Unmarshal(data, out) != nil {
		return ErrUpstream
	}
	return nil
}

type envelope[T any] struct {
	Data    T     `json:"data"`
	Success *bool `json:"success"`
}

func (e envelope[T]) check(err error) error {
	if err != nil {
		return err
	}
	if e.Success != nil && !*e.Success {
		return ErrUpstream
	}
	return nil
}
