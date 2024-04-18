package cmds

import (
	"FinalTaskAppGoBasic/internal/api"
	"FinalTaskAppGoBasic/internal/configs"
	"FinalTaskAppGoBasic/internal/database"
	"FinalTaskAppGoBasic/internal/handlers"
	"FinalTaskAppGoBasic/internal/logs"

	"github.com/sirupsen/logrus"
)

func Cmd() {
	err := logs.Init()
	if err != nil {
		logrus.Fatalf("init logger failed: %v", err)
	}

	logs.Log.Info("logrus init successfully")

	cfg, err := configs.LoadConfig("./configs/config.yaml")
	if err != nil {
		logs.Log.Fatalf("init config failed: %v", err)
	}

	logs.Log.Info("init app config successfully")

	dbConnection := database.New()

	err = dbConnection.Connect(&cfg.DataBase)
	if err != nil {
		logs.Log.Fatalf("init connection to database failed: %v", err)
	}

	err = dbConnection.Migrate()
	if err != nil {
		logs.Log.Fatalf("migrate failed: %v", err)
	}

	apiHandlers := handlers.New(dbConnection.Gorm())

	server := api.New(&cfg.Restapi, apiHandlers)
	server.Init()

	err = server.ListenAndServe()
	if err != nil {
		logs.Log.WithError(err).Error("server closed")
	}

	err = dbConnection.Close()
	if err != nil {
		logs.Log.Fatalf("close database failed: %v", err)
	}
}
