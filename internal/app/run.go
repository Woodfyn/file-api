package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Woodfyn/file-api/internal/config"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/redis"
	"github.com/Woodfyn/file-api/internal/repository/storage"
	"github.com/Woodfyn/file-api/internal/service"
	"github.com/Woodfyn/file-api/internal/transport/rest"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/fbstorage"
	"github.com/Woodfyn/file-api/pkg/hash"
	"github.com/Woodfyn/file-api/pkg/mdb"
	"github.com/Woodfyn/file-api/pkg/rdb"
	"github.com/Woodfyn/file-api/pkg/signaler"
	"github.com/Woodfyn/file-api/pkg/srv"
)

const (
	cfg_folder = "configs"
	cfg_file   = "prod"
)

func init() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(h))
}

func Run() {
	// init config
	cfg, err := config.InitConfig(cfg_folder, cfg_file)
	if err != nil {
		panic(err)
	}

	slog.Info("config loaded", "cfg", cfg)

	// init all clients
	rdbClient := rdb.NewRdbClient(rdb.ConnInfo{
		Addr: cfg.RDB.Addr,
	})

	if err := rdbClient.Ping().Err(); err != nil {
		panic(err)
	}

	defer rdbClient.Close()

	mongoClient, err := mdb.NewMongoClient(context.Background(), mdb.ConnInfo{
		URI:      cfg.Mongo.URI,
		Username: cfg.Mongo.Username,
		Password: cfg.Mongo.Password,
	})
	if err != nil {
		panic(err)
	}
	mongoDB := mongoClient.Database(cfg.Mongo.Database)
	defer mongoClient.Disconnect(context.Background())

	hasher := hash.NewSHA1Hasher(cfg.Password.Salt)
	tokenManager, err := auth.NewManager(cfg.JWT.Secret)
	if err != nil {
		panic(err)
	}

	fbClient, err := fbstorage.NewFBStorageClient(context.Background(), cfg.Firebase.FileName)
	if err != nil {
		panic(err)
	}

	// init dependencies
	redis := redis.NewRepository(rdbClient)
	storage := storage.NewRepository(fbClient.Bucket(cfg.Firebase.BucketName))
	mongo := mongo.NewRepository(mongoDB)

	deps := service.Deps{
		RedisRepos:  redis,
		MongoRepo:   mongo,
		StorageRepo: storage,
		Hasher:      hasher,

		TokenManager: tokenManager,
		AcssTokenTTL: cfg.JWT.AccessTokenTTL,
		RefreshTTL:   cfg.JWT.RefreshTokenTTL,
	}

	service := service.NewService(deps)
	handler := rest.NewHandler(service, tokenManager)

	srv := srv.NewServer(handler.Init())

	// run server
	go func() {
		if err := srv.Run(cfg.Server.Port); err != nil {
			panic(err)
		}
	}()

	slog.Info("server started...")

	// graceful shutdown
	signaler.Wait()

	slog.Info("server stopped...")

	if err := fbClient.Close(); err != nil {
		panic(err)
	}

	// shutdown server
	srv.Shutdown()
}
