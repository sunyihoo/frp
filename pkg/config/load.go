package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/sunyihoo/frp/pkg/config/legacy"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/util/util"
	"gopkg.in/ini.v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"strings"
	"text/template"
)

var glbEnvs map[string]string

func init() {
	glbEnvs = make(map[string]string)
	envs := os.Environ()
	for _, env := range envs {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		glbEnvs[pair[0]] = pair[1]
	}
}

type Values struct {
	Envs map[string]string // 环境变量
}

func GetValues() *Values {
	return &Values{
		Envs: glbEnvs,
	}
}

func LoadServerConfig(path string, strict bool) (*v1.ServerConfig, bool, error) {
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
		legacyCfg, err := legacy.UnmarshalServerConfFromIni(content)
		if err != nil {
			return nil, true, err
		}
		svrCfg = legacy.Convert_ServerCommonConf_To_v1(&legacyCfg)
		isLegacyFormat = true
	} else {
		svrCfg = &v1.ServerConfig{}
		if err := LoadConfigureFromFile(path, svrCfg, strict); err != nil {
			return nil, false, err
		}
	}
	if svrCfg != nil {
		svrCfg.Complete()
	}
	return svrCfg, isLegacyFormat, nil
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

func RenderWithTemplate(in []byte, values *Values) ([]byte, error) {
	tmpl, err := template.New("frp").Funcs(template.FuncMap{
		"parseNumberRange":     parseNumberRange,
		"parseNumberRangePair": parseNumberRangePair,
	}).Parse(string(in))
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBufferString("")
	if err := tmpl.Execute(buffer, values); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func LoadFileContentWithTemplate(path string, values *Values) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return RenderWithTemplate(b, values)
}

func LoadConfigureFromFile(path string, c any, strict bool) error {
	content, err := LoadFileContentWithTemplate(path, GetValues())
	if err != nil {
		return err
	}
	return LoadConfigure(content, c, strict)
}

func LoadConfigure(b []byte, c any, strict bool) error {
	v1.DisallowUnKnownFieldsMu.Lock()
	defer v1.DisallowUnKnownFieldsMu.Unlock()
	v1.DisallowUnknownFields = strict

	var tomlObj interface{}
	// 首先尝试TOML unmarshal;吞下错误（假设它不是有效的 TOML）。
	// 首先尝试将数据字符串按照 TOML格式进行解析或反序列化。如果解析失败，不要让错误影响程序的其余部分，优雅地处理这个错误。
	if err := toml.Unmarshal(b, &tomlObj); err == nil {
		b, err = json.Marshal(&tomlObj)
		if err != nil {
			return err
		}
	}
	// 如果数据缓冲区的第一个非空白字符是 {，则认为数据可能是 JSON 格式的，并直接尝试将其解析为 JSON。
	if yaml.IsJSONBuffer(b) {
		decoder := json.NewDecoder(bytes.NewBuffer(b))
		if strict {
			decoder.DisallowUnknownFields()
		}
		return decoder.Decode(c)
	}
	// 如果数据不是 JSON 格式，那么尝试将其解析为 YAML 格式。
	if strict {
		return yaml.UnmarshalStrict(b, c)
	}
	return yaml.Unmarshal(b, c)
}

func NewProxyConfigurerFromMsg(m *msg.NewProxy, serverCfg *v1.ServerConfig) (v1.ProxyConfigurer, error) {
	m.ProxyType = util.EmptyOr(m.ProxyType, string(v1.ProxyTypeTCP))

	configurer := v1.NewProxyConfigurerByType(v1.ProxyType(m.ProxyType))
	if configurer == nil {
		return nil, fmt.Errorf("unknown proxy type: %s", m.ProxyType)
	}

	configurer.UnmarshalFromMsg(m)
	configurer.Complete("")

	if err := valitation.Va
}
