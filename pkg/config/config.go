package config

import (
	"encoding/json"
	"github.com/go-ini/ini"
	"letsgo/pkg/log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// todo: 文件监视: watch file, delivering events to a channel

type logging struct {
	LogLevel map[string]string
}

type Appsettings struct {
	Logging           logging
	ConnectionStrings map[string]string
}

const defaultSettingFile = "c:/users/peizhong/source/repos/letsgo/config/appsettings.json"

var (
	// json
	once     sync.Once
	settings atomic.Value

	// ini
	cfg       *ini.File
	RunMode   string
	JwtSecret string

	// database
	DBType     string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string

	// server
	GatewayPort int
)

func init() {
	var err error
	cfg, err = ini.Load("c:/users/peizhong/source/repos/letsgo/config/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	loadBase()
	loadDatabase()
	loadApp()
}

func loadBase() {
	RunMode = cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func loadDatabase() {
	DBType = cfg.Section("database").Key("TYPE").MustString("sqlite3")
	DBName = cfg.Section("database").Key("NAME").MustString("dev")
	DBHost = cfg.Section("database").Key("HOST").MustString("localhost")
	DBUser = cfg.Section("database").Key("USER").MustString("root")
	DBPassword = cfg.Section("database").Key("PASSWORD").MustString("root")
}

func loadApp() {
	sec, err := cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")

	GatewayPort = cfg.Section("server").Key("GATEWAY_PORT").MustInt(8010)
}

// GetConnectionString get connection string from appsettings.json
func GetConnectionString(key string) string {
	s := GetAppsettings(defaultSettingFile)
	value := s.ConnectionStrings[key]
	return value
}

func GetAppsettings(settingFile string) (s Appsettings) {
	once.Do(func() {
		if settingFile == "" {
			settingFile = defaultSettingFile
		}
		settingFile = filepath.ToSlash(settingFile)
		file, err := os.Open(settingFile)
		// 函数返回时，关闭文件
		defer file.Close()
		if err == nil {
			s := Appsettings{}
			err := json.NewDecoder(file).Decode(&s)
			if err == nil {
				settings.Store(s)
			}
		}
	})
	return settings.Load().(Appsettings)
}
