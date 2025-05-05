package main

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志配置
func initLogger() *zap.Logger {
	// 1. 配置编码器
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder   // 可读的时间格式
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder // 大写日志级别

	// 2. 配置输出位置（文件 + 控制台）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,   // MB
		MaxBackups: 3,    // 保留旧文件数量
		MaxAge:     30,   // 保留天数
		Compress:   true, // 启用压缩
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	// 3. 创建多输出
	multiWriter := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)

	// 4. 创建Core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg), // 控制台友好格式
		multiWriter,
		zap.DebugLevel, // 日志级别
	)

	// 5. 创建Logger并添加调用信息
	logger := zap.New(core,
		zap.AddCaller(),                   // 记录调用位置
		zap.AddStacktrace(zap.ErrorLevel), // 错误级别记录堆栈
	)

	return logger
}

// HTTP日志中间件
func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 包装ResponseWriter获取状态码
			wrapped := &responseWrapper{w, 0}
			next.ServeHTTP(wrapped, r)

			logger.Info("HTTP请求",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.status),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

// 自定义ResponseWriter包装
type responseWrapper struct {
	http.ResponseWriter
	status int
}

func (w *responseWrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func main() {
	// 初始化日志
	logger := initLogger()
	defer logger.Sync()

	// 设置全局Logger
	zap.ReplaceGlobals(logger)

	// 使用示例
	logger.Debug("调试信息",
		zap.String("service", "startup"),
		zap.Int("version", 1),
	)

	logger.Info("服务初始化完成",
		zap.String("environment", "production"),
	)

	logger.Warn("数据库连接池接近限制",
		zap.Int("current_connections", 95),
		zap.Int("max_connections", 100),
	)

	// 使用SugaredLogger
	sugar := logger.Sugar()
	sugar.Infow("订单处理",
		"orderID", 12345,
		"amount", 199.99,
		"currency", "USD",
	)

	// 启动HTTP服务
	http.Handle("/", loggingMiddleware(logger)(http.HandlerFunc(handler)))
	logger.Info("启动HTTP服务", zap.String("addr", ":8080"))
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("处理请求开始") // 使用全局Logger

	// 模拟业务处理
	time.Sleep(100 * time.Millisecond)

	w.Write([]byte("OK"))
	zap.L().Info("处理请求完成")
}
