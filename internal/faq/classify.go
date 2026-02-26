package faq

func Classify(msg string) string {
	switch {
	case msg == "":
		return "unknown"
	default:
		return "runtime-prerun"
	}
}
