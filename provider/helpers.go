package provider

import "math/rand"

func randStringNumber(lenght int) string {
	const charset = "0123456789"
	str := make([]byte, lenght)

	for i := 0; i < lenght; i++ {
		str[i] = charset[rand.Intn(len(charset))]
	}
	return string(str)
}
