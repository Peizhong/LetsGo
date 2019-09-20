package config

import (
	"encoding/json"
	"github.com/go-ini/ini"
	"github.com/peizhong/letsgo/pkg/log"
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

const defaultSettingFile = "c:/users/peizhong/source/repos/letsgo/playground/conf/appsettings.json"

var (
	// json
	once     sync.Once
	settings atomic.Value

	// ini
	HomeDir string
	WorkspaceDir string

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

	HTTPApp1Port int
	GrpcApp1Port int

	CertCrt string
	CertKey string
	CertName string
)

func init() {
	var err error
	HomeDir, err = os.UserHomeDir()
	WorkspaceDir = filepath.Join(HomeDir, "source/repos/letsgo")
	if err != nil {
		log.Fatalf("Fail to load user home dir: %v", err)
	}
	cfg, err = ini.Load(filepath.Join(WorkspaceDir, filepath.ToSlash("playground/conf/app.ini")))
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

	GatewayPort = sec.Key("GATEWAY_PORT").MustInt(8010)
	HTTPApp1Port = sec.Key("HTTP_APP1_PORT").MustInt(8081)
	GrpcApp1Port = sec.Key("GRPC_APP1_PORT").MustInt(8091)

	CertName = sec.Key("CERT_NAME").MustString("LetsGo")
	CertCrt = sec.Key("CERT_CRT").MustString("source/repos/letsgo/playground/key/server.crt")
	CertKey = sec.Key("CERT_KEY").MustString("source/repos/letsgo/playground/key/server.key")
}

// GetConnectionString get connection string from appsettings.json
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
