package config

import (
	"github.com/ToolPackage/fse/common/utils"
	"os"
)

type Config struct {
	RunMode   string // FSE_RunMode: server的运行模式，dev 或 prod，默认是dev
	FileDir   string // FSE_FileDir: 上传的文件的存储路径
	Host      string // FSE_Host
	Port      string // FSE_Port
	MongoHost string
	MongoPort string
}

const (
	ModeDev  = "dev"
	ModeProd = "prod"

	MongoDbName      = "fse"
	FileInfoMongoCol = "fileInfo"
)

var (
	Conf = New()
)

func New() *Config {
	runMode := getEnvOrDefault("FSE_RunMode", ModeDev)
	fileDir := getEnvOrDefault("FSE_FileDir", "upload/")
	host := getEnvOrDefault("FSE_Host", "0.0.0.0")
	port := getEnvOrDefault("FSE_Port", "8000")
	config := &Config{
		RunMode:   runMode,
		FileDir:   fileDir,
		Host:      host,
		Port:      port,
		MongoHost: "localhost", // mongo for docker run
		MongoPort: "27017",
	}
	return config
}

func getEnvOrDefault(envName, defaultValue string) string {
	return utils.OrString(os.Getenv(envName), defaultValue)
}
