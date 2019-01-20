package main

//When we want to set up website import the website

import (
	//"io"
	"fmt"
	//"os"
	//"net/http"
	"github.com/hashgraph/hedera-sdk-go"
	
)

func main(){
		// Target account to get the balance for
	accountID := hedera.AccountID{Account: 1001}

	client, err := hedera.Dial("testnet.hedera.com:51005")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(accountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString("302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79")
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})

	defer client.Close()

	// Get the _answer_ for the query of getting the account balance
	balance, err := client.Account(accountID).Balance().Get()
	if err != nil {
		
	}

	fmt.Printf("balance = %v tinybars\n", balance)
	fmt.Printf("balance = %.5f hbars\n", float64(balance)/100000000.0)


/*	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)*/
}



