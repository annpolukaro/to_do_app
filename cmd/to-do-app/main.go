package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/annpolukaro/to_do_app/internal/core/logger"
	core_postgres_pool "github.com/annpolukaro/to_do_app/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/annpolukaro/to_do_app/internal/core/transport/http/middleware"

	core_http_server "github.com/annpolukaro/to_do_app/internal/core/transport/http/server"
	users_postgres_repository "github.com/annpolukaro/to_do_app/internal/features/users/repository/postgres"
	users_service "github.com/annpolukaro/to_do_app/internal/features/users/service"
	user_transport_http "github.com/annpolukaro/to_do_app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer cancel()
	fmt.Println("Hello, ToDo App")
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initialing postgres connection postgres connection pool")

	pool, err := core_postgres_pool.NewConnectionPoll(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}

	defer pool.Close()

	logger.Debug("initialazing feature", zap.String("feature", "user"))
	usersRepository := users_postgres_repository.NewuserRepository(pool)
	userService := users_service.NewUsersService(usersRepository)

	logger.Debug("Starting ToDo application!")

	usersTransportHTTP := user_transport_http.NewUsersHTTPHandler(userService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
