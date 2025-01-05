package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	url, err := s.urlFor(svc, pth)
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
	replaceHeaders(rsp.Header, w.Header())
	w.WriteHeader(rsp.StatusCode)
	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	proxyPath := fmt.Sprintf("%s/%s", s.urlPrefix, svc)
	rspString := string(rspBody)
	rspString = strings.ReplaceAll(rspString, "href=\"/", fmt.Sprintf("href=\"%s/", proxyPath))
	rspString = strings.ReplaceAll(rspString, "src=\"/", fmt.Sprintf("src=\"%s/", proxyPath))
	size, err := w.Write([]byte(rspString))
	logger.Infof("response [%d] received [%s] from [%s] url [%s]", rsp.StatusCode, util.ByteSizeSI(int64(size)), svc, req.URL.String())
	return nil
}

func (s *Service) urlFor(svc string, pth string) (string, error) {
	u, ok := s.proxies[svc]
	if !ok {
		return "", errors.Errorf("service [%s] is not registered", svc)
	}
	if strings.HasPrefix(u, "/") {
		u = strings.TrimPrefix(u, "/")
	}
	if !strings.HasPrefix(pth, "/") {
		pth = "/" + pth
	}
	u += pth
	return u, nil
}

var badHeaders = []string{"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization", "Te", "Trailers", "Transfer-Encoding", "Upgrade"}

func replaceHeaders(src http.Header, dst http.Header) {
	for k := range dst {
		dst.Del(k)
	}
	for k, vv := range src {
		for _, v := range vv {
			if !slices.Contains(badHeaders, v) {
				println(k, v)
				dst.Add(k, v)
			}
		}
	}
}
