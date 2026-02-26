package faq

type Event struct {
	Type      string `json:"type"`
	RootCause string `json:"root_cause"`
	Fix       string `json:"fix"`
}
