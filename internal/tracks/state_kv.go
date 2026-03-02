package tracks

import (
	"encoding/json"
)

// KVGet retrieves a single key value from state
// Returns empty string if key doesn't exist
func KVGet(projectDir, key string) (string, error) {
	s, err := ReadState(projectDir)
	if err != nil {
		return "", err
	}
	if s.Metrics == nil {
		return "", nil
	}
	return s.Metrics[key], nil
}

// KVSet sets a single key-value pair in state
func KVSet(projectDir, key, value string) error {
	s, err := ReadState(projectDir)
	if err != nil {
		return err
	}
	if s.Metrics == nil {
		s.Metrics = make(map[string]string)
	}
	s.Metrics[key] = value
	return WriteState(projectDir, s)
}

// KVUpdate performs atomic merge of multiple key-value pairs
func KVUpdate(projectDir string, updates map[string]string) error {
	s, err := ReadState(projectDir)
	if err != nil {
		return err
	}
	if s.Metrics == nil {
		s.Metrics = make(map[string]string)
	}
	for k, v := range updates {
		s.Metrics[k] = v
	}
	return WriteState(projectDir, s)
}

// KVGetAll returns all key-value pairs from state
func KVGetAll(projectDir string) (map[string]string, error) {
	s, err := ReadState(projectDir)
	if err != nil {
		return nil, err
	}
	if s.Metrics == nil {
		return map[string]string{}, nil
	}
	// Return a copy to prevent mutation
	result := make(map[string]string, len(s.Metrics))
	for k, v := range s.Metrics {
		result[k] = v
	}
	return result, nil
}

// KVDelete removes a key from state
func KVDelete(projectDir, key string) error {
	s, err := ReadState(projectDir)
	if err != nil {
		return err
	}
	if s.Metrics == nil {
		return nil
	}
	delete(s.Metrics, key)
	return WriteState(projectDir, s)
}

// KVGetJSON returns the entire state as JSON bytes
func KVGetJSON(projectDir string) ([]byte, error) {
	s, err := ReadState(projectDir)
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// KVUpdateFromJSON updates state from JSON bytes (atomic merge)
func KVUpdateFromJSON(projectDir string, jsonData []byte) error {
	var updates map[string]string
	if err := json.Unmarshal(jsonData, &updates); err != nil {
		return err
	}
	return KVUpdate(projectDir, updates)
}
