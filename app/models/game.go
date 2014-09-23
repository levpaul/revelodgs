package models

type Game struct {
	GameId      int
	Name        string
	Description string
	AmiId       string `json:"-"`
	Type        string
}

const (
	GameTypeSteam     string = "steam"
	GameTypeMinecraft string = "minecraft"
)
