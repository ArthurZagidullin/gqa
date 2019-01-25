package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gqa/config"
	"gqa/lib"
	"net"
	"sync"
	"time"
)

type Aggregator struct {
	conf    config.AgregatorConfig
	storage config.Saver
	conn    net.Conn
	buffer  map[string][]int
	group   *sync.WaitGroup
	mu      *sync.Mutex
}

func NewAggregator(group *sync.WaitGroup, conf config.AgregatorConfig, storage config.Saver, conn net.Conn) *Aggregator {
	return &Aggregator{
		conf:    conf,
		storage: storage,
		conn:    conn,
		group:   group,
		buffer:  make(map[string][]int),
		mu:      &sync.Mutex{},
	}
}

func (a *Aggregator) Run() {
	defer a.conn.Close()
	defer a.group.Done()

	go a.calcAvg()

	input := bufio.NewScanner(a.conn)

	for {

		if ok := input.Scan(); !ok {
			break
		}
		newMess := lib.Message{}
		if err := json.Unmarshal(input.Bytes(), &newMess); err != nil {
			fmt.Printf("Error message unmarshaling: %s \n", err)
			continue
		}

		for _, sub_id := range a.conf.Sub_ids {
			if sub_id == newMess.Id {
				a.mu.Lock()
				a.buffer[newMess.Id] = append(a.buffer[newMess.Id], newMess.Value)

				fmt.Printf("aggregate: %+v. Buffer size: %d \n", newMess, len(a.buffer[newMess.Id]))
				a.mu.Unlock()
				break
			}
		}

	}
}

func (a *Aggregator) calcAvg() {
	ticker := time.NewTicker(time.Second * time.Duration(a.conf.Agregate_period_s))

	for _ = range ticker.C {
		a.mu.Lock()

		for id, buf := range a.buffer {
			sum := 0
			for i := range buf {
				sum += buf[i]
			}
			avg := sum / len(buf)

			go a.storage.Save(id, avg)
			delete(a.buffer, id)
		}
		a.mu.Unlock()
	}

}
