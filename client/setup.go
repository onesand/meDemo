package client

import (
	"gorm.io/gorm"
	"meDemo/model"
	"sync"
)

func SetupConnections() {
	SetupConnectionsWithDBConfig(&gorm.Config{})
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
	wg.Add(2)

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

	go func() {
		defer wg.Done()
		// logrus.Info("Connecting to database...")
		if err := ConnectDBWithConfig(gormOption); err != nil {
			println("Fail to connect to database," + err.Error())
			panic(err)
		} else {
			err := DB().AutoMigrate(model.UserAddress{})
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
	wg.Wait()
}
