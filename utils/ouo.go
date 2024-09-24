package utils

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"

	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func OuoBypass(ouoURL string) (string, error) {
	tempURL := strings.Replace(ouoURL, "ouo.press", "ouo.io", 1)
	var res string
	u, err := url.Parse(tempURL)
	if err != nil {
		return "", err
	}

	id := tempURL[strings.LastIndex(tempURL, "/")+1:]
	jar := tlsclient.NewCookieJar()
	options := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithClientProfile(profiles.Chrome_110),
		tlsclient.WithNotFollowRedirects(),
		tlsclient.WithCookieJar(jar),
	}

	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger(), options...)
	if err != nil {
		return "", err
	}

	getReq, err := http.NewRequest(http.MethodGet, tempURL, nil)
	if err != nil {
		return "", err
	}

	const chrome110UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"
	const accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	const acceptEncoding = "gzip, deflate, br, zstd"
	const acceptLang = "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"
	getReq.Header = http.Header{
		"accept":                    {accept},
		"accept-encoding":           {acceptEncoding},
		"accept-language":           {acceptLang},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {chrome110UserAgent},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}

	resp, err := client.Do(getReq)
	if err != nil {
		return "", err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode == 403 {
		return "", errors.New("ouo.io is blocking the request")
	}
	readBytes, _ := io.ReadAll(resp.Body)
	data := url.Values{}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(readBytes)))
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if strings.HasSuffix(name, "token") {
			value, _ := s.Attr("value")
			data.Add(name, value)
		}
	})
	nextURL := fmt.Sprintf("%s://%s/go/%s", u.Scheme, u.Host, id)

	recaptchaV3, err := ReCaptcha("https://www.google.com/recaptcha/api2/anchor?ar=1&k=6Lcr1ncUAAAAAH3cghg6cOTPGARa8adOf-y9zv2x&co=aHR0cHM6Ly9vdW8uaW86NDQz&hl=zh-CN&v=rKbTvxTxwcw5VqzrtN-ICwWt&size=invisible&cb=cuzyb4r7cdyg")
	if err != nil {
		return "", err
	}
	data.Set("x-token", recaptchaV3)
	for i := 0; i < 2; i++ {
		postReq, err := http.NewRequest(http.MethodPost, nextURL, strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}
		postReq.Header = http.Header{
			"accept":                    {accept},
			"content-type":              {"application/x-www-form-urlencoded"},
			"accept-encoding":           {acceptEncoding},
			"accept-language":           {acceptLang},
			"upgrade-insecure-requests": {"1"},
			"user-agent":                {chrome110UserAgent},
		}
		resp, err := client.Do(postReq)
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		defer func() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()
		if resp.StatusCode == 302 {
			res = resp.Header.Get("Location")
			break
		} else if resp.StatusCode == 403 {
			return "", errors.New("ouo.io is blocking the request")
		}
		nextURL = fmt.Sprintf("%s://%s/xreallcygo/%s", u.Scheme, u.Host, id)
	}
	return res, nil
}
