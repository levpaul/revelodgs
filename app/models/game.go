package models

type Game struct {
	GameId    int
	Name      string
	ShortDesc string
	LongDesc  string
	AmiId     string `json:"-"`
	Type      string
}

const (
	GameTypeSteam     string = "steam"
	GameTypeMinecraft string = "minecraft"
)
