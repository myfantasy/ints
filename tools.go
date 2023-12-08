package ints

var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
	"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

var digitsLower = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var digitsByte []byte

const hyphenStr string = "-"

var hyphen byte = []byte(hyphenStr)[0]

var searchDigits map[byte]uint64
var searchDigitsLower map[byte]uint64
var searchDigitsLowerSize int

func init() {
	searchDigits = make(map[byte]uint64)
	searchDigitsLower = make(map[byte]uint64)
	digitsByte = make([]byte, len(digits))
	for k, v := range digits {
		searchDigits[v[0]] = uint64(k)
		digitsByte[k] = v[0]
	}

	for k, v := range digitsLower {
		searchDigitsLower[v[0]] = uint64(k)
	}
	searchDigitsLowerSize = len(searchDigitsLower)

	err := DefaultUuidGenerator.Init()
	if err != nil {
		panic("INIT DefaultUuidGenerator.Init() error: " + err.Error())
	}

	err = DefaultRandomUuidGenerator.Init()
	if err != nil {
		panic("INIT DefaultRandomUuidGenerator.Init() error: " + err.Error())
	}
}

func BytesAppendForward(buf []byte, l int) []byte {
	res := make([]byte, l, l+len(buf))

	res = append(res, buf...)

	return res
}
