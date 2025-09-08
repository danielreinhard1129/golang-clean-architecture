package configs

import (
	"os"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/joho/godotenv"
)

type Config interface {
	Get(key string) string
}

type configImpl struct{}

func (config *configImpl) Get(key string) string {
	return os.Getenv(key)
}

func New(filenames ...string) Config {
	err := godotenv.Load(filenames...)
	exception.PanicLogging(err)
	return &configImpl{}
}
