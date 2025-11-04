package proxy

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

var (
	rnd *rand.Rand
	mtx sync.Mutex
)

func initRandSeed() {
	source := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(source)
}

func PickService(c *api.Client, name string) (*api.AgentService, error) {
	entries, _, err := c.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("Counsul search failed: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("%s exists no instance", name)
	}

	mtx.Lock()
	idx := rnd.Intn(len(entries))
	mtx.Unlock()

	return entries[idx].Service, nil
}
