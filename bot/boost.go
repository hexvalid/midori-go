package bot

import "time"

type Boost struct {
	ID     string    `json:"id"`
	Type   string    `json:"type"`
	Expire time.Time `json:"expire"`
}
