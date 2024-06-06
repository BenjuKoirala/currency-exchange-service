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

var baseCurrency, targetCurrency string

// getRateCmd represents the getrate command
var getRateCmd = &cobra.Command{
	Use:   "getrate",
	Short: "Get exchange rate for a currency pair",
	Long:  `Usage : $exchange-rate <src-currency> <trg-currency>`,
	Run: func(cmd *cobra.Command, args []string) {
		// Dial the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		client := currency.NewCurrencyServiceClient(conn)

		// Set up context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Call the GetExchangeRate method
		res, err := client.GetExchangeRate(ctx, &currency.ExchangeRateRequest{
			BaseCurrency:   baseCurrency,
			TargetCurrency: targetCurrency,
		})
		if err != nil {
			log.Printf(err.Error())
		} else {
			fmt.Printf("Exchange rate from %s to %s : %f\n", baseCurrency, targetCurrency, res.GetRate())
		}
	},
}

func init() {
	rootCmd.AddCommand(getRateCmd)

	// Define flags and configuration settings
	getRateCmd.Flags().StringVarP(&baseCurrency, "base", "b", "USD", "Base currency")
	getRateCmd.Flags().StringVarP(&targetCurrency, "target", "t", "EUR", "Target currency")
}
