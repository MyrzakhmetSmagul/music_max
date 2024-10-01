package musicmax

import "time"

type Song struct {
	Id      int        `json:"id,omitempty"`
	Group   string     `json:"group"`
	Name    string     `json:"song"`
	Release *time.Time `json:"releaseDate"`
}
