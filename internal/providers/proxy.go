package providers

import "os"

func ProxyEnv() []string {
	http := os.Getenv("MULTIPOWERS_HTTP_PROXY")
	https := os.Getenv("MULTIPOWERS_HTTPS_PROXY")
	if http == "" && https == "" {
		return nil
	}
	out := make([]string, 0, 2)
	if http != "" {
		out = append(out, "HTTP_PROXY="+http)
	}
	if https != "" {
		out = append(out, "HTTPS_PROXY="+https)
	}
	return out
}
