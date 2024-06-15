package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"yc-backend/common"
	"yc-backend/models"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateEmployeeRequest struct {
	FirstName        string  `json:"firstName,omitempty" validate:"required"`
	LastName         string  `json:"lastName,omitempty" validate:"required"`
	MiddleName       string  `json:"middleName,omitempty"`
	Email            string  `json:"email,omitempty" validate:"required,email"`
	BVN              string  `json:"bvn,omitempty"`
	DOB              string  `json:"dob,omitempty"`
	Address          string  `json:"address,omitempty"`
	Phone            string  `json:"phone,omitempty"`
	Country          string  `json:"country,omitempty"`
	IDNumber         string  `json:"idNumber,omitempty"`
	IDType           string  `json:"idType,omitempty"`
	AdditionalIDType string  `json:"additionalIdType,omitempty"`
	Salary           float64 `json:"salary,omitempty" validate:"required"`
	AccountName      string  `json:"account_name,omitempty" validate:"required"`
	BankName         string  `json:"bank_name,omitempty" validate:"required"`
	AccountType      string  `json:"account_type,omitempty" validate:"required"`
}

type UpdateEmployeeRequest struct {
	FirstName        string  `json:"firstName,omitempty" validate:"required"`
	LastName         string  `json:"lastName,omitempty" validate:"required"`
	MiddleName       string  `json:"middleName,omitempty"`
	Address          string  `json:"address,omitempty"`
	Phone            string  `json:"phone,omitempty"`
	Country          string  `json:"country,omitempty"`
	IDNumber         string  `json:"idNumber,omitempty"`
	IDType           string  `json:"idType,omitempty"`
	AdditionalIDType string  `json:"additionalIdType,omitempty"`
	Salary           float64 `json:"salary,omitempty" validate:"required"`
	AccountName      string  `json:"account_name,omitempty" validate:"required"`
	BankName         string  `json:"bank_name,omitempty" validate:"required"`
	AccountType      string  `json:"account_type,omitempty" validate:"required"`
	Bvn              string  `json:"bvn,omitempty" validate:"required"`
}

func AddEmployee(ctx *gin.Context) {
	logger := common.GetLoggerFromCtx(ctx)
	repo := common.GetReposFromCtx(ctx)

	user, ok := ctx.MustGet(common.UserKey).(*models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("internal server error")))
		return
	}

	var employeeRequest CreateEmployeeRequest
	if err := ctx.ShouldBindJSON(&employeeRequest); err != nil {
		logger.Errorf("bind request to CreateEmployeeRequest failed: %v", err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	logger.Infof("Received employee request: %+v", employeeRequest)

	employee := models.Employee{
		Email:            employeeRequest.Email,
		FirstName:        employeeRequest.FirstName,
		LastName:         employeeRequest.LastName,
		BVN:              employeeRequest.BVN,
		UpdatedAt:        time.Now(),
		CreatedAt:        time.Now(),
		DOB:              employeeRequest.DOB,
		IDType:           employeeRequest.IDType,
		IDNumber:         employeeRequest.IDNumber,
		Salary:           employeeRequest.Salary,
		Phone:            employeeRequest.Phone,
		AdditionalIDType: employeeRequest.AdditionalIDType,
		Address:          employeeRequest.Address,
		BankName:         employeeRequest.BankName,
		Country:          employeeRequest.Country,
		AccountName:      employeeRequest.AccountName,
		AccountType:      employeeRequest.AccountType,
		UserID:           user.ID,
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Infof("Checking if employee with email %s already exists", employeeRequest.Email)
	_, err := repo.EmployeeRepository.FindOne(ctxWithTimeout, bson.D{{Key: "email", Value: employeeRequest.Email}})
	if !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Errorf("Employee with the provided email exists already: %v", err)
		ctx.JSON(http.StatusConflict, utils.ErrorResponse(errors.New("employee with the provided email exists already")))
		return
	}

	id, err := repo.EmployeeRepository.Create(ctxWithTimeout, employee)
	if err != nil {
		logger.Errorf("Error occurred while creating employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	employeeId, ok := id.(primitive.ObjectID)
	if !ok {
		logger.Errorf("Invalid type assertion for employee ID: %v", id)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("internal server error")))
		return
	}

	employee.ID = employeeId
	ctx.JSON(http.StatusOK, utils.SuccessResponse("employee created successfully", employee))
}

func DeleteEmployee(ctx *gin.Context) {
	repo := common.GetReposFromCtx(ctx)

	user := ctx.MustGet(common.UserKey).(*models.User)

	employeeId, err := primitive.ObjectIDFromHex(ctx.Param("employeeId"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	query := primitive.D{
		{Key: "user_id", Value: user.ID},
		{Key: "_id", Value: employeeId},
	}

	// Delete many is not proper here:, but it make ease witht the model definition
	if err := repo.EmployeeRepository.DeleteMany(ctx, query); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("could not delete employee with id [%v]", employeeId.String())))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("", nil))
}

func UpdateEmployee(ctx *gin.Context) {
	repo := common.GetReposFromCtx(ctx)
	logger := common.GetLoggerFromCtx(ctx)

	employeeId, err := primitive.ObjectIDFromHex(ctx.Param("employeeId"))

	var employeeeRequest UpdateEmployeeRequest

	if err := ctx.ShouldBindJSON(&employeeeRequest); err != nil {
		logger.Infof("bind request to createUserRequest failed : %v", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	employee := models.Employee{
		FirstName:        employeeeRequest.FirstName,
		LastName:         employeeeRequest.LastName,
		UpdatedAt:        time.Now(),
		IDType:           employeeeRequest.IDType,
		IDNumber:         employeeeRequest.IDNumber,
		Salary:           employeeeRequest.Salary,
		Phone:            employeeeRequest.Phone,
		AdditionalIDType: employeeeRequest.AdditionalIDType,
		Address:          employeeeRequest.Address,
		AccountName:      employeeeRequest.AccountName,
		AccountType:      employeeeRequest.AccountType,
		BankName:         employeeeRequest.BankName,
		BVN:              employeeeRequest.Bvn,
	}

	if err := repo.EmployeeRepository.UpdateOneById(ctx, employeeId, employee); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("could not update employee with id [%v]", employeeId.String())))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse("", nil))
}
