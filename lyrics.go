package musicmax

type Lyrics struct {
	Song
	Id   int    `json:"-"`
	Text string `json:"text"`
}
