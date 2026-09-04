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
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	core_storage_minio "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage/s3/minio"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	auth_config "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth"
	auth_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/repository/postgres"
	auth_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/service"
	auth_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/transport/http"
	users_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository/postgres"
	users_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/service"
	users_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/transport/http"
	"github.com/google/uuid"
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

	logger.Debug("initializing s3 storage")
	S3StorageConfig := core_storage.NewConfigMust()
	S3Storage, err := core_storage_minio.NewStorage(ctx, S3StorageConfig)
	if err != nil {
		logger.Fatal("failed to init minio storage", zap.Error(err))
	}

	logger.Debug("initializing auth feature")
	authConfig := auth_config.NewConfigMust()
	jwtSigner := auth_service.NewJWTSigner(
		authConfig.JWTSecret,
		authConfig.JWTAccessTTL,
		authConfig.JWTRefreshTTL,
	)
	authMW := core_http_middleware.Auth(func(token string) (uuid.UUID, error) {
		claims, err := jwtSigner.ParseAccessToken(token)
		if err != nil {
			return uuid.Nil, err
		}

		return uuid.Parse(claims.Subject)
	})
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository, jwtSigner)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	logger.Debug("initializing users feature")
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, S3Storage)
	usersRoutes := usersTransportHTTP.Routes()
	for i := range usersRoutes {
		usersRoutes[i].Middleware = append(usersRoutes[i].Middleware, authMW)
	}

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
	apiVersionRouter.RegisterRouters(authTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(usersRoutes...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
