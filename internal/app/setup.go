package app

import (
	"order-service/internal/config"
	"order-service/internal/handlers"
	"order-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)


func SetRoutes(router *gin.Engine, handler *handlers.Handler) {
	router.POST("/orders", handler.CreateOrder)
	router.GET("/order-items", handler.GetOrderItems)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func SetupLogger(cfg config.LoggingConfig) (logger.Logger, error) {
	lvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	var sampling *zap.SamplingConfig
	if cfg.Sampling != nil {
		sampling = &zap.SamplingConfig{
			Initial:    cfg.Sampling.Initial,
			Thereafter: cfg.Sampling.Thereafter,
		}
	}

	opts := []logger.Option{
		logger.WithMode(cfg.Mode),
		logger.WithLevel(lvl),
		logger.WithEncoding(cfg.Encoding),
		logger.WithDisableCaller(cfg.DisableCaller),
		logger.WithDisableStacktrace(cfg.DisableStacktrace),
		logger.WithOutputPaths(cfg.OutputPaths...),
		logger.WithErrorOutputPaths(cfg.ErrorOutputPaths...),
		logger.WithEncoderConfig(func(ec *zapcore.EncoderConfig) {
			ec.TimeKey = cfg.TimestampKey
			ec.MessageKey = "M"
			ec.LevelKey = "L"
			ec.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
			if cfg.CapitalizeLevel {
				ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
			}
		}),
	}

	if sampling != nil {
		opts = append(opts, logger.WithSampling(sampling))
	}

	if len(cfg.InitialFields) > 0 {
		opts = append(opts, logger.WithInitialFields(cfg.InitialFields))
	}

	return logger.NewLogger(opts...)
}
