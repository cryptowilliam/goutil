package youtube

// youtube的auth方法
// https://godoc.org/google.golang.org/api/youtube/v3#hdr-Other_authentication_options

import (
	"context"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
	"google.golang.org/api/youtube/v3"
	"net"
	"net/http"
	"net/url"
	"time"
)

type (
	Youtube struct {
		service *youtube.Service
	}
)

// https://github.com/googleapis/google-api-go-client/issues/434#issuecomment-563428458
func newServiceWithProxy(ctx context.Context, proxy string, opts ...option.ClientOption) (*youtube.Service, error) {
	scopesOption := option.WithScopes(
		"https://www.googleapis.com/auth/youtube",
		"https://www.googleapis.com/auth/youtube.force-ssl",
		"https://www.googleapis.com/auth/youtube.readonly",
		"https://www.googleapis.com/auth/youtube.upload",
		"https://www.googleapis.com/auth/youtubepartner",
		"https://www.googleapis.com/auth/youtubepartner-channel-audit",
	)
	// NOTE: prepend, so we don't override user-specified scopes.
	opts = append([]option.ClientOption{scopesOption}, opts...)
	client, endpoint, err := htransport.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}

		baseTransportWithProxy := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		client.Transport.(*transport.APIKey).Transport = baseTransportWithProxy
	}
	s, err := youtube.New(client)
	if err != nil {
		return nil, err
	}
	if endpoint != "" {
		s.BasePath = endpoint
	}
	return s, nil
}

// NewYoutubeService create youtube client
func NewWithKey(key, proxy string) (*Youtube, error) {
	service, err := newServiceWithProxy(context.Background(), proxy, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	return &Youtube{service: service}, nil
}
