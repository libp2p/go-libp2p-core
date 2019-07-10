package crypto

import (
	"fmt"
)

const MinRsaKeyBits = 2048

// ErrRsaKeyTooSmall is returned when trying to generate or parse an RSA key
// that's smaller than 512 bits. Keys need to be larger enough to sign a 256bit
// hash so this is a reasonable absolute minimum.
var ErrRsaKeyTooSmall = fmt.Errorf("rsa keys must be >= %d bits to be useful", MinRsaKeyBits)
