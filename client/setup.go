package client

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"meDemo/model"
	"sync"
)

func SetupConnections() {
	SetupConnectionsWithDBConfig(&gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
}

func SetUpEthClient() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// logrus.Info("Connecting to ethereum node...")
		if err := ConnectEthNode(); err != nil {
			println("Fail to connect to ethereum node," + err.Error())
			panic(err)
		} else {
			println("Connected to ethereum node")
		}
	}()
	wg.Wait()
}

func SetupConnectionsWithDBConfig(gormOption gorm.Option) {
	var wg sync.WaitGroup
	wg.Add(3)

	//go func() {
	//	defer wg.Done()
	//	// logrus.Info("Connecting to redis...")
	//	if err := ConnectRedis(); err != nil {
	//		logrus.WithError(err).Error("Fail to connect to redis")
	//		panic(err)
	//	} else {
	//		logrus.Info("Connected to redis")
	//	}
	//}()

	//go func() {
	//	defer wg.Done()
	//	if err := ConnectionBot(); err != nil {
	//		println("Fail to connect to dc bot," + err.Error())
	//		panic(err)
	//	} else {
	//		println("Connected to dc bot")
	//	}
	//}()

	go func() {
		defer wg.Done()
		// logrus.Info("Connecting to database...")
		if err := ConnectDBWithConfig(gormOption); err != nil {
			println("Fail to connect to database," + err.Error())
			panic(err)
		} else {
			err := DB().AutoMigrate(model.UserAddress{})
			err = DB().AutoMigrate(model.FreeMintMode{})
			err = DB().AutoMigrate(model.Mints{})
			err = DB().AutoMigrate(model.Nft{})
			if err != nil {
				return
			} else {
				println("UserAddress tablet created successful.")
			}
			println("Connected to database")
		}
	}()

	go func() {
		defer wg.Done()
		// logrus.Info("Connecting to ethereum node...")
		if err := ConnectEthNode(); err != nil {
			println("Fail to connect to ethereum node," + err.Error())
			panic(err)
		} else {
			println("Connected to ethereum node")
		}
	}()

	go func() {
		defer wg.Done()
		// logrus.Info("Connecting to ethereum node...")
		if err := ConnectEthWsNode(); err != nil {
			println("Fail to connect to ethereum ws node," + err.Error())
			panic(err)
		} else {
			println("Connected to ethereum ws node")
		}
	}()

	wg.Wait()
}
