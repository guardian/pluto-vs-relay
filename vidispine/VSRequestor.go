package vidispine

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type VSRequestor struct {
	url    url.URL
	auth   string
	client http.Client
}

/**
initialise a new VSRequestor object
*/
func NewVSRequestor(url url.URL, user string, passwd string) *VSRequestor {
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
		// Set this value so that the underlying transport round-tripper
		// doesn't try to auto decode the body of objects with
		// content-encoding set to `gzip`.
		//
		// Refer:
		//    https://golang.org/src/net/http/transport.go?h=roundTrip#L1843
		DisableCompression: true,
	}
	authpart := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, passwd)))
	computedAuthString := fmt.Sprintf("Basic %s", authpart)

	return &VSRequestor{
		url:  url,
		auth: computedAuthString,
		client: http.Client{
			Transport: transport,
		},
	}
}

func (r *VSRequestor) url_to_call(subpath string) (*url.URL, error) {
	var urlToCall url.URL

	if strings.Contains(subpath, "://") {
		parsed_url, url_err := url.Parse(subpath)
		if url_err != nil {
			return nil, url_err
		}
		if parsed_url.Host != r.url.Host {
			return nil, errors.New("Absolute URL was not to the designated Vidispine host")
		}
		urlToCall = r.url
		urlToCall.Path = parsed_url.Path
	} else {
		urlToCall = r.url
		urlToCall.Path = subpath
	}
	return &urlToCall, nil
}

func (r *VSRequestor) Do(method string, subpath string, accept string, bodyContentType string, body io.Reader) (io.ReadCloser, error) {
	urlToCall, url_err := r.url_to_call(subpath)
	if url_err != nil {
		return nil, url_err
	}

	log.Printf("Performing %s to %s", method, urlToCall.String())

	req, _ := http.NewRequest("GET", urlToCall.String(), body)
	req.Header.Add("Authorization", r.auth)
	req.Header.Add("Accept", accept)
	if bodyContentType != "" {
		req.Header.Add("Content-Type", bodyContentType)
	}
	resp, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp.Body, nil
	}

	bodyContent, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, errors.New(fmt.Sprintf("Server returned a %d but could not read the result because of %s", resp.StatusCode, readErr))
	}

	switch resp.StatusCode {
	case 400:
		fallthrough
	case 500:
		return nil, errors.New(string(bodyContent))
	case 503:
		fallthrough
	case 504:
		log.Print("Server is not responding, retrying after a few seconds...")
		time.Sleep(5 * time.Second)
		return r.Do(method, subpath, accept, bodyContentType, body)
	default:
		return nil, errors.New(fmt.Sprintf("Unexpected error code %d, server said %s", resp.StatusCode, string(bodyContent)))
	}
}

func (r *VSRequestor) Get(subpath string, contentType string) (io.ReadCloser, error) {
	return r.Do("GET", subpath, contentType, "", nil)
}

func (r *VSRequestor) Post(subpath string, accept string, bodyContentType string, body io.Reader) (io.ReadCloser, error) {
	return r.Do("POST", subpath, accept, bodyContentType, body)
}
