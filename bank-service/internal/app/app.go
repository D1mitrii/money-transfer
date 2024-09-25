package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	grpcapp "github.com/d1mitrii/money-transfer/bank-service/internal/app/grpc"
	"github.com/d1mitrii/money-transfer/bank-service/internal/config"
	"github.com/d1mitrii/money-transfer/bank-service/internal/repository/pgdb"
	"github.com/d1mitrii/money-transfer/bank-service/internal/services/bank"
	"github.com/d1mitrii/money-transfer/bank-service/pkg/logger"
	"github.com/d1mitrii/money-transfer/bank-service/pkg/postgres"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *config.Config) {
	const op = "app - Run"

	// Logger
	log := logger.SetupLogger("local")
	log.Info(cfg.Postgres.URL)

	// Database
	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.MaxPoolSize))
	if err != nil {
		log.Error(fmt.Sprintf("%s - postgres.New: %v", op, err))
		return
	}
	defer pg.Close()

	// Repositories
	bankRepo := pgdb.New(pg)

	// Services
	b := bank.New(
		log,
		bankRepo,
		bankRepo,
	)

	// grpc server
	grpcApp := grpcapp.New(log, b, cfg.GRPC.Port)

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer done()

	// Run apps
	g, _ := errgroup.WithContext(ctx)
	g.Go(grpcApp.Run)

	g.Wait()

	// Graceful shutdown
	grpcApp.Stop()
}
