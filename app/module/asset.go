package module

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

const (
	assetURL    = "https://api.github.com/repos/kyleu/projectforge/releases/latest"
	assetPrefix = "projectforge_module_"
	assetSuffix = ".zip"
)

var assetMap map[string]string

func (s *Service) AssetURL(ctx context.Context, key string, logger util.Logger) (string, error) {
	if assetMap == nil {
		if err := loadAssetMap(ctx, logger); err != nil {
			return "", err
		}
	}
	ret, ok := assetMap[key]
	if !ok {
		msg := "no URL available for module [%s] among candidates [%s]"
		keys := make([]string, 0, len(assetMap))
		for k := range assetMap {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		return "", errors.Errorf(msg, key, strings.Join(keys, ", "))
	}
	return ret, nil
}

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
	Size int    `json:"size"`
}

type ghRsp struct {
	Assets []*ghAsset `json:"assets"`
}

func loadAssetMap(ctx context.Context, logger util.Logger) error {
	logger.Infof("loading assets from [%s]", assetURL)
	assetMap = map[string]string{}
	httpClient := telemetry.WrapHTTPClient(http.DefaultClient)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, assetURL, http.NoBody)
	if err != nil {
		return errors.Wrapf(err, "unable to create request from [%s]", assetURL)
	}
	rsp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "unable to get release asset from [%s]", assetURL)
	}
	if rsp.StatusCode != 200 {
		return errors.Errorf("release asset [%s] returned status [%s]", assetURL, rsp.Status)
	}
	bts, err := io.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrapf(err, "unable to read release asset from [%s]", assetURL)
	}
	x := &ghRsp{}
	err = util.FromJSON(bts, &x)
	if err != nil {
		return errors.Wrapf(err, "release asset at [%s] returned invalid JSON", assetURL)
	}
	for _, asset := range x.Assets {
		if strings.HasPrefix(asset.Name, assetPrefix) {
			key := strings.TrimSuffix(asset.Name[len(assetPrefix):], assetSuffix)
			assetMap[key] = asset.URL
		}
	}
	return nil
}
