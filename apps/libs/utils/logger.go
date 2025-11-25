package utils

import (
	"go.uber.org/zap"
)

var Log *zap.Logger
var SugarLog *zap.SugaredLogger

func Initialize(environment string) error {
	var err error
	var logger *zap.Logger

	if environment == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	Log = logger
	SugarLog = logger.Sugar()

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

