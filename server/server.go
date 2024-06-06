package main

import (
	"CurrencyExchangeService/currency"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

//----------------------------------------------------------------------------------------------------------------------

// Logrus logger instance
var log = logrus.New()
var apiKey string
var apiURL string

//----------------------------------------------------------------------------------------------------------------------

func init() {
	// Log initialization
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

//----------------------------------------------------------------------------------------------------------------------

func main() {
	// Load configuration during initialization
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	apiKey = config.APIKey
	apiURL = config.ExchangeRateApiUrl

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	currency.RegisterCurrencyServiceServer(s, &server{})
	log.Infof("server is running at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error starting server : %v", err)
	}
}

//----------------------------------------------------------------------------------------------------------------------

type server struct {
	currency.UnimplementedCurrencyServiceServer
}

func (s *server) GetExchangeRate(ctx context.Context, req *currency.ExchangeRateRequest) (*currency.ExchangeRateResponse, error) {
	return s.getExchangeRate(ctx, req, apiURL, apiKey)
}

//----------------------------------------------------------------------------------------------------------------------

func (s *server) getExchangeRate(ctx context.Context, req *currency.ExchangeRateRequest, apiURL, apiKey string) (*currency.ExchangeRateResponse, error) {
	log.Info("Incoming request to the exchange server")
	client := resty.New()
	url := apiURL + apiKey + "/pair/" + req.BaseCurrency + "/" + req.TargetCurrency
	log.Infof("Sending request to url %s", url)

	resp, err := client.R().Get(url)
	if err != nil {
		log.Errorf("error executing request : %s", err.Error())
		return nil, fmt.Errorf("error executing request : unable to connect to exchange rate api")
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	status := result["result"]
	log.Infof("response from exchange server : %s", status)
	if status == "success" {
		rate := result["conversion_rate"].(float64)
		return &currency.ExchangeRateResponse{Rate: rate}, nil
	}
	errorType := result["error-type"]
	log.Errorf("error performing exchange : %s", errorType)
	return nil, fmt.Errorf("error performing exchange , reason : %s", errorType)
}

//----------------------------------------------------------------------------------------------------------------------

/**CONFIG**/

// Config struct to hold configuration parameters
type Config struct {
	ExchangeRateApiUrl string `json:"exchangeRateApiUrl"`
	APIKey             string `json:"apiKey"`
}

// LoadConfig reads the configuration from a JSON file
func loadConfig(filePath string) (*Config, error) {
	log.Infof("Loading config file...")
	if filePath != "" {
		configFile, err := os.Open(filePath)
		if err != nil {
			log.Errorf("Error while opening config file %s", err.Error())
			return nil, err
		}
		defer func(configFile *os.File) {
			err := configFile.Close()
			if err != nil {
				log.Errorf("Error while closing config file %s", err.Error())
			}
		}(configFile)

		config := &Config{}
		jsonParser := json.NewDecoder(configFile)
		err = jsonParser.Decode(config)
		if err != nil {
			log.Errorf("Error while marshalling config file %s", err.Error())
			return nil, err
		}
		return config, nil
	}

	// This is for checking in a test environment
	if os.Getenv("TEST_ENV") == "true" {
		return &Config{
			ExchangeRateApiUrl: "http://mockedAPIURL/",
			APIKey:             "mockedAPIKey",
		}, nil
	}

	// If neither file path nor test environment is set, return an error
	return nil, errors.New("no configuration provided")
}

//----------------------------------------------------------------------------------------------------------------------
