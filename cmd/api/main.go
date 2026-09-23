package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_redis_cache "github.com/M1sterZag/Dont_Play_Separately/internal/core/cache/redis"
	core_config "github.com/M1sterZag/Dont_Play_Separately/internal/core/config"
	core_logger "github.com/M1sterZag/Dont_Play_Separately/internal/core/logger"
	core_igdb_provider "github.com/M1sterZag/Dont_Play_Separately/internal/core/provider/igdb"
	core_pgx_pool "github.com/M1sterZag/Dont_Play_Separately/internal/core/repository/postgres/pool/pgx"
	core_storage "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage"
	core_storage_minio "github.com/M1sterZag/Dont_Play_Separately/internal/core/storage/s3/minio"
	core_http_middleware "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/http/middleware"
	core_http_server "github.com/M1sterZag/Dont_Play_Separately/internal/core/transport/server"
	auth_config "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth"
	auth_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/repository/postgres"
	auth_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/service"
	auth_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/auth/transport/http"
	games_config "github.com/M1sterZag/Dont_Play_Separately/internal/features/games"
	games_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/games/repository/postgres"
	games_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/games/service"
	games_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/games/transport/http"
	platforms_config "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms"
	platforms_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms/repository/postgres"
	platforms_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms/service"
	platforms_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/platforms/transport/http"
	users_postgres_repository "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/repository/postgres"
	users_service "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/service"
	users_transport_http "github.com/M1sterZag/Dont_Play_Separately/internal/features/users/transport/http"
	"github.com/google/uuid"
	"go.uber.org/zap"

	_ "github.com/M1sterZag/Dont_Play_Separately/docs"
)

// @title DPS Backend API
// @version 1.0
// @description DPS Backend REST-API Schema
// @host 127.0.0.1:5050
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите access-токен в формате "Bearer <token>"
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

	logger.Debug("initializing cache")
	CacheConfig := core_redis_cache.NewConfigMust()
	RedisClient, err := core_redis_cache.NewRedisCache(ctx, CacheConfig)
	if err != nil {
		logger.Fatal("failed to init redis cache", zap.Error(err))
	}
	defer RedisClient.Close()

	logger.Debug("initializing provider")
	ProviderConfig := core_igdb_provider.NewConfigMust()
	ProviderClient := core_igdb_provider.NewClient(ProviderConfig)

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

	logger.Debug("initializing games feature")
	gamesConfig := games_config.NewConfigMust()
	gamesRepository := games_postgres_repository.NewGamesRepository(pool)
	gamesService := games_service.NewGameService(gamesRepository, RedisClient, ProviderClient, gamesConfig.CacheTTL)
	gamesTransportHTTP := games_transport_http.NewGamesHTTPHandler(gamesService)

	logger.Debug("initializing platforms feature")
	platformsConfig := platforms_config.NewConfigMust()
	platformsRepository := platforms_postgres_repository.NewPlatformsRepository(pool)
	platformsService := platforms_service.NewPlatformsService(platformsRepository, RedisClient, ProviderClient, platformsConfig.CacheTTL)
	platformsTransportHTTP := platforms_transport_http.NewPlatformsHTTPHandler(platformsService)

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
	apiVersionRouter.RegisterRouters(gamesTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(platformsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
