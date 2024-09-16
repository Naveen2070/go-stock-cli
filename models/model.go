package models

type Stock struct {
	Ticker       string
	Gap          float64
	OpeningPrice float64
}

type Position struct {
	EntryPrice       float64
	Shares           int
	TakeProfilePrice float64
	StopLossPrice    float64
	Profit           float64
}

type Article struct {
	PublishedOn string
	Headline    string
}

type Selection struct {
	Ticker string
	Position
	Articles []Article
}
