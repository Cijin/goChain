package utils

import "math/big"

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

/*
 * Ecoding Algorithm
 *	Till quotient is not zero:
 *		q = input mod base
 *		use remainder to get string of base to be converted to
 *		ex: if remainder 1 and base 2, resulting value is 1
 *		ex: if remainder 2 and base 58, resulting value is 2 (check string above)
 *
 *	Once reulting string found:
 *		Reverse string (mod returns remainder, i.e. last character first)
 *		Account for any leading zeros
 */
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0)
	x.SetBytes(input)

	zero := big.NewInt(0)
	base := big.NewInt(int64((len(b58Alphabet))))
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	// Account for leading zeros
	for b := range input {
		if b != 0x00 {
			break
		}

		result = append([]byte{b58Alphabet[0]}, result...)
	}

	return result
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
