package config

import (
	"encoding/json"
	"flag"
	"github.com/jinzhu/configor"
	log "github.com/wonderivan/logger"
	"os"
	"ticket-service/utils"
)

var Conf *Config

// 默认配置

var DefaultConfig = Config{
	APP: APP{
		Name:               "ticket-service",
		IP:                 "0.0.0.0",
		Port:               "8041",
		Mode:               "debug",
		SkipAuthentication: false,
		ContextPath:        "/api",
		UploadBasePath:     "files/any_files/",
		UploadFileSize:     10485760,
	},
	DB: DB{
		Name:            "ticket",
		Host:            "120.46.48.255",
		User:            "root",
		Password:        "root",
		Port:            "3306",
		MaxIdleConnects: 10,
		MaxOpenConnects: 1024,
		InitTable:       true,
	},
	Redis: Redis{
		Host:     "120.46.48.255",
		Port:     "6379",
		Password: "",
		//MaTeachProgressKey: "ma_teach_progress",
	},
}

type Config struct {
	APP   APP   `json:"app" yaml:"app"`
	DB    DB    `json:"db" yaml:"db"`
	Redis Redis `json:"redis" yaml:"redis"`
}

type APP struct {
	Name               string `yaml:"name" json:"name"`
	IP                 string `yaml:"ip" json:"ip"`
	Port               string `yaml:"port" json:"port"`
	Mode               string `yaml:"mode" json:"mode"`
	SkipAuthentication bool   `yaml:"skip_authentication" json:"skip_authentication"`
	ContextPath        string `yaml:"context_path" json:"context_path"`
	UploadBasePath     string `yaml:"upload_base_path" json:"upload_base_path"`
	UploadFileSize     int    `yaml:"upload_file_size" json:"upload_file_size"`
}

type DB struct {
	Name            string `yaml:"name" json:"name"`
	Host            string `yaml:"host" json:"host"`
	User            string `yaml:"user" json:"user"`
	Password        string `yaml:"password" json:"password"`
	Port            string `yaml:"port" json:"port"`
	MaxIdleConnects int    `yaml:"max_idle_connects" json:"max_idle_connects"`
	MaxOpenConnects int    `yaml:"max_open_connects" json:"max_open_connects"`
	InitTable       bool   `yaml:"init_table" json:"init_table"`
}

type Redis struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	//MaTeachProgressKey string `yaml:"ma_teach_progress_key" json:"ma_teach_progress_key"`
}

func InitConfig() error {
	Conf = &DefaultConfig
	confPath := "./conf/config.yml"
	if utils.FileExist(confPath) {
		c := initConfLoader()
		log.Debug("加载用户自定义配置...")
		err := c.Load(Conf, confPath)
		if err != nil {
			return err
		}
	}
	// 启动命令参数覆盖默认配置
	appIP := flag.String("app_ip", "", "输入app的ip地址")
	appPort := flag.String("app_port", "", "输入app的端口号")
	dbHost := flag.String("db_host", "", "输入db的ip地址")
	flag.Parse()
	if *appIP != "" {
		Conf.APP.IP = *appIP
	}
	if *appPort != "" {
		Conf.APP.Port = *appPort
	}
	if *dbHost != "" {
		Conf.DB.Host = *dbHost
	}
	//LoadConfFromEnv(Conf)
	log.Info("启动配置参数：")
	PrettyPrint(Conf)
	if !utils.Exists(Conf.APP.UploadBasePath) {
		err := os.MkdirAll(Conf.APP.UploadBasePath, 0777)
		if err != nil {
			log.Error("上传文件目录创建失败。err:[%#v]", err)
		}
	}
	return nil
}

func initConfLoader() *configor.Configor {
	config := configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
			log.Info("配置文件热加载：")
			PrettyPrint(config)
		},
	}
	c := configor.New(&config)
	return c
}

//func LoadConfFromEnv(conf *Config) {
//	log.Debug("读取环境变量配置参数.env:[%#v]", os.Environ())
//	if appIp, ok := os.LookupEnv("APP_IP"); ok {
//		conf.APP.IP = appIp
//	}
//	if appPort, ok := os.LookupEnv("APP_PORT"); ok {
//		conf.APP.Port = appPort
//	}
//	if dbHost, ok := os.LookupEnv("DB_HOST"); ok {
//		conf.DB.Host = dbHost
//	}
//	if dbPort, ok := os.LookupEnv("DB_PORT"); ok {
//		conf.DB.Port = dbPort
//	}
//	if dbInit, ok := os.LookupEnv("DB_INIT_TABLE"); ok {
//		conf.DB.InitTable = dbInit == "true"
//	}
//
//	if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok {
//		conf.Redis.Host = redisHost
//	}
//	if redisPort, ok := os.LookupEnv("REDIS_PORT"); ok {
//		conf.Redis.Port = redisPort
//	}
//}

func PrettyPrint(data interface{}) {
	p, _ := json.MarshalIndent(data, "", "\t")
	log.Info("%s \n", p)
}
