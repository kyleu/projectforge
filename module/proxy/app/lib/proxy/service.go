package proxy

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

type Service struct {
	urlPrefix string
	proxies   map[string]string
}

func NewService(urlPrefix string, initialProxies map[string]string) *Service {
	return &Service{urlPrefix: urlPrefix, proxies: initialProxies}
}

func (s *Service) SetURL(svc string, url string) {
	s.proxies[svc] = url
}

func (s *Service) Remove(svc string) {
	delete(s.proxies, svc)
}

func (s *Service) List() []string {
	return util.MapKeysSorted(s.proxies)
}

func (s *Service) Handle(ctx context.Context, svc string, w http.ResponseWriter, r *http.Request, pth string, logger util.Logger) error {
	url, err := s.urlFor(svc, pth, r.URL.Query())
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, r.Method, url, r.Body)
	if err != nil {
		return err
	}
	replaceHeaders(r.Header, req.Header)
	req.Header.Set("Proxied", util.AppKey)
	logger.Infof("handling request [%s] for service [%s]", pth, svc)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Check if response is compressed
	var reader io.Reader = rsp.Body
	contentEncoding := rsp.Header.Get("Content-Encoding")
	if contentEncoding == "gzip" {
		gzReader, err := gzip.NewReader(rsp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to create gzip reader")
		}
		defer gzReader.Close()
		reader = gzReader
		// Remove Content-Encoding header since we're sending uncompressed data
		rsp.Header.Del("Content-Encoding")
	}

	rspBody, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	rspString := string(rspBody)
	proxyPath := fmt.Sprintf("%s/%s", s.urlPrefix, svc)
	rspString = strings.ReplaceAll(rspString, "href=\"/", fmt.Sprintf("href=\"%s/", proxyPath))
	rspString = strings.ReplaceAll(rspString, "src=\"/", fmt.Sprintf("src=\"%s/", proxyPath))
	replaceHeaders(rsp.Header, w.Header(), len(rspString))
	w.WriteHeader(rsp.StatusCode)

	size, err := w.Write([]byte(rspString))
	if err != nil {
		return err
	}

	rSz, sz := util.ByteSizeSI(int64(len(rspString))), util.ByteSizeSI(int64(size))
	logger.Infof("response [%d] received [%s/%s] from [%s] url [%s]", rsp.StatusCode, rSz, sz, svc, req.URL.String())
	return nil
}

func (s *Service) urlFor(svc string, pth string, q url.Values) (string, error) {
	u, ok := s.proxies[svc]
	if !ok {
		return "", errors.Errorf("service [%s] is not registered", svc)
	}
	u = strings.TrimPrefix(u, "/")
	if !strings.HasPrefix(pth, "/") {
		pth = "/" + pth
	}
	u += pth
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	return u, nil
}

var badHeaders = []string{"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization", "Te", "Trailers", "Transfer-Encoding", "Upgrade"}

func replaceHeaders(src http.Header, dst http.Header, contentLength ...int) {
	for k := range dst {
		dst.Del(k)
	}
	cl := -1
	if len(contentLength) == 1 {
		cl = contentLength[0]
	}
	for k, vv := range src {
		for _, v := range vv {
			if cl >= 0 && k == "Content-Length" {
				dst.Add(k, fmt.Sprintf("%d", cl))
			} else if !slices.Contains(badHeaders, v) {
				dst.Add(k, v)
			}
		}
	}
}
