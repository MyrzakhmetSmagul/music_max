package musicmax

import "time"

type Song struct {
	Id      string     `json:"id,omitempty"`
	Group   string     `json:"group"`
	Name    string     `json:"song"`
	Release *time.Time `json:"releaseDate,omitempty"`
	Lyrics  *Lyrics    `json:"-"`
}

type SongRequest struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type SongResponse struct {
	Id      string     `json:"id"`
	Group   string     `json:"group"`
	Name    string     `json:"song"`
	Release *time.Time `json:"releaseDate,omitempty"`
}

type SongsResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Songs []SongResponse `json:"songs"`
}
