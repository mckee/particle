# particle
--
    import "github.com/mckee/particle"


## Usage

```go
const URL string = "https://api.particle.io/v1/devices/events/"
```

#### func  Subscribe

```go
func Subscribe(eventPrefix string, token string) <-chan Event
```
subscribes to a particle.io event stream and returns a channel to receive them

#### type Event

```go
type Event struct {
	Name string
	Data struct {
		Data      string `json:"data"`
		TTL       string `json:"ttl"`
		Timestamp string `json:"published_at"`
		CoreID    string `json:"coreid"`
	}
}
```
