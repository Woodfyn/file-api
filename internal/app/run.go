package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/Woodfyn/file-api/internal/config"
	"github.com/Woodfyn/file-api/internal/repository/mongo"
	"github.com/Woodfyn/file-api/internal/repository/storage"
	"github.com/Woodfyn/file-api/internal/service"
	"github.com/Woodfyn/file-api/internal/transport"
	"github.com/Woodfyn/file-api/internal/transport/rest"
	"github.com/Woodfyn/file-api/pkg/auth"
	"github.com/Woodfyn/file-api/pkg/hash"
	"github.com/Woodfyn/file-api/pkg/mdb"
	"github.com/Woodfyn/file-api/pkg/signaler"
	"github.com/Woodfyn/file-api/pkg/srv"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	CFG_FOLDER = "configs"
	CFG_FILE   = "prod"
)

var appCtx = context.Background()

func init() {
	log := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(log))
}

func Run() {
	cfg, err := config.InitConfig(CFG_FOLDER, CFG_FILE)
	if err != nil {
		panic(err)
	}

	db, err := mdb.NewMongoClient(appCtx, mdb.ConnInfo{
		URI:      cfg.Mongo.URI,
		Username: cfg.Mongo.Username,
		Password: cfg.Mongo.Password,
		Database: cfg.Mongo.Database,
	})
	if err != nil {
		panic(err)
	}

	awsCfg, err := awsCfg.LoadDefaultConfig(appCtx, awsCfg.WithRegion(cfg.AWS.Region))
	if err != nil {
		panic(err)
	}

	storageS3 := s3.NewFromConfig(awsCfg)

	presignS3 := s3.NewPresignClient(storageS3)

	hasher := hash.NewSHA256Hasher(cfg.Password.Salt)

	authManager, err := auth.NewManager(cfg.JWT.Secret)
	if err != nil {
		panic(err)
	}

	handler := transport.NewHandler(rest.NewHandler(service.NewFile(mongo.NewFile(db), storage.NewFile(storageS3, presignS3, cfg.AWS.BucketName)),
		service.NewAuth(mongo.NewAuth(db), hasher, authManager, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)))

	slog.Info("Starting server...")
	server := srv.NewServer(handler.Init())
	server.Run(cfg.Server.Port)

	signaler.Wait()

	slog.Info("Shutting down...")
	server.Shutdown()
}
