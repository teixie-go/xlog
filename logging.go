package xlog

import (
	"container/list"
	"io"
	"os"
	"time"

	"github.com/op/go-logging"
)

var (
	_ Configurator = (*Configuration)(nil)
)

var (
	fileList  = list.New()
	logFormat = "%{time:15:04:05.000} %{shortfile} >%{level:.5s} - %{message}"
	stdFormat = "%{color}%{time:15:04:05.000} %{shortfile} >%{level:.5s}%{color:reset} - %{message}"
)

type Configurator interface {
	GetLogPath() string
	GetLogLevel() string
	GetLogFormat() string
	GetStdoutLevel() string
	GetStderrLevel() string
	GetStdFormat() string
}

type Configuration struct {
	LogPath     string `yaml:"log_path" json:"log_path"`
	LogLevel    string `yaml:"log_level" json:"log_level"`
	LogFormat   string `yaml:"log_format" json:"log_format"`
	StdoutLevel string `yaml:"stdout_level" json:"stdout_level"`
	StderrLevel string `yaml:"stderr_level" json:"stderr_level"`
	StdFormat   string `yaml:"std_format" json:"std_format"`
}

func (c *Configuration) GetLogPath() string {
	return c.LogPath
}

func (c *Configuration) GetLogLevel() string {
	return c.LogLevel
}

func (c *Configuration) GetLogFormat() string {
	return c.LogFormat
}

func (c *Configuration) GetStdoutLevel() string {
	return c.StdoutLevel
}

func (c *Configuration) GetStderrLevel() string {
	return c.StderrLevel
}

func (c *Configuration) GetStdFormat() string {
	return c.StdFormat
}

//------------------------------------------------------------------------------

func closeOldLogFile(isNewOpen bool) {
	expectedFileNum := 0
	if isNewOpen {
		expectedFileNum++
	}
	if fileList.Len() > expectedFileNum {
		element := fileList.Front()
		if element == nil {
			return
		}
		if fp, ok := element.Value.(*os.File); ok {
			fileList.Remove(element)
			time.Sleep(time.Second * 5)
			Notice("start close old log file")
			if err := fp.Close(); err != nil {
				Error("file close error: %v", err)
			}
		} else {
			Error("file type error")
		}
	}
}

func initLogging(out io.Writer, level logging.Level, formatter logging.Formatter) {
	logBackend := logging.NewLogBackend(out, "", 1)
	backendFormatter := logging.NewBackendFormatter(logBackend, formatter)
	leveledBackend := logging.AddModuleLevel(backendFormatter)
	leveledBackend.SetLevel(level, "")
	logging.SetBackend(leveledBackend)
}

func initFileLogging(path string, level logging.Level, formatter logging.Formatter) error {
	if len(path) > 0 {
		fp, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		fileList.PushBack(fp)
		initLogging(fp, level, formatter)
	}
	go closeOldLogFile(len(path) > 0)
	return nil
}

func Init(cfg Configurator) error {
	if cfg == nil {
		initLogging(os.Stdout, logging.INFO, MustStringFormatter(stdFormat))
		return nil
	}

	// 处理日志formatter
	var logFormatter, stdFormatter Formatter
	if len(cfg.GetLogFormat()) > 0 {
		logFormatter = MustStringFormatter(cfg.GetLogFormat())
	} else {
		logFormatter = MustStringFormatter(logFormat)
	}
	if len(cfg.GetStdFormat()) > 0 {
		stdFormatter = MustStringFormatter(cfg.GetStdFormat())
	} else {
		stdFormatter = MustStringFormatter(stdFormat)
	}

	// 判断是否输出stderr日志
	if len(cfg.GetStderrLevel()) > 0 {
		stderrLevel, err := logging.LogLevel(cfg.GetStderrLevel())
		if err != nil {
			return err
		}
		initLogging(os.Stderr, stderrLevel, stdFormatter)
	}

	logPath := cfg.GetLogPath()
	logLevel, err := logging.LogLevel(cfg.GetLogLevel())
	if err != nil {
		logLevel = logging.INFO
	}

	// 若设置了stdout level或日志文件路径为空，则输出stdout日志，且优先stdout level
	if len(cfg.GetStdoutLevel()) > 0 {
		stdoutLevel, err := logging.LogLevel(cfg.GetStdoutLevel())
		if err != nil {
			return err
		}
		initLogging(os.Stdout, stdoutLevel, stdFormatter)
	} else if len(logPath) <= 0 {
		initLogging(os.Stdout, logLevel, stdFormatter)
	}

	return initFileLogging(logPath, logLevel, logFormatter)
}
