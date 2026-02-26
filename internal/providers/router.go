package providers

func AvailableProviders() []Provider {
	out := make([]Provider, 0)
	for _, p := range Registry() {
		if p.Available() {
			out = append(out, p)
		}
	}
	return out
}
