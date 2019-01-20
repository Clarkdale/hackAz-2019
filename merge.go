package main

//When we want to set up website import the website

import (
	//"io"
	"fmt"
	//"reflect"
	//"os"
	//"net/http"
	"time"
	"github.com/hashgraph/hedera-sdk-go"
	
)

func main(){
	operatorAccountNumber := 1001
	targetAccountNumber := 1003
	operatorSecretKey := "302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79"
	donationAmount := 1
	client, err := hedera.Dial("testnet.hedera.com:51005")
	if err != nil {
		panic(err)
	}
	//checkBalance(client, 1001, "302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79")
	//createNewAccount(client, 1001, "302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79")
	transferMoney(client, operatorSecretKey, int64(operatorAccountNumber), int64(targetAccountNumber),int64(donationAmount))

	defer client.Close()

	//checkBalance(client, accountID)


	// Get the _answer_ for the query of getting the account balance


/*	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)*/
}

func createNewAccount(client hedera.Client, operatorAccountNumber int64, secretKey string){
	nodeAccountID := hedera.AccountID{Account: 3}
	operatorAccountID := hedera.AccountID{Account: operatorAccountNumber}
	operatorSecret, err := hedera.SecretKeyFromString(secretKey)
	if err != nil {
		panic(err)
	}

	// Generate a new keypair for the new account
	secret, _ := hedera.GenerateSecretKey()
	public := secret.Public()

	fmt.Printf("secret = %v\n", secret)
	fmt.Printf("public = %v\n", public)

	response, err := client.CreateAccount().
		Key(public).
		InitialBalance(0).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Sign(operatorSecret).
		Execute()

	if err != nil {
		panic(err)
	}

	transactionID := response.ID
	fmt.Printf("created account; transaction = %v\n", transactionID)

	//
	// Get receipt to prove we created it ok
	//

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}

	fmt.Printf("account = %v\n", *receipt.AccountID)
		// Target account to get the balance for
}
/*func createReciept(client hedera.Client){
	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}
}*/
/*func setOperator(client hedera.client, operatorBalance int64, operatorSecretKey string){
	operatorAccountID := hedera.AccountID{Account: operatorAccountNumber}
	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(operatorSecretKey)
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})
}*/

func checkBalance(client hedera.Client, operatorAccountNumber int64, operatorSecretKey string){
	operatorAccountID := hedera.AccountID{Account: operatorAccountNumber}
	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(operatorSecretKey)
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})
	balance, err := client.Account(operatorAccountID).Balance().Get()
	if err != nil {
		
	}

	fmt.Printf("balance = %v tinybars\n", balance)
	fmt.Printf("balance = %.5f hbars\n", float64(balance)/100000000.0)

}

func transferMoney(client hedera.Client, operatorSecretKey string,  operatorAccountNumber int64, targetAccountNumber int64, donationAmount int64){
		// Read and decode the operator secret key
	operatorAccountID := hedera.AccountID{Account: operatorAccountNumber}
	operatorSecret, err := hedera.SecretKeyFromString(operatorSecretKey)
	if err != nil {
		panic(err)
	}

	// Read and decode target account
	targetAccountID := hedera.AccountID{Account: targetAccountNumber}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(operatorSecretKey)
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})

	//
	// Get balance for target account
	//

	operatorBalance, err := client.Account(operatorAccountID).Balance().Get()
	if err != nil{
		panic(err)
	}
	fmt.Printf("Operator account balance = %v\n", operatorBalance)

	targetBalance, err := client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}


	fmt.Printf("Target account balance = %v\n", targetBalance)

	//
	// Transfer 100 cryptos to target
	//

	nodeAccountID := hedera.AccountID{Account: 3}
	response, err := client.TransferCrypto().
		// Move 100 out of operator account
		Transfer(operatorAccountID, -donationAmount).
		// And place in our new account
		Transfer(targetAccountID, donationAmount).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Sign(operatorSecret). // Sign it once as operator
		Sign(operatorSecret). // And again as sender
		Execute()

	if err != nil {
		panic(err)
	}

	transactionID := response.ID
	fmt.Printf("transferred; transaction = %v\n", transactionID)

	//
	// Get receipt to prove we sent ok
	//

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	//
	// Get balance for target account (again)
	//
	newOperatorBalance, err := client.Account(operatorAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("new operator account balance = %v\n", newOperatorBalance)

	newTargetBalance, err := client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("new target account balance = %v\n", newTargetBalance)
}




