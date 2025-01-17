package web

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Web{
	MaxSize: 1024,
	Except: Route{
		Uri:    "",
		Method: "",
	},
	System: System{
		Tls:        false,
		Level:      "debug",
		Addr:       "127.0.0.1:8085",
		DbType:     "mysql",
		TimeFormat: "2006-01-02 15:04:05",
	},
	Limit: Limit{
		Disable: true,
		Limit:   0,
		Burst:   5,
	},
	Captcha: Captcha{
		KeyLong:   4,
		ImgWidth:  240,
		ImgHeight: 80,
	},
}

type Web struct {
	MaxSize int64   `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	Except  Route   `mapstructure:"except" json:"except" yaml:"except"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit   `mapstructure:"limit" json:"limit" yaml:"limit"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}
type Route struct {
	Uri    string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Method string `mapstructure:"method" json:"method" yaml:"method"`
}

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
}

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}

type System struct {
	Tls          bool   `mapstructure:"tls" json:"tls" yaml:"tls"`       // debug,release,test
	Level        string `mapstructure:"level" json:"level" yaml:"level"` // debug,release,test
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`
	StaticPrefix string `mapstructure:"static-prefix" json:"staticPrefix" yaml:"static-prefix"`
	WebPrefix    string `mapstructure:"web-prefix" json:"webPrefix" yaml:"web-prefix"`
	DbType       string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	TimeFormat   string `mapstructure:"time-format" json:"timeFormat" yaml:"time-format"`
}

// Verfiy
func Verfiy() {
	if CONFIG.System.Addr == "" { // 默认 8085
		CONFIG.System.Addr = "127.0.0.1:8085"
	}

	if CONFIG.System.TimeFormat == "" { // 默认 80
		CONFIG.System.TimeFormat = "2006-01-02 15:04:05"
	}
}

// ToStaticUrl 静态地址url
func ToStaticUrl(uri string) string {
	path := filepath.Join(CONFIG.System.Addr, CONFIG.System.StaticPrefix, uri)
	if CONFIG.System.Tls {
		return filepath.ToSlash(str.Join("https://", path))
	}
	return filepath.ToSlash(str.Join("http://", path))
}

// IsExist 配置文件是否存在
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove 删除配置文件
func Remove() error {
	err := getViperConfig().Remove()
	if err != nil {
		return err
	}
	return nil
}

// getViperConfig 获取初始化配置
func getViperConfig() viper_server.ViperConfig {
	maxSize := strconv.FormatInt(CONFIG.MaxSize, 10)
	keyLong := strconv.FormatInt(int64(CONFIG.Captcha.KeyLong), 10)
	imgWidth := strconv.FormatInt(int64(CONFIG.Captcha.ImgWidth), 10)
	imgHeight := strconv.FormatInt(int64(CONFIG.Captcha.ImgHeight), 10)
	limit := strconv.FormatInt(int64(CONFIG.Limit.Limit), 10)
	burst := strconv.FormatInt(int64(CONFIG.Limit.Burst), 10)
	disable := strconv.FormatBool(CONFIG.Limit.Disable)
	tls := strconv.FormatBool(CONFIG.System.Tls)
	configName := "web"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("反序列化错误: %v", err)
			}
			// 监控配置文件变化
			vi.SetConfigName(configName)
			vi.WatchConfig()
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("配置发生变化:", e.Name)
				if err := vi.Unmarshal(&CONFIG); err != nil {
					fmt.Printf("反序列化错误: %v \n", err)
				}
			})
			return nil
		},
		// 注意:设置默认配置值的时候,前面不能有空格等其他符号.必须紧贴左侧.
		Default: []byte(`
{
	"max-size": ` + maxSize + `,
	"except":
		{ 
			"uri": "` + CONFIG.Except.Uri + `",
			"method": "` + CONFIG.Except.Method + `"
		},
	"captcha":
		{
		"key-long": ` + keyLong + `,
		"img-width": ` + imgWidth + `,
		"img-height": ` + imgHeight + `
		},
	"limit":
		{
			"limit": ` + limit + `,
			"disable": ` + disable + `,
			"burst": ` + burst + `
		},
	"system":
		{
			"tls": ` + tls + `,
			"level": "` + CONFIG.System.Level + `",
			"addr": "` + CONFIG.System.Addr + `",
			"db-type": "` + CONFIG.System.DbType + `",
			"time-format": "` + CONFIG.System.TimeFormat + `"
		}
 }`),
	}
}
