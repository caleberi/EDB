package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"yc-backend/common"
	"yc-backend/models"
	"yc-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/square/go-jose/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	hasher = utils.NewHasher(bcrypt.DefaultCost)
)

type CreateUserRequest struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	BVN                string `json:"bvn,omitempty"`
	DOB                string `json:"dob,omitempty"`
	Address            string `json:"address,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Country            string `json:"country,omitempty"`
	IDNumber           string `json:"idNumber,omitempty"`
	IDType             string `json:"idType,omitempty"`
	AdditionalIDType   string `json:"additionalIdType,omitempty"`
	AdditionalIdNumber string `json:"additionalIdNumber,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func RegisterUser(c *gin.Context) {
	repo := common.ReposFromCtx(c)
	logger := common.LoggerFromCtx(c)

	var createUserRequest CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		logger.Infof("bind request to createUserRequest failed : %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx, cancelFunc := context.WithTimeout(c, 5*time.Second)
	defer cancelFunc()

	if _, err := repo.User.FindOne(ctx, bson.D{{Key: "email", Value: createUserRequest.Email}}); !errors.Is(err, mongo.ErrNoDocuments) {
		logger.Infof("an error occurred : %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("user with the provided email exist")))
		return
	}

	hashedPassword, err := hasher.HashPassword(createUserRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	dobTIme, err := time.Parse("2006-01-02", createUserRequest.DOB)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("invalid date of birth")))
		return
	}
	dobTimeStr := dobTIme.Format("04/03/2016")

	user := models.User{
		FirstName:          strings.TrimSpace(createUserRequest.FirstName),
		LastName:           strings.TrimSpace(createUserRequest.LastName),
		Email:              strings.ToLower(createUserRequest.Email),
		Password:           hashedPassword,
		DOB:                dobTimeStr,
		IdType:             createUserRequest.IDType,
		IdNumber:           createUserRequest.IDNumber,
		Phone:              createUserRequest.Phone,
		AdditionalIdType:   createUserRequest.AdditionalIDType,
		AdditionalIdNumber: createUserRequest.AdditionalIdNumber,
		Address:            createUserRequest.Address,
		Country:            createUserRequest.Country,
	}

	id, err := repo.User.Create(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	userId, ok := id.(primitive.ObjectID)
	if !ok {
		err := fmt.Errorf("error occurred while creating user")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	user.ID = userId
	user, err = user.Omit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("user created successfully", user))
}

func LoginUser(c *gin.Context) {
	cfg := common.ConfigFromCtx(c)
	repo := common.ReposFromCtx(c)
	logger := common.LoggerFromCtx(c)

	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		logger.Infof("bind request to createUserRequest failed : %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx, cancelFunc := context.WithTimeout(c, 5*time.Second)
	defer cancelFunc()

	user, err := repo.User.FindOne(ctx, primitive.D{{Key: "email", Value: loginRequest.Email}})
	if err != nil {
		logger.Infof("error during login : %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("user with this email does not exist")))
		return
	}

	err = hasher.CheckPassword(loginRequest.Password, user.Password)
	if err != nil {
		logger.Infof("error during login : %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("password does not match")))
		return
	}

	token, err := utils.Sign(utils.SigningPayload{
		Algorithm: jose.HS256,
		Payload:   user.ID,
		Issuer:    cfg.JWTCredentials.AccessTokenClaim.Issuer,
		Audience:  cfg.JWTCredentials.AccessTokenClaim.Audience,
		Subject:   user.FirstName + ":" + user.LastName,
		Expiry:    cfg.JWTCredentials.AccessTokenTTL,
		Secret:    cfg.JWTCredentials.AccessTokenSecret,
	})

	if err != nil {
		logger.Infof("error during login : %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("error occurred")))
		return
	}

	var loginResponse LoginResponse
	loginResponse.AccessToken = token

	c.JSON(http.StatusOK, utils.SuccessResponse("login successfully", loginResponse))
}

func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("logged out", nil))
}
