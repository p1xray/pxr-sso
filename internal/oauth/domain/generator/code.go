package generator

import "math/rand"

const symbols = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHILKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
const authorizationCodeLength = 10

func AuthorizationCode() string {
	code := randomString(authorizationCodeLength)
	return code
}

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
