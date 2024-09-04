package zlog

import (
	"encoding/json"

	"go.uber.org/zap"
)

// Info ...
func Info(content string, data LogData) {
	jsonData, _ := json.Marshal(data)
	zapLogger.Info(content, zap.String("data", string(jsonData)))
}

// Error ...
func Error(content string, data LogData) {
	jsonData, _ := json.Marshal(data)
	zapLogger.Error(content, zap.String("data", string(jsonData)))
}
