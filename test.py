import hedera

secret, mnemonic = hedera.SecretKey.generate("")
public = secret.public

print(secret)