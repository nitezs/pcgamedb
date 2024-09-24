package utils

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"regexp"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
)

func ReCaptcha(anchorUrl string) (string, error) {
	urlBase := "https://www.google.com/recaptcha/api2/"

	matches := regexp.MustCompile(`/anchor\?(.*)`).FindStringSubmatch(anchorUrl)
	if len(matches) < 2 {
		return "", fmt.Errorf("no matches found in ANCHOR_URL")
	}
	params := matches[1]

	client, err := tlsclient.NewHttpClient(tlsclient.NewNoopLogger())
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, anchorUrl, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("recaptcha status code is not 200")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	tokenMatches := regexp.MustCompile(`"recaptcha-token" value="(.*?)"`).FindStringSubmatch(string(body))
	if len(tokenMatches) < 2 {
		return "", errors.New("no token found in response")
	}
	token := tokenMatches[1]
	paramsMap, err := url.ParseQuery(params)
	if err != nil {
		return "", err
	}
	paramsMap.Set("c", token)
	paramsMap.Set("reason", "q")
	reloadUrl := urlBase + "reload?k=" + paramsMap.Get("k")
	postReq, err := http.NewRequest(http.MethodPost, reloadUrl, strings.NewReader(paramsMap.Encode()))
	if err != nil {
		return "", err
	}
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(postReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	answerMatches := regexp.MustCompile(`"rresp","(.*?)"`).FindStringSubmatch(string(body))
	if len(answerMatches) < 2 {
		return "", fmt.Errorf("no answer found in reCaptcha response: %s", string(body))
	}
	return answerMatches[1], nil
}

//https://www.google.com/recaptcha/api2/anchor?ar=1&k=6Lf-ZrEUAAAAAEtmR70o2Rb9JM2QUBCH4j7EzIWX&co=aHR0cHM6Ly93d3cua2VlcGxpbmtzLm9yZzo0NDM.&hl=zh-CN&v=-80zvSY9h4i8O-ocN2P5qTJk&size=normal&cb=xm7hg2ftd5e2
