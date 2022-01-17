package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	orangeSDK "github.com/orange-protocol/orange-provider-go-sdk"
	orangeOnt "github.com/orange-protocol/orange-provider-go-sdk/ont"
)

type BalanceReq struct {
	UserDID string `json:"user_did" `
	Address string `json:"address"`
	Chain string `json:"chain" `
	Encrypt bool `json:"encrypt"`
}



type BalanceData struct {
	Balance string `json:"balance"`
}

type RespData struct {
	Data BalanceData `json:"data"`
	Sig string `json:"sig"`
}

var didsdk *orangeSDK.OrangeProviderSdk
var selfDID = "did:ont:ASwHNVY8jvtuJoxbFKDcz1KkVCxcYUvSj2"

func main() {
	didsdk,err:= orangeOnt.NewOrangeProviderOntSdk("./wallet.dat","123456","TESTNET")
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/balance",func(c *gin.Context){
		fmt.Println("=========================")
		requestJson:=&BalanceReq{}
		if err := c.ShouldBindJSON(requestJson);err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		fmt.Printf("address:%s\n",requestJson.Address)
		fmt.Printf("chain:%s\n",requestJson.Chain)
		fmt.Printf("encrypt:%v\n",requestJson.Encrypt)

		balanceData := BalanceData{Balance:"1000000"}
		dataToSign ,err:= json.Marshal(balanceData)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		sig, err := didsdk.SignData(dataToSign)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}

		dataWithSig := RespData{
			Data: balanceData,
			Sig:  hex.EncodeToString(sig),
		}

		if requestJson.Encrypt {
			databytes, err := json.Marshal(dataWithSig)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
				return
			}

			enctrypted, err := didsdk.EncryptDataWithDID(databytes, requestJson.UserDID)
			if err != nil {
				c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
				return
			}

			enhex := hex.EncodeToString(enctrypted)
			c.JSON(200, gin.H{
				"provider_did":selfDID,
				"data": nil,
				"encrypted":enhex,
			})
		}else{
			c.JSON(200, gin.H{
				"provider_did":selfDID,
				"data": dataWithSig,
				"encrypted":nil,
			})
		}

	})
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}