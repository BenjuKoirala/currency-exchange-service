package cmd

import (
	"CurrencyExchangeService/currency"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	baseCurrency, targetCurrency string
)

// getRateCmd represents the getrate command
var getRateCmd = &cobra.Command{
	Use:   "getrate",
	Short: "Get exchange rate for a currency pair",
	Long:  `Usage : $exchange-cli getrate -b <src-currency> -t <trg-currency>`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting getrate command")

		log.Printf("Dialing gRPC server at %s", GrpcServerEndpoint)
		conn, err := grpc.Dial(GrpcServerEndpoint, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		client := currency.NewCurrencyServiceClient(conn) // creates a new client for the CurrencyService

		// Setting up context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Call the GetExchangeRate method
		res, err := client.GetExchangeRate(ctx, &currency.ExchangeRateRequest{
			BaseCurrency:   baseCurrency,
			TargetCurrency: targetCurrency,
		})
		if err != nil {
			log.Printf("Error while getting exchage rate: %v", err.Error())
		} else {
			log.Printf("Received exchange rate: %f", res.GetRate())
			fmt.Printf("Exchange rate from %s to %s : %f\n", baseCurrency, targetCurrency, res.GetRate())
		}
	},
}

func init() {
	rootCmd.AddCommand(getRateCmd)

	// Defining flags
	getRateCmd.Flags().StringVarP(&baseCurrency, "base", "b", "", "Base currency")
	getRateCmd.Flags().StringVarP(&targetCurrency, "target", "t", "", "Target currency")

	//Making flags required
	getRateCmd.MarkFlagsRequiredTogether("base", "target")

	log.Println("Initialized getrate command")
}
