package config

type Options struct {
	HTTPListenPort int
	Debug          bool
	AuditLogPath   string
	Db             DbServer
	Cache          Cache
}

type DbServer struct {
	Ip     string
	Port   int
	User   string
	Passwd string
}

type Cache struct {
	Addr   string
	Passwd string
	DB     int
}
