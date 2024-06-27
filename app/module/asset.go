package module

import (
	"context"
	"net/http"
	"path"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

const (
	assetURL       = "https://api.github.com/repos/kyleu/projectforge/releases/latest"
	backupAssetURL = "https://projectforge.dev/assets/module"
	assetPrefix    = "projectforge_module_"
	assetSuffix    = ".zip"
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
		return "", errors.Errorf(msg, key, strings.Join(util.ArraySorted(lo.Keys(assetMap)), ", "))
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
	assetMap = map[string]string{}
	if runtime.GOOS == "js" {
		return loadBackupAssetMap(logger)
	}
	logger.Infof("loading assets from [%s]", assetURL)
	req := util.NewHTTPRequest(ctx, http.MethodGet, assetURL)
	req.Header.Set("Access-Control-Allow-Origin", "*")
	rsp, b, err := req.WithClient(telemetry.HTTPClient()).RunSimple()
	if err != nil {
		return loadBackupAssetMap(logger)
	}
	if rsp.StatusCode != 200 {
		return errors.Errorf("release asset [%s] returned status [%d]", assetURL, rsp.StatusCode)
	}
	x := &ghRsp{}
	err = util.FromJSON(b, &x)
	if err != nil {
		return errors.Wrapf(err, "release asset at [%s] returned invalid JSON", assetURL)
	}
	lo.ForEach(x.Assets, func(asset *ghAsset, _ int) {
		if strings.HasPrefix(asset.Name, assetPrefix) {
			key := strings.TrimSuffix(asset.Name[len(assetPrefix):], assetSuffix)
			assetMap[key] = asset.URL
		}
	})
	return nil
}

func loadBackupAssetMap(logger util.Logger) error {
	logger.Info("unable to download assets from github, using backup option")
	lo.ForEach(nativeModuleKeys, func(key string, _ int) {
		assetMap[key] = path.Join(backupAssetURL, key+".zip")
	})
	return nil
}
