package extension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cucumber/godog"
)

type (
	// WireMock holds the request and response we will store on the WireMock server
	WireMock struct {
		Request   WireMockRequest  `json:"request"`
		Response  WireMockResponse `json:"response"`
		serverURL string
		client    http.Client
	}

	// WireMockRequest represents the request structure of a mock
	WireMockRequest struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	}

	// WireMockResponse represents the response structure of a mock
	WireMockResponse struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
)

var wm WireMock

// NewWireMock sets the server url on the WireMock struct and resets mocks before scenarios
func NewWireMock(s *godog.Suite, serverURL string) {
	wm.serverURL = serverURL

	s.BeforeScenario(func(interface{}) {
		if err := wm.ResetMocks(); err != nil {
			log.Fatalf("failed to reset mocks before scenarios: %+v", err.Error())
		}
	})
}

func WireMockClient() WireMock {
	return wm
}

// ResetMocks will reset all stored mocks on the WireMock server
func (w *WireMock) ResetMocks() error {
	request, err := http.NewRequest(http.MethodPost, w.serverURL+"/__admin/mappings/reset", nil)
	if err != nil {
		return err
	}

	response, err := w.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf(
				"failed to retrieve the body of the response with status code: %d while resetting mocks",
				response.StatusCode,
			)
		}

		return fmt.Errorf(
			"received unexpected status code while resetting mocks; status code: %d, body: %s",
			response.StatusCode,
			body,
		)
	}

	return nil
}

// SendMocks submits the mocks to the WireMock server
func (w *WireMock) SendMocks() error {
	mocks, err := json.Marshal(w)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, w.serverURL+"/__admin/mappings", bytes.NewReader(mocks))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := w.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf(
				"failed to retrieve the body of the response with status code: %d while resetting mocks",
				response.StatusCode,
			)
		}

		return fmt.Errorf(
			"received unexpected status code while resetting mocks; status code: %d, body: %s",
			response.StatusCode,
			body,
		)
	}

	return nil
}
