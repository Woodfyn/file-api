package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Woodfyn/file-api/internal/config"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/redis"
	"github.com/Woodfyn/file-api/internal/service"
	"github.com/Woodfyn/file-api/internal/transport/rest"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
	"github.com/Woodfyn/file-api/pkg/mdb"
	"github.com/Woodfyn/file-api/pkg/rdb"
	"github.com/Woodfyn/file-api/pkg/signaler"
	server "github.com/Woodfyn/file-api/pkg/srv"
)

const (
	CFG_FOLDER = "configs"
	CFG_FILE   = "prod"
)

func init() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(h))
}

func Run() {
	// init config
	cfg, err := config.InitConfig(CFG_FOLDER, CFG_FILE)
	if err != nil {
		panic(err)
	}

	log.Print(cfg)

	rdbClient := rdb.NewRdbClient(rdb.ConnInfo{})
	defer rdbClient.Close()

	mongoClient, err := mdb.NewMongoClient(context.Background(), mdb.ConnInfo{})
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	hasher := hash.NewSHA1Hasher("salt")
	tokenManager, err := auth.NewManager("secret")
	if err != nil {
		panic(err)
	}

	// init dependencies
	redis := redis.NewRepository(rdbClient)
	mongo := mongo.NewRepository(mongoClient)

	deps := service.Deps{
		RedisRepos: redis,
		MongoRepo:  mongo,
		Hasher:     hasher,

		TokenManager: tokenManager,
		AcssTokenTTL: 5 * time.Minute,
		RefreshTTL:   60 * time.Minute,
	}

	service := service.NewService(deps)
	handler := rest.NewHandler(service, tokenManager)

	go func() {
		if err := server.Run(cfg.Server.Port, handler.Init()); err != nil {
			panic(err)
		}
	}()

	slog.Info("server started...")

	signaler.Wait()

	slog.Info("server stopped...")

	server.Shutdown(handler.Init())
}
