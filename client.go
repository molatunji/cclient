package cclient

import (
	"log"
	"time"

	http "github.com/Carcraftz/fhttp"
	"github.com/Carcraftz/fhttp/cookiejar"

	utls "github.com/Carcraftz/utls"
	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"
)

func NewClient(clientHello utls.ClientHelloID, proxyUrl string, allowRedirect bool, timeout time.Duration) (http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	if err != nil {
		log.Fatal(err)
	}
	if len(proxyUrl) > 0 {
		dialer, err := NewConnectDialer(proxyUrl)
		if err != nil {
			if allowRedirect {
				return http.Client{
					Transport: NewRoundTripper(clientHello, dialer),
					Timeout:   time.Second * timeout,
					Jar:       jar,
				}, err
			}
			return http.Client{
				Transport: NewRoundTripper(clientHello, dialer),
				Timeout:   time.Second * timeout,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Jar: jar,
			}, err
		}
		if allowRedirect {
			return http.Client{
				Transport: NewRoundTripper(clientHello, dialer),
				Timeout:   time.Second * timeout,
				Jar:       jar,
			}, nil
		}
		return http.Client{
			Transport: NewRoundTripper(clientHello, dialer),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: jar,
		}, nil
	} else {
		if allowRedirect {
			return http.Client{
				Transport: NewRoundTripper(clientHello, proxy.Direct),
				Timeout:   time.Second * timeout,
				Jar:       jar,
			}, nil
		}
		return http.Client{
			Transport: NewRoundTripper(clientHello, proxy.Direct),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: jar,
		}, nil

	}
}
