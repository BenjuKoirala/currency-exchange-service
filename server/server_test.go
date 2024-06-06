package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"CurrencyExchangeService/currency"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func startTestServer(apiURL string) {
	s := grpc.NewServer()
	srv := &server{}
	currency.RegisterCurrencyServiceServer(s, srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func TestMain(m *testing.M) {
	os.Setenv("TEST_ENV", "true")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetExchangeRate_Success(t *testing.T) {
	// Mock the external API
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mockedAPIKey/pair/USD/CAD", r.URL.Path)
		response := map[string]interface{}{
			"result":          "success",
			"conversion_rate": 1.25,
		}
		json.NewEncoder(w).Encode(response)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Override the apiURL with the mock server URL
	apiURL = server.URL + "/"
	apiKey = "mockedAPIKey"

	// Start the gRPC test server with the mock API URL
	startTestServer(server.URL + "/")

	// Set up the gRPC client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := currency.NewCurrencyServiceClient(conn)

	// Call the GetExchangeRate method
	req := &currency.ExchangeRateRequest{BaseCurrency: "USD", TargetCurrency: "CAD"}
	res, err := client.GetExchangeRate(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, 1.25, res.GetRate())
}

func TestGetExchangeRate_Failure(t *testing.T) {
	// Mock the external API
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mockedAPIKey/pair/USD/CD", r.URL.Path)
		response := map[string]interface{}{
			"result":     "error",
			"error-type": "unsupported-code",
		}
		json.NewEncoder(w).Encode(response)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Override the apiURL with the mock server URL
	apiURL = server.URL + "/"
	apiKey = "mockedAPIKey"

	// Start the gRPC test server with the mock API URL
	startTestServer(server.URL + "/v6/")

	// Set up the gRPC client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := currency.NewCurrencyServiceClient(conn)

	// Call the GetExchangeRate method
	req := &currency.ExchangeRateRequest{BaseCurrency: "USD", TargetCurrency: "CD"}
	_, err = client.GetExchangeRate(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rpc error: code = Unknown desc = error performing exchange , reason : unsupported-code")
}

func TestLoadConfig_NormalScenario(t *testing.T) {
	// Explicitly set TEST_ENV to false
	os.Setenv("TEST_ENV", "false")
	defer os.Unsetenv("TEST_ENV")

	// Create a temporary configuration file for testing
	filePath := "test_config.json"
	defer os.Remove(filePath)
	fileContent := `{"exchangeRateApiUrl": "http://example.com/", "apiKey": "testAPIKey"}`
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	assert.NoError(t, err)

	// Load the configuration
	config, err := loadConfig(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/", config.ExchangeRateApiUrl)
	assert.Equal(t, "testAPIKey", config.APIKey)
}

func TestLoadConfig_TestScenario(t *testing.T) {
	// Set the TEST_ENV variable to true
	os.Setenv("TEST_ENV", "true")
	defer os.Unsetenv("TEST_ENV")

	// Load the configuration
	config, err := loadConfig("")
	assert.NoError(t, err)
	assert.Equal(t, "http://mockedAPIURL/", config.ExchangeRateApiUrl)
	assert.Equal(t, "mockedAPIKey", config.APIKey)
}
