package logger

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// 自定义日志写入器
type dynamicLogWriter struct {
	mu         sync.Mutex
	currentDay string
	file       *os.File
	logDir     string
}

func (w *dynamicLogWriter) Write(p []byte) (n int, err error) {
	// 加锁, 让各个进程互不影响
	w.mu.Lock()
	defer w.mu.Unlock()
	//按天来存储日只能记录
	currentDay := time.Now().Format("2006-01-02")
	if w.currentDay != currentDay {
		w.file.Close()

		// 创建新的日志记录
		if err = os.MkdirAll(w.logDir, 0755); err != nil {
			fmt.Println("make directory failed", err.Error(), w.logDir)
			return 0, err
		}

		// 写到文件上
		filepath := w.logDir + "/app-" + currentDay + ".log"
		file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("open file failed, err: ", err.Error())
			return 0, err
		}

		w.file = file
		w.currentDay = currentDay
	}
	//写入日志
	return w.file.Write(p)
}

func Init(mode string) {
	// 发布模式的话不需要把日志文件写到终端上面
	var core zapcore.Core
	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	if mode == "release" {
		core = zapcore.NewCore(jsonEncoder, getLogWriter(), zapcore.DebugLevel) //文件
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(getEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel), //终端
			zapcore.NewCore(jsonEncoder, getLogWriter(), zapcore.DebugLevel),              // 文件
		)
	}
	// 创建logger
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	// 替代全局变量
	zap.ReplaceGlobals(logger)
}

// 获取Encoder
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()                             // 日志编码器配置
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") //定义日志输出时间格式
	encoderConfig.TimeKey = "Time"                                                //时间key
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                  //等级颜色设置
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder                 // 定义耗时
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder                       // 显示文件名和行号
	return zapcore.NewConsoleEncoder(encoderConfig)                               // 写到终端上面的比较好看
}

// 日志同步器 --> 写到文件和终端上
func getLogWriter() zapcore.WriteSyncer {
	writer := &dynamicLogWriter{
		logDir: viper.GetString("log.logDir"),
	}
	// 将日志文件写到文件当中去
	return zapcore.AddSync(writer)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
