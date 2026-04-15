package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config 对应 config.yaml 的结构
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	LLM struct {
		OpenAI struct {
			APIKey  string `mapstructure:"api_key"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"openai"`
	} `mapstructure:"llm"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")   // 配置文件名 (无需后缀)
	viper.SetConfigType("yaml")     // 如果没有后缀，指定类型
	viper.AddConfigPath(".")        // 在当前目录查找
	viper.AddConfigPath("./config") // 或者在 config 目录下查找

	// 支持环境变量映射
	// 例如：export TINY_LLM_SERVER_PORT=9090
	viper.SetEnvPrefix("TINY_LLM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// 如果找不到配置文件，也可以通过环境变量启动
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	// 设置默认值
	if conf.Server.Port == "" {
		conf.Server.Port = ":8080"
	}

	return &conf, nil
}
