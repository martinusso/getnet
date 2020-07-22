package getnet

func maxLength(s string, l int) string {
	if l > len(s) {
		l = len(s)
	}
	return s[0:l]
}
