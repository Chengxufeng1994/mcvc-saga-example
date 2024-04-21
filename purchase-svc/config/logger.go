package config

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	ContextLogger *log.Entry
	GormLogger    gormlogger.Interface
}

// ginOnce is a wrapper around gin global var changes. This is a workaround
// against the lack of concurrency safety of these vars in the gin package.
var ginOnce sync.Once

var gormOnce sync.Once

// initLogger creates the logger instance
func InitLogger(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *Logger {
	writer := os.Stderr

	logger := log.New()
	logger.Out = writer
	logger.Formatter = &log.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	// logger.ReportCaller = true
	level, err := logrus.ParseLevel(appCfg.LogConfig.Level)
	if err != nil {
		log.WithError(err).Error("parsing log level error")
		logger.Level = log.InfoLevel
	}
	logger.Level = level

	contextLogger := log.NewEntry(logger).WithField("application", bootCfg.Application)

	ginOnce.Do(func() {
		if bootCfg.GinMode == "debug" {
			gin.DefaultWriter = logger.Writer()
			gin.SetMode(gin.DebugMode)
		} else {
			gin.DefaultWriter = io.Discard
			gin.SetMode(gin.ReleaseMode)
		}
	})

	var gormLogger gormlogger.Interface
	gormOnce.Do(func() {
		var gormLogLevel gormlogger.LogLevel
		if bootCfg.GinMode == "debug" {
			gormLogLevel = gormlogger.Info
		} else {
			gormLogLevel = gormlogger.Silent
		}

		gormLogger = gormlogger.New(logger,
			gormlogger.Config{
				SlowThreshold: time.Second,  // Slow SQL threshold
				LogLevel:      gormLogLevel, // Log level
				Colorful:      true,         // Disable color
				// IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
				// ParameterizedQueries:      true,         // Don't include params in the SQL log
			})
	})

	return &Logger{ContextLogger: contextLogger, GormLogger: gormLogger}
}
