package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yc-backend/common"
	"yc-backend/models"
	"yc-backend/pkg"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MakeDisbursmentToEmployee(ctx *gin.Context) {
	repo := common.ReposFromCtx(ctx)
	logger := common.LoggerFromCtx(ctx)
	cfg := common.ConfigFromCtx(ctx)

	user, ok := ctx.MustGet(common.UserKey).(*models.User)
	if !ok {
		err := fmt.Errorf("error occurred while creating user")
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	employeeId, err := primitive.ObjectIDFromHex(ctx.Param("employeeId"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	query := primitive.D{{Key: "_id", Value: employeeId}}

	employee, err := repo.Employee.FindOne(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	client := pkg.NewYellowClient(
		cfg.YellowCardCredentials.BaseUrl,
		cfg.YellowCardCredentials.ApiKey,
		cfg.YellowCardCredentials.SecretKey)

	paymentDetails := map[string]interface{}{
		"channelId":   "fe8f4989-3bf6-41ca-9621-ffe2bc127569",
		"sequenceId":  uuid.New().String(),
		"localAmount": employee.Salary,
		"reason":      "other",
		"sender": map[string]interface{}{
			"name":               user.FirstName + " " + user.LastName,
			"phone":              user.Phone,
			"country":            user.Country,
			"address":            user.Address,
			"dob":                user.DOB,
			"email":              user.Email,
			"idNumber":           user.IdNumber,
			"idType":             user.IdType,
			"businessId":         "B1234567",
			"businessName":       "Example Inc.",
			"additionalIdType":   user.AdditionalIdType,
			"additionalIdNumber": user.AdditionalIdNumber,
		},
		"destination": map[string]interface{}{
			"accountNumber": employee.AccountName,
			"accountType":   "bank",
			"networkId":     "31cfcc77-8904-4f86-879c-a0d18b4b9365",
			"accountBank":   employee.BankName,
			"networkName":   "Guaranty Trust Bank",
			"country":       employee.Country,
			"accountName":   employee.FirstName + " " + employee.LastName,
			"phoneNumber":   employee.Phone,
		},
		"forceAccept":  true,
		"customerType": "retail",
	}

	var payment models.Payment
	resp, err := client.MakeRequest(http.MethodPost, "/business/payments", paymentDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	err = json.Unmarshal(body, &payment)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	logger.Infof("Payment = %+v", payment)
	timeNow := time.Now()
	disbursment := models.Disbursement{
		ReceiverID:   employee.ID,
		SenderID:     user.ID,
		CreatedAt:    &timeNow,
		SalaryAmount: employee.Salary,
		Status:       "processing",
		Payment:      payment,
	}

	_, err = repo.Disbursement.Create(ctx, disbursment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("disbursement submitted successfully", disbursment))
}
