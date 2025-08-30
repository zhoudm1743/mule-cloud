package database

import (
	"errors"
	"mule-cloud/pkg/services/config"

	"github.com/zhoudm1743/torm/db"
)

func InitTorm() error {
	cfg := config.GetConfig()
	config := &db.Config{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		Charset:  cfg.Database.Charset,
	}

	err := db.AddConnection("default", config)
	if err != nil {
		return err
	}

	inst, err := db.DB("default")
	if err != nil {
		return err
	}
	if inst == nil {
		return errors.New("数据库连接失败")
	}

	return nil
}
