package models

import "time"

type Url struct {
	Link    string
	Expires time.Duration
}
