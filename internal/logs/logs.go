package logs

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Log *logrus.Logger

func Init() error {
	viper.SetConfigName("logsetup")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	return configureLogrus()
}

func configureLogrus() error {
	Log = logrus.New()

	levelString := viper.GetString("logger.logLevel")

	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return fmt.Errorf("invalid logger level %s: %v", levelString, err)
	}

	Log.SetLevel(level)
	Log.SetOutput(os.Stdout)

	logFile := viper.GetString("logger.logFile")

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open log file %s: %v", logFile, err)
	}

	Log.SetOutput(file)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	return nil
}
