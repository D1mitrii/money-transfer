package main

import (
	"github.com/d1mitrii/money-transfer/bank-service/internal/app"
	"github.com/d1mitrii/money-transfer/bank-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.Run(cfg)
}
