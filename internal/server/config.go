package server

import (
	"fmt"
	"time"
)

// Config contains information necessary to set up a Server.
type Config struct {
	Port         int           `envvar:"PORT,default=8080"`
	ReadTimeout  time.Duration `envvar:"READ_TIMEOUT,default=3s"`
	WriteTimeout time.Duration `envvar:"WRITE_TIMEOUT,default=3s"`

	PicLife time.Duration `envvar:"PICTURE_LIFE,default=9s"`

	Logger Logger `envvar:"-"`

	Picker PicPicker
}

// The Logger interface defines the methods required by the Server for logging.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type PicPicker interface {
	Next() ([]byte, error)
}

// CfgBuildInit initializes the Logger.  It should only be called by a cfgbuild.Builder.
func (cfg *Config) CfgBuildInit() error {
	if cfg.Logger == nil {
		cfg.Logger = defaultLogger{}
	}
	return nil
}

// CfgBuildValidate checks the specified values.  It should only be called by a cfgbuild.Builder.
func (cfg *Config) CfgBuildValidate() error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("%d is not a valid port", cfg.Port)
	}
	return nil
}
