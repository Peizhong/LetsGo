package config

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/go-ini/ini"
	"github.com/peizhong/letsgo/pkg/log"
)

// todo: 文件监视: watch file, delivering events to a channel

var (
	// json
	once     sync.Once
	settings atomic.Value

	// ini
	HomeDir      string
	WorkspaceDir string

	cfg       *ini.File
	RunMode   string
	JwtSecret string

	// database
	DBType      string
	DBHost      string
	DBName      string
	DBUser      string
	DBPassword  string
	MongoDBHost string

	// server
	GatewayPort int

	HTTPApp1Port int
	GrpcApp1Port int

	CertCrt  string
	CertKey  string
	CertName string
)

func init() {
	var err error
	HomeDir, err = os.UserHomeDir()
	WorkspaceDir = filepath.Join("", "/home/peizhong/source/repos/letsgo")
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
	MongoDBHost = cfg.Section("database").Key("MONGOHOST").MustString("localhost")
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
