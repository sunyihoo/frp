package config

import (
	"github.com/sunyihoo/frp/pkg/config/legacy"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"gopkg.in/ini.v1"
	"os"
)

func LoadServerConfig(path string, strict string) (*v1.ServerConfig, bool, error) {
	var (
		svrCfg         *v1.ServerConfig
		isLegacyFormat bool
	)
	// 检测 Legacy ini格式
	if DetectLegacyINIFormatFromFile(path) {
		content, err := legacy.GetRenderedConfFormFile(path)
		if err != nil {
			return nil, true, err
		}
		legacyCfg, err :=
	}

}

func DetectLegacyINIFormat(content []byte) bool {
	f, err := ini.Load(content)
	if err != nil {
		return false
	}
	if _, err := f.GetSection("common"); err == nil {
		return true
	}
	return false
}

func DetectLegacyINIFormatFromFile(path string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return DetectLegacyINIFormat(b)
}
