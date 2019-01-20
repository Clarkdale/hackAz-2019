package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
)

func main() {

	secret, phrase := hedera.GenerateSecretKey()
	fmt.Printf("secret = %v\n", secret)
	fmt.Printf("phrase = %v\n", phrase)

	public := secret.Public()
	fmt.Printf("public = %v\n", public)

}