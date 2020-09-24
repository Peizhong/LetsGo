package config

type DBConfig struct {
	UserName, Password string
	Address            string
	Database           string
	IdleConnections    int
	MaxConnections     int
}

var db = &DBConfig{
	UserName:        "app_user",
	Password:        "app_user",
	Database:        "kzz",
	Address:         "localhost:13306",
	IdleConnections: 2,
	MaxConnections:  5,
}

func DB() *DBConfig {
	return db
}
