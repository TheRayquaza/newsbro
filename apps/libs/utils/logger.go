package utils

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Initialize(environment string) error {
	var err error
	
	if environment == "dev" {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}
	
	if err != nil {
		return err
	}
	
	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func WithContext(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}

func WithRequestID(requestID string) *zap.Logger {
	return Log.With(zap.String("request_id", requestID))
}

func WithUserID(userID uint) *zap.Logger {
	return Log.With(zap.Uint("user_id", userID))
}
