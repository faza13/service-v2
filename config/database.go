package config

type Database struct {
	Driver          string `json:"driver"` // mariadb,postgres
	Host            string `json:"host"`
	Port            string `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Name            string `json:"database"`
	MaxIdleConn     string `json:"max_idle"`
	MaxOpenConn     string `json:"max_open"`
	MaxConnLifetime string `json:"max_conn_lifetime"`
}
