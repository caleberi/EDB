package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"yc-backend/common"
	"yc-backend/models"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ProcessingEvent string = "PAYMENT.PROCESSING"
	PendingEvent    string = "PAYMENT.PENDING"
	FailedEvent     string = "PAYMENT.FAILED"
	CompletedEvent  string = "PAYMENT.COMPLETE"
)

type Webhook struct {
	ID         string `json:"id"`
	SequenceID string `json:"sequenceId"`
	Status     string `json:"status"`
	ApiKey     string `json:"apiKey"`
	Event      string `json:"event"`
	ExecutedAt int64  `json:"executedAt"`
}

func YellowCardWebHook(ctx *gin.Context) {
	repo := common.ReposFromCtx(ctx)
	logger := common.LoggerFromCtx(ctx)

	if validateSignature(ctx) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("validating request to webhook payload failed")))
		return
	}

	var hook Webhook
	if err := ctx.ShouldBindJSON(&hook); err != nil {
		logger.Errorf("bind request to webhook failed: %v", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	disbursement, err := repo.DisbursementRepository.FindOne(ctx, primitive.D{{Key: "payment.sequenceid", Value: hook.SequenceID}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	switch hook.Event {
	case PendingEvent, PendingEvent, CompletedEvent, FailedEvent:
		err = repo.DisbursementRepository.UpdateOneById(ctx, disbursement.ID, models.Disbursement{Status: hook.Status})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	default:
	}

}

func validateSignature(ctx *gin.Context) bool {
	cfg := common.ConfigFromCtx(ctx)
	receivedSignature := ctx.GetHeader("X-YC-Signature")
	if receivedSignature == "" {
		return false
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return false
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	h := hmac.New(sha256.New, []byte(cfg.YellowCardCredentials.SecretKey))
	h.Write(body)
	computedHash := h.Sum(nil)
	computedSignature := base64.StdEncoding.EncodeToString(computedHash)
	return hmac.Equal([]byte(receivedSignature), []byte(computedSignature))
}
