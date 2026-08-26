package generator

import "crypto/rand"

const pkgTag = "generator"
const symbols = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHILKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return b
}
