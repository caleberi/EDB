package common

import (
	"errors"
	"net/http"
	"strings"
	"yc-backend/internals"
	"yc-backend/repository"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	configContextKey     = "_yc_config"
	requestIdContextKey  = "_yc_request_id"
	dbContextKey         = "__yc_db"
	repositoryContextKey = "__yc_repo"
	poolContextKey       = "__yc_pool"
	loggerContextKey     = "__yc_logger"
	UserKey              = "__user"
)

func GetRequestIdFromCtx(ctx *gin.Context) string {
	return ctx.MustGet(requestIdContextKey).(string)
}

func GetConfigFromCtx(ctx *gin.Context) *Config {
	return ctx.MustGet(configContextKey).(*Config)
}

func GetDbFromCtx(ctx *gin.Context) *mongo.Client {
	return ctx.MustGet(dbContextKey).(*mongo.Client)
}

func GetReposFromCtx(ctx *gin.Context) *repository.Repositories {
	return ctx.MustGet(repositoryContextKey).(*repository.Repositories)
}

func GetLoggerFromCtx(ctx *gin.Context) internals.Logger {
	return ctx.MustGet(loggerContextKey).(internals.Logger)
}

func AddConfigMiddleware(cfg *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(configContextKey, cfg)
		ctx.Next()
	}
}

func AddLoggerMiddleware(logger internals.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(loggerContextKey, logger)
		ctx.Next()
	}
}

func AddReposToMiddleware(db *mongo.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		databaseName := ctx.MustGet(configContextKey).(*Config).MongoDB.DatabaseName
		collection := db.Database(databaseName)
		repos := repository.InitRepositories(collection)
		ctx.Set(repositoryContextKey, repos)
		ctx.Set(dbContextKey, db)
		ctx.Next()
	}
}

func AddRequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New().String()
		ctx.Set(requestIdContextKey, id)
		ctx.Next()
	}
}

func AuthorizeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := GetConfigFromCtx(ctx)
		repo := GetReposFromCtx(ctx)

		value := ctx.GetHeader("Authorization")

		if value == "" || !strings.HasPrefix(value, "Bearer") {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("missing Bearer token"))
			return
		}

		parts := strings.Split(value, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("missing Bearer token"))
			return
		}

		token := parts[1]

		userId, err := utils.Verify(utils.VerificationPayload{
			Token:  token,
			Secret: cfg.JWTCredentials.AccessTokenSecret,
			Issuer: cfg.JWTCredentials.AccessTokenClaim.Issuer,
		})

		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid access token"))
			return
		}
		id, err := primitive.ObjectIDFromHex(userId.(string))
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid access token"))
			return
		}

		user, err := repo.UserRepository.FindOne(ctx, primitive.D{{Key: "_id", Value: id}})
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("user with this email does not exist"))
			return
		}

		ctx.Set(UserKey, user)
		ctx.Next()
	}
}
