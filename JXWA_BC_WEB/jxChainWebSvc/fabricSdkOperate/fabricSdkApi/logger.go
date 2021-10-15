/*
 * 资源管理实体相关的API
 */

package fabricSdkApi

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

//在使用SDK的相关API之前一定要先初始化logger
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	zapLogger := zap.New(core)
	logger = zapLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	//如果想要追加写入可以查看我的博客文件操作那一章
	file, err := os.Create("./test.log")
	if err != nil {
		panic("create file failed!")
	}
	return zapcore.AddSync(file)
}
