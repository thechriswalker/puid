package puid

var base36chars = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

// this maps bytes to base36 normalised bytes
// our encoding does not need to be lossless
var base36translator [256]byte
var isBase36byte map[byte]struct{}

func init() {
	for i := 0; i < 256; i++ {
		m := i % BASE
		base36translator[i] = base36chars[m]
	}
	isBase36byte = make(map[byte]struct{}, len(base36chars))
	for _, b := range base36chars {
		isBase36byte[b] = struct{}{}
	}
}

func base36convert(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = base36translator[b[i]]
	}
}

func isAllBase36(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if _, ok := isBase36byte[b[i]]; !ok {
			return false
		}
	}
	return true
}
