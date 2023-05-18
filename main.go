package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bndr/gotabulate"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var authorization string
var userVIPLevel uint8

var vipDepositQuotaIn24h = map[uint8]decimal.Decimal{
	1: decimal.NewFromFloat(150000),
	2: decimal.NewFromFloat(500000),
	3: decimal.NewFromFloat(1000000),
	4: decimal.NewFromFloat(2000000),
}

var vipDepositQuotaIn30d = map[uint8]decimal.Decimal{
	1: decimal.NewFromFloat(2000000),
	2: decimal.NewFromFloat(4500000),
	3: decimal.NewFromFloat(6000000),
	4: decimal.NewFromFloat(10000000),
}

var codeMap = map[int]error{
	0:        nil,
	40101631: fmt.Errorf("empty token"),
	40101633: fmt.Errorf("invalid token"),
}

/*
   "code": 0,
    "message": "success",
    "data": {
        "records": [
            {
                "coin": "USDT",
                "complete_time": 1684297260058,
                "create_time": 1684211888224,
                "exchange_rate": "30.855",
                "failure_reason": "ROLLING_QUOTE_AMOUNT_EXCEEDED,AMOUNT_GT_AVG_INCOME",
                "order_id": "d5f4f1fd2f3c8eef7706af2e1783ddd3",
                "reason_code": "DEPOSIT_EXCESS",
                "status": "refunded",
                "transfer_time": 1684211876000,
                "twd_amount": "600000",
                "twd_fee": "0",
                "usdt_amount": "19445.79484686"
            },
*/

type depositRecordResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Records []depositRecord `json:"records"`
	} `json:"data"`
}

type depositRecord struct {
	Coin          string          `json:"coin"`
	CompleteTime  int64           `json:"complete_time"`
	CreateTime    int64           `json:"create_time"`
	ExchangeRate  string          `json:"exchange_rate"`
	FailureReason string          `json:"failure_reason"`
	OrderID       string          `json:"order_id"`
	ReasonCode    string          `json:"reason_code"`
	Status        string          `json:"status"`
	TransferTime  int64           `json:"transfer_time"`
	TwdAmount     decimal.Decimal `json:"twd_amount"`
	TwdFee        decimal.Decimal `json:"twd_fee"`
	UsdtAmount    decimal.Decimal `json:"usdt_amount"`
}

func init() {
	authorization = os.Getenv("AUTH_TOKEN")
	vipLevelStr := os.Getenv("VIP_LEVEL")
	vipUint, err := strconv.ParseUint(vipLevelStr, 10, 8)
	if err != nil {
		panic(err)
	}
	userVIPLevel = uint8(vipUint)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rqt",
		Short: "A command line tool for rybit quota calculation",
		// SilenceUsage is an option to silence usage when an error occurs.
		SilenceUsage: true,
	}

	rootCmd.AddCommand(quotaCmd)

	quotaCmd.Flags().StringVar(&authorization, "auth", authorization, "authorization token (Bearer ...)")
	quotaCmd.Flags().Uint8Var(&userVIPLevel, "vip", userVIPLevel, "VIP level (1, 2, 3, 4)")

	rootCmd.Execute()
}

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "calculate quota in 24 hours and 30 days according your vip level",
	Run: func(cmd *cobra.Command, args []string) {
		quota()
	},
}

func quota() {

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatal(err)
		return
	}

	authorization = strings.Replace(authorization, "\n", "", -1)
	authorization = strings.TrimPrefix(authorization, " ")
	authorization = strings.TrimSuffix(authorization, " ")

	httpReq, err := http.NewRequest(http.MethodGet, "https://www.rybit.com/wallet-api/v1/kgi/deposits", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	httpReq.Header.Set("Authorization", authorization)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	resp := depositRecordResp{}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Fatal(err)
		return
	}

	if err, ok := codeMap[resp.Code]; !ok {
		log.Fatal("unknown code:", resp.Code)
		log.Fatal(resp.Message)
		return
	} else if err != nil {
		log.Fatal(err)
		return
	}

	var depositIn24h []depositRecord
	var depositIn30d []depositRecord

	now := time.Now()
	for _, record := range resp.Data.Records {
		if record.Coin != "USDT" || record.Status != "success" {
			continue
		}
		completedTime := time.UnixMilli(record.CompleteTime)

		switch {
		case now.Sub(completedTime).Seconds() < 60*60*24:
			depositIn24h = append(depositIn24h, record)
			fallthrough
		case now.Sub(completedTime).Seconds() < 60*60*24*30:
			depositIn30d = append(depositIn30d, record)
		}
	}
	sortDepositByTime(depositIn24h)
	sortDepositByTime(depositIn30d)

	var totalDepositIn24h decimal.Decimal
	var totalDepositIn30d decimal.Decimal

	for _, record := range depositIn24h {
		totalDepositIn24h = totalDepositIn24h.Add(record.TwdAmount)
	}

	for _, record := range depositIn30d {
		totalDepositIn30d = totalDepositIn30d.Add(record.TwdAmount)
	}

	remainQuotaIn24h := vipDepositQuotaIn24h[userVIPLevel].Sub(totalDepositIn24h)
	remainQuotaIn30d := vipDepositQuotaIn30d[userVIPLevel].Sub(totalDepositIn30d)

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("*** current remaining deposit quota in 24h:", remainQuotaIn24h.String())

	depositIn24hData := [][]string{}
	for _, record := range depositIn24h {
		depositIn24hData = append(depositIn24hData, []string{
			time.UnixMilli(record.CompleteTime).AddDate(0, 0, 1).In(location).Format("2006-01-02 15:04:05"),
			record.TwdAmount.String(),
		})
	}
	if len(depositIn24hData) != 0 {
		prettyPrint([]string{"unlock time", "amount"}, depositIn24hData)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("*** current remaining deposit quota in 30 days:", remainQuotaIn30d.String())
	depositIn30dData := [][]string{}
	for _, record := range depositIn30d {
		depositIn30dData = append(depositIn30dData, []string{
			time.UnixMilli(record.CompleteTime).AddDate(0, 0, 30).In(location).Format("2006-01-02 15:04:05"),
			record.TwdAmount.String(),
		})
	}

	if len(depositIn30dData) != 0 {
		prettyPrint([]string{"unlock time", "amount"}, depositIn30dData)
	}
}

func prettyPrint(title []string, data [][]string) {
	// Create Object
	tabulate := gotabulate.Create(data)

	// Set Headers
	tabulate.SetHeaders(title)

	// Render
	fmt.Println(tabulate.Render("simple"))
}

func sortDepositByTime(deposits []depositRecord) {
	sort.Slice(deposits, func(i, j int) bool {
		return deposits[i].CompleteTime <= deposits[j].CompleteTime
	})
}
