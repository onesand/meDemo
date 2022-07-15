module meDemo

go 1.16

//go mod tidy
//heroku logs --tail  --app fringuante-mandarine-50948

require (
	github.com/ethereum/go-ethereum v1.10.20
	github.com/gin-gonic/gin v1.8.1
	gorm.io/driver/postgres v1.3.8
	gorm.io/gorm v1.23.8
)
