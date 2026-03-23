package zapLogger

import (
	"gin_demo/internal/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() (zLogger *zap.Logger, zSugarLogger *zap.SugaredLogger) {
	var logger *zap.Logger
	var sugarLogger *zap.SugaredLogger
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	// NewCore创建一个向WriteSyncer写入日志的Core。

	// A WriteSyncer is an io.Writer that can also flush any buffered data. Note
	// that *os.File (and thus, os.Stderr and os.Stdout) implement WriteSyncer.

	// LevelEnabler决定在记录消息时是否启用给定的日志级别。
	// Each concrete Level value implements a static LevelEnabler which returns
	// true for itself and all higher logging levels. For example WarnLevel.Enabled()
	// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
	// FatalLevel, but return false for InfoLevel and DebugLevel.
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// New constructs a new Logger from the provided zap core.Core and Options. If
	// the passed zap core.Core is nil, it falls back to using a no-op
	// implementation.

	// AddCaller configures the Logger to annotate each message with the filename,
	// line number, and function name of zap's caller. See also WithCaller.
	logger = zap.New(core, zap.AddCaller())
	// Sugar封装了Logger，以提供更符合人体工程学的API，但速度略慢。糖化一个Logger的成本非常低，
	// 因此一个应用程序同时使用Loggers和SugaredLoggers是合理的，在性能敏感代码的边界上在它们之间进行转换。
	sugarLogger = logger.Sugar()
	return logger, sugarLogger
}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next() // 执行后续中间件

		// Since returns the time elapsed since t.
		// It is shorthand for time.Now().Sub(t).
		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost), // 运行时间
		)
	}
}

func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
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

func getEncoder() zapcore.Encoder {
	// NewJSONEncoder创建了一个快速、低分配的JSON编码器。编码器适当地转义所有字段键和值。
	// NewProductionEncoderConfig returns an opinionated EncoderConfig for
	// production environments.
	//return zap core.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// NewConsoleEncoder创建一个编码器，其输出是为人类而不是机器设计的。
	// 它以纯文本格式序列化核心日志条目数据(消息、级别、时间戳等)，并将结构化上下文保留为JSON。
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	// Logger is an io.WriteCloser that writes to the specified filename.
	// 日志记录器在第一次写入时打开或创建日志文件。如果文件存在并且小于MaxSize兆字节，则lumberjack将打开并追加该文件。
	// 如果该文件存在并且其大小为>= MaxSize兆字节，
	// 则通过将当前时间放在文件扩展名(或者如果没有扩展名则放在文件名的末尾)的名称中的时间戳中来重命名该文件。
	// 然后使用原始文件名创建一个新的日志文件。
	// 每当写操作导致当前日志文件超过MaxSize兆字节时，将关闭当前文件，重新命名，并使用原始名称创建新的日志文件。
	// 因此，您给Logger的文件名始终是“当前”日志文件。
	// 如果MaxBackups和MaxAge均为0，则不会删除旧的日志文件。
	lumberJackLogger := &lumberjack.Logger{
		// Filename是要写入日志的文件。备份日志文件将保留在同一目录下
		Filename: config.GetLogPath() + "zapLog" + time.Now().Format("2006-01-02") + ".log",
		// MaxSize是日志文件旋转之前的最大大小(以兆字节为单位)。默认为100兆字节。
		MaxSize: 1, // M
		// MaxBackups是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件(尽管MaxAge仍然可能导致它们被删除)。
		MaxBackups: 5, // 备份数量
		// MaxAge是根据文件名中编码的时间戳保留旧日志文件的最大天数。
		// 请注意，一天被定义为24小时，由于夏令时、闰秒等原因，可能与日历日不完全对应。默认情况下，不根据时间删除旧的日志文件。
		MaxAge: 7, // 备份天数
		// Compress决定是否应该使用gzip压缩旋转的日志文件。默认情况下不执行压缩。
		Compress:  true, //是否压缩
		LocalTime: false,
	}
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),        //日志同时输出到控制台
		zapcore.AddSync(lumberJackLogger), //配置的hook
	)
}
