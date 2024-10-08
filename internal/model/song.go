package model

import "time"

type Song struct {
	Id      string     `json:"id,omitempty"`
	Group   string     `json:"group"`
	Name    string     `json:"song"`
	Release *time.Time `json:"releaseDate,omitempty"`
	Link    string     `json:"link,omitempty"`
	Lyrics  *Lyrics    `json:"-"`
}

type SongRequest struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type SongPatchRequest struct {
	Group   *string `json:"group,omitempty"`
	Name    *string `json:"song,omitempty"`
	Release *string `json:"releaseDate,omitempty"`
	Text    *string `json:"text,omitempty"`
	Link    *string `json:"link,omitempty"`
}

type SongResponse struct {
	Id      string `json:"id"`
	Group   string `json:"group"`
	Name    string `json:"song"`
	Release string `json:"releaseDate,omitempty"`
	Link    string `json:"link,omitempty"`
}

type SongsResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Songs []SongResponse `json:"songs"`
}
