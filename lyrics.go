package musicmax

import "time"

type Lyrics struct {
	Id     string `json:"-"`
	Text   string `json:"text"`
	SongId string `json:"-"`
	Song   string `json:"-"`
	Group  string `json:"-"`
}

type LyricsResponse struct {
	Group string   `json:"group"`
	Song  string   `json:"song"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Text  []string `json:"text"`
}

type LyricsAPIResponse struct {
	Release *time.Time `json:"releaseDate"`
	Link    string     `json:"link"`
	Text    string     `json:"text"`
}
