package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/M1sterZag/Dont_Play_Separately/internal/core/config"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_pgx_pool "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	users_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository/postgres"
	users_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/service"
	users_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing users feature")
	usersRepo := users_postgres_repository.New(pool)
	usersService := users_service.NewUsersService(usersRepo)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
