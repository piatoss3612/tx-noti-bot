package helpers

import "regexp"

func IsValidEthereumAddress(s string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(s)
}
