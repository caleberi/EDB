package pkg

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
	"yc-backend/common"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/samber/lo"
)

var (
	apiKey    string
	apiSecret string
	baseUrl   string
	webhook   string = "https://webhook.site/0c0594d3-6d4b-4008-831b-fbbaa6b77946"
)

var (
	processingEvent string = "payment.PROCESSING"
	pendingEvent    string = "payment.PENDING"
	failedEvent     string = "payment.FAILED"
	completedEvent  string = "payment.COMPLETE"
)

func init() {
	config, err := common.LoadConfiguration(common.ConfEnvSetting{YamlFilePath: []string{"./../dev.yml"}})

	if err != nil {
		log.Panic(err)
	}
	apiKey = config.YellowCardCredentials.ApiKey
	apiSecret = config.YellowCardCredentials.SecretKey
	baseUrl = config.YellowCardCredentials.BaseUrl

	log.Printf("apiKey= %s apiSecret= %s baseUrl= %s", apiKey, apiSecret, baseUrl)
}

func TestYellowCardChannelEndpoint(t *testing.T) {
	yc := NewYellowClient(baseUrl, apiKey, apiSecret)
	path := "/business/channels"
	method := http.MethodGet
	response, err := yc.MakeRequest(method, path, nil)

	if !assert.NoError(t, err, "channel fetching failed") ||
		!assert.Equal(t, response.StatusCode, http.StatusOK) {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if !assert.NoError(t, err, "error occurred while reading body") {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	var channelResponse ChannelResponse
	err = json.Unmarshal(body, &channelResponse)
	assert.NoError(t, err, "unmarshalling channel failed")
}

func TestYellowCardAccountEndpoint(t *testing.T) {
	yc := NewYellowClient(baseUrl, apiKey, apiSecret)
	path := "/business/account"
	method := http.MethodGet
	response, err := yc.MakeRequest(method, path, nil)

	if !assert.NoError(t, err, "account fetching failed") ||
		!assert.Equal(t, response.StatusCode, http.StatusOK) {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if !assert.NoError(t, err, "error occurred while reading body") {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	var accountDetailResponse AccountDetailResponse
	err = json.Unmarshal(body, &accountDetailResponse)
	log.Printf("%v", accountDetailResponse)
	assert.NoError(t, err, "unmarshalling account response failed")
}

func TestYellowCardCreateWebHookEndpoint(t *testing.T) {
	yc := NewYellowClient(baseUrl, apiKey, apiSecret)
	path := "/business/webhooks"
	method := http.MethodPost
	events := []string{processingEvent, failedEvent, completedEvent, pendingEvent}

	lo.ForEach(events, func(event string, idx int) {
		response, err := yc.MakeRequest(method, path, map[string]interface{}{
			"url":    webhook,
			"state":  event,
			"active": true,
		})

		if !assert.NoError(t, err, "webhook creation failed") ||
			!assert.Equal(t, response.StatusCode, http.StatusOK) {
			t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
		}
		response.Body.Close()
	})
}

func TestYellowCardListWebhooksEndpoint(t *testing.T) {
	yc := NewYellowClient(baseUrl, apiKey, apiSecret)
	path := "/business/webhooks"
	method := http.MethodGet
	response, err := yc.MakeRequest(method, path, nil)

	if !assert.NoError(t, err, "webhooks fetching failed") ||
		!assert.Equal(t, response.StatusCode, http.StatusOK) {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if !assert.NoError(t, err, "error occurred while reading body") {
		t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
	}
	var webhooks map[string]interface{}
	err = json.Unmarshal(body, &webhooks)
	log.Printf("%v", webhooks)
	assert.NoError(t, err, "unmarshalling account response failed")
}

func TestYellowCardDeleteWebHookEndpoint(t *testing.T) {
	yc := NewYellowClient(baseUrl, apiKey, apiSecret)
	path := "/business/webhooks"
	method := http.MethodDelete
	ids := []string{"c39147f4-a941-416b-a8a4-cd7e199982b7"}

	lo.ForEach(ids, func(id string, idx int) {
		response, err := yc.MakeRequest(method, path+"/"+id, nil)
		if !assert.NoError(t, err, "webhook deletion failed") ||
			!assert.Equal(t, response.StatusCode, http.StatusNoContent) {
			t.Fatalf("expected %v got: %v", http.StatusOK, response.StatusCode)
		}
		response.Body.Close()
	})
}
