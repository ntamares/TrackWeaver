package parser

type Song struct {
	Time   string // "9p", "10p"
	Title  string
	Artist string
}

type Playlist struct {
	Date  string
	Songs []Song
}
