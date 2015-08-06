package particle

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	Name string
	Data struct {
		Data      string `json:"data"`
		TTL       string `json:"ttl"`
		Timestamp string `json:"published_at"`
		CoreID    string `json:"coreid"`
	}
}

const URL string = "https://api.particle.io/v1/devices/events/"

// subscribes to a particle.io event stream and returns a channel to receive them
func Subscribe(eventPrefix string, token string) <-chan Event {
	out := make(chan Event)

	var client *http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{tr, nil, nil, 0 * time.Second}

	req, err := http.NewRequest("GET", URL+eventPrefix+"?access_token="+token, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)

	// check for :ok as first event on stream
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	} else if line != ":ok\n" {
		log.Fatal(line)
	}

	go func() {
		for {
			var event Event

			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			} else if strings.HasPrefix(line, "event:") {
				event.Name = strings.TrimPrefix(strings.TrimSuffix(line, "\n"), "event: ")
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				} else if strings.HasPrefix(line, "data:") {
					jsonBlob := strings.TrimPrefix(strings.TrimSuffix(line, "\n"), "data: ")
					err := json.Unmarshal([]byte(jsonBlob), &event.Data)
					if err != nil {
						log.Fatal(err)
					}
					out <- event

				} else {
					log.Fatal("Expected event data, got: " + line)
				}
			} else if line == "\n" {
				// next
			} else {
				log.Fatal("Expected event name, got: " + line)
			}
		}
		resp.Body.Close()
	}()

	return out
}
