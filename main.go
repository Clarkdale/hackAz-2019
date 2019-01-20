package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	//"honnef.co/go/js/dom"
)

var masterId int = -1;

func checkBalance(client hedera.Client, operatorAccountNumber int64, operatorSecretKey string, recieverAccountNumber int64, w http.ResponseWriter, r *http.Request){
	operatorAccountID := hedera.AccountID{Account: operatorAccountNumber}
	targetAccount := hedera.AccountID{Account: recieverAccountNumber}
	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(operatorSecretKey)
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})
	balance, err := client.Account(targetAccount).Balance().Get()
	if err != nil {
		fmt.Printf("noooooooooo\n")
	}
	//stringToWrite := "balance = " + strconv.FormatUint(balance, 10)
	//d := dom.GetWindow().Document()
	//e1 := d.GetElementByID("accountBalance")
	//e1.SetInnerHTML(stringToWrite)

	fmt.Printf("balance = %v tinybars\n", balance)
	fmt.Printf("balance = %.5f hbars\n", float64(balance)/100000000.0)

	fmt.Fprintf(w, "balance = %v tinybars\n", balance)
	fmt.Fprintf(w, "balance = %.5f hbars\n", float64(balance)/100000000.0)

}

func check(w http.ResponseWriter, r *http.Request) {
	var (
		id int
		email string 
		password string 
		isCharity int
		accountNumber string
		operatorAccountNumber int64
		operatorSecretKey string
	)

	db, err := sql.Open("mysql", "root:@/giveCoin"); if err != nil {
		panic(err)
	}

	acct := r.FormValue("acct")

	rows, err := db.Query("SELECT * FROM users WHERE ID=" + acct); if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &isCharity, &accountNumber, &operatorSecretKey); if err != nil {
			panic(err)
		}
		operatorAccountNumber, err = strconv.ParseInt(accountNumber, 10, 64); if err != nil {
			panic(err)
		}
	}	
	
	client, err := hedera.Dial("testnet.hedera.com:51005")
	if err != nil {
		panic(err)
	}

	checkBalance(client, 1001, "302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79", operatorAccountNumber, w, r)
}

func transferMoney(client hedera.Client, operatorSecretKey string,  operatorAccountNumber int64, targetAccountNumber int64, donationAmount int64, w http.ResponseWriter, r *http.Request){
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

	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		//http.ServeFile(w, r, "failure.html")
		fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String())
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

	http.ServeFile(w, r, "success.html")

	fmt.Printf("new target account balance = %v\n", newTargetBalance)
}

func transfers(w http.ResponseWriter, r *http.Request) {
	var (
		id int
		email string 
		password string 
		isCharity int
		accountNumber string
		operatorAccountNumber int64
		operatorSecretKey string
		targetAccountNumber int64
		secretKey string 
		donationAmount int64
	)

	db, err := sql.Open("mysql", "root:@/giveCoin"); if err != nil {
		panic(err)
	}

	from := r.FormValue("from")
	to := r.FormValue("to")
	amount := r.FormValue("amt")
	donationAmount, err = strconv.ParseInt(amount, 10, 64); if err != nil {
		fmt.Printf("one\n")
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE ID=" + from); if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &isCharity, &accountNumber, &operatorSecretKey); if err != nil {
			panic(err)
		}
		operatorAccountNumber, err = strconv.ParseInt(accountNumber, 10, 64); if err != nil {
			fmt.Printf("two\n")
			panic(err)
		}
	}

	defer db.Close()

	rows, err2 := db.Query("SELECT * FROM users WHERE ID=" + to); if err != nil {
		panic(err2)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &isCharity, &accountNumber, &secretKey); if err != nil {
			panic(err)
		}
		targetAccountNumber, err = strconv.ParseInt(accountNumber, 10, 64); if err != nil {
			fmt.Printf("3\n")
			panic(err)
		}
	}

	client, err := hedera.Dial("testnet.hedera.com:51005")
	if err != nil {
		panic(err)
	}
	
	transferMoney(client, operatorSecretKey, int64(operatorAccountNumber), int64(targetAccountNumber),int64(donationAmount), w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	var (
		id int
		email string 
		password string 
		isCharity int
		accountNumber string
		operatorAccountNumber int64
		operatorSecretKey string

	)
	username := r.FormValue("user")
	passwordq := r.FormValue("pass")

	db, err := sql.Open("mysql", "root:@/giveCoin"); if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM users WHERE email='" + username + "' AND password='" + passwordq + "'"); if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &password, &isCharity, &accountNumber, &operatorSecretKey); if err != nil {
			panic(err)
		}
	}

	client, err := hedera.Dial("testnet.hedera.com:51005")
	if err != nil {
		panic(err)
	}

	operatorAccountNumber, err = strconv.ParseInt(accountNumber, 10, 64); if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", id)
	if (email != username) {
		http.ServeFile(w, r, "failure.html")
	} else {
		masterId = id;
		fmt.Printf("%v\n", masterId)
		http.ServeFile(w, r, "profile.html")
		checkBalance(client, 1001, "302e020100300506032b657004220420aaa58cd91d6d5bbaac4e713f96021712804467c105438f8ed970f950e8cd1c79", operatorAccountNumber, w, r)
	}
}

func main() {
	// Target account to get the balance for
	// ex /transfers?from=id1&to=id2
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.HandleFunc("/transfers", transfers)
	http.HandleFunc("/login", login)
	// ex /check?acct=id
	http.HandleFunc("/check", check)
	if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
	}
}