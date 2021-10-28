package gencrypt

type Alg string

const (
	/*
		AES key length must be one of 128, 192, 256 bits.
	*/
	ALG_AES_128_CBC Alg = "ALG_AES_128_CBC"
	ALG_AES_128_CTR Alg = "ALG_AES_128_CTR"
	ALG_AES_256_CBC Alg = "ALG_AES_256_CBC"
	ALG_AES_256_CTR Alg = "ALG_AES_256_CTR"
)
