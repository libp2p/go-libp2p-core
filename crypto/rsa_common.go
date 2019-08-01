package crypto

import (
	"fmt"
	"os"
)

// UnsafeRsaKeyEnv is an environment variable which, when set, lowers the
// minimum required bits of RSA keys to 512. This should be used exclusively in
// test situations.
const UnsafeRsaKeyEnv = "LIBP2P_ALLOW_WEAK_RSA_KEYS"

var MinRsaKeyBits = 2048

// ErrRsaKeyTooSmall is returned when trying to generate or parse an RSA key
// that's smaller than MinRsaKeyBits bits. In test
var ErrRsaKeyTooSmall error

func init() {
	if _, ok := os.LookupEnv(UnsafeRsaKeyEnv); ok {
		MinRsaKeyBits = 512
	}

	ErrRsaKeyTooSmall = fmt.Errorf("rsa keys must be >= %d bits to be useful", MinRsaKeyBits)
}
