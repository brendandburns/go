package extensions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

// Result is a watch result
type Result struct {
	Type string `json:"type"`
	Object map[string]interface{} `json:"object"`
}

// WatchClient is a client for Watching the Kubernetes API
type WatchClient struct {
	Client *http.Client
	URL string
}

func (w *WatchClient) Connect(output chan<- Result) error {
	// TODO: this is hacky and brittle
	res, err := w.Client.Get(w.URL + "?watch=1")
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Error connecting watch (%d: %s)", res.StatusCode, res.Status)
	}
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		watchObj := Result{}
		if err := json.Unmarshal([]byte(line), &watchObj); err != nil {
			return err
		} else {
			output <- watchObj
		}
	}

	return scanner.Err()
}