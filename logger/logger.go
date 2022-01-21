package logger

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Logger = log.New()
)

type CustomFormatter struct{}

//自定义日志输出格式
func (s *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}
	msg := fmt.Sprintf("%s [%s] [%s:%d] %s\n", timestamp, strings.ToUpper(entry.Level.String()), file, line, entry.Message)
	return []byte(msg), nil
}

// Setup initialize the log instance
func Setup() {
	_ = getLogFilePath()
	fileName := getLogFileName()

	writer, _ := rotateLogs.New(
		"%Y-%m-%d_%H-%M"+fileName,
		//rotateLogs.WithLinkName(logFile),
		rotateLogs.WithMaxAge(1*24*time.Hour),
		rotateLogs.WithRotationTime(1*time.Hour),
	)
	mw := io.MultiWriter(os.Stdout, writer)
	Logger.SetOutput(mw)
	//Logger.SetFormatter(&log.TextFormatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})
	Logger.SetFormatter(new(CustomFormatter))
	Logger.SetReportCaller(true)
	switch strings.ToLower("info") {
	case "panic":
		Logger.SetLevel(log.PanicLevel)
	case "fatal":
		Logger.SetLevel(log.FatalLevel)
	case "warn":
		Logger.SetLevel(log.WarnLevel)
	case "info":
		Logger.SetLevel(log.InfoLevel)
	case "debug":
		Logger.SetLevel(log.DebugLevel)
	case "trace":
		Logger.SetLevel(log.TraceLevel)
	default:
		Logger.SetLevel(log.InfoLevel)
	}

	//logger.Hooks.Add(NewContextHook())
}

func getLogFilePath() string {
	var (
		logPath                 = "logs"
		currentPath, logAbsPath string
		err                     error
	)
	if currentPath, err = filepath.Abs("."); err != nil {
		panic(err)
	}
	logAbsPath = filepath.Join(currentPath, logPath)
	if _, err := os.Stat(logAbsPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logAbsPath, 0755)
			if err != nil {
				fmt.Printf("%s 目录已存在", logAbsPath)
			}
		}
	}
	return logPath
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s.%s",
		time.Now().Format(setting.AppSetting.TimeFormat),
		"log",
	)
}
