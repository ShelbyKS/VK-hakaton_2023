package app

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	auth_http "hakaton/auth/delivery/http"
	auth_postgres "hakaton/auth/repository/postgresql"
	auth_redis "hakaton/auth/repository/redis"
	auth_usecase "hakaton/auth/usecase"
	"hakaton/domain"

	"hakaton/db/connector/postgres"
	"hakaton/db/connector/redis"
	logs "hakaton/logger"
	"hakaton/middleware"

	csrf_http "hakaton/csrf/delivery/http"

	_ "github.com/lib/pq"
)

func StartServer() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	pc := postgres.Connect(ctx)
	defer pc.Close()

	rc := redis.Connect()
	defer rc.Close()

	tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sessionRepository := auth_redis.NewSessionRedisRepository(rc)
	authRepository := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	authUsecase := auth_usecase.NewAuthUsecase(authRepository, sessionRepository)

	auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	csrf_http.NewCsrfHandler(mainRouter, tokens)

	mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
