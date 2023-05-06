package configloader

import (
	"gopkg.in/yaml.v3"
)

// 支持yaml格式的配置文件
func init() {
	RegisterExtensionHandler("yaml", yamlLoader)
	RegisterExtensionHandler("yml", yamlLoader)
}

func yamlLoader(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
