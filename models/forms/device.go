package forms

import (
	"time"
)

type ReturnByCreatedAt []*Return

func (p ReturnByCreatedAt) Len() int           { return len(p) }
func (p ReturnByCreatedAt) Less(i, j int) bool { return p[i].CreatedAt.After(p[j].CreatedAt) }
func (p ReturnByCreatedAt) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Return struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
