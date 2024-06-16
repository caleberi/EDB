package controllers

import (
	"net/http"
	"yc-backend/common"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func PingDb(ctx *gin.Context) {
	db := common.DbFromCtx(ctx)
	logger := common.LoggerFromCtx(ctx)

	err := db.Ping(ctx, &readpref.ReadPref{})
	if err != nil {
		logger.Infof("error pinging db : %#v", err)
		ctx.JSON(http.StatusInternalServerError, utils.SuccessResponse("Error occurred", nil))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse("Pinged successfully", nil))
}
