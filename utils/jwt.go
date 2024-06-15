package utils

import (
	"time"
	"yc-backend/internals"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

var logger internals.Logger

func init() {
	logger = internals.GetLogger()
}

type SigningPayload struct {
	Payload                   any
	Secret                    string
	Alogrithm                 jose.SignatureAlgorithm
	Issuer, Subject, Audience string
	Expiry                    time.Duration
}

type VerificationPayload struct {
	Token            string
	Secret           string
	Issuer, Audience string
}

func Sign(payload SigningPayload) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: payload.Alogrithm, Key: []byte(payload.Secret)}, nil)
	if err != nil {
		logger.Errorf("Failed to create signer: %v\n", err)
		return "", err
	}

	claims := jwt.Claims{
		Issuer:   payload.Issuer,
		Subject:  payload.Subject,
		Audience: jwt.Audience{payload.Audience},
		Expiry:   jwt.NewNumericDate(time.Now().Add(payload.Expiry)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	data := struct {
		jwt.Claims
		Payload any `json:"payload"`
	}{
		Claims:  claims,
		Payload: payload.Payload,
	}

	rawJwt, err := jwt.Signed(signer).Claims(data).CompactSerialize()
	if err != nil {
		logger.Errorf("Failed to create JWT: %v\n", err)
		return "", err
	}

	return rawJwt, nil
}

func Verify(payload VerificationPayload) (any, error) {
	parsedToken, err := jwt.ParseSigned(payload.Token)

	if err != nil {
		logger.Errorf("Failed to parse JWT: %v\n", err)
		return nil, err
	}

	claims := struct {
		jwt.Claims
		Payload any `json:"payload"`
	}{}

	err = parsedToken.Claims([]byte(payload.Secret), &claims)
	if err != nil {
		logger.Errorf("Failed to verify JWT: %v\n", err)
		return nil, err
	}

	err = claims.Validate(jwt.Expected{
		Issuer:   payload.Issuer,
		Audience: jwt.Audience{payload.Audience},
		Time:     time.Now(),
	})
	if err != nil {
		logger.Errorf("JWT claims validation failed: %v\n", err)
		return nil, err
	}

	return claims.Payload, nil
}
