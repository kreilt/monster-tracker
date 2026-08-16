package model

type Rarity string

type Statuses string

const (
	RarityOrdinary  Rarity = "ordinary"
	RarityUnusual   Rarity = "unusual"
	RarityRare      Rarity = "rare"
	RarityVeryRare  Rarity = "very rare"
	RarityLegendary Rarity = "legendary"

	StatusesDisc    Statuses = "disc"
	StatusesLim     Statuses = "lim"
	StatusesNew     Statuses = "new"
	StatusesReg     Statuses = "reg"
	StatusesAlc     Statuses = "alc"
	StatusesCurrent Statuses = "current"
)

func IsValidRarity(s string) bool {
	switch Rarity(s) {
	case RarityOrdinary, RarityUnusual, RarityRare, RarityVeryRare, RarityLegendary:
		return true
	}
	return false
}

func IsValidStatuses(s string) bool {
	switch Statuses(s) {
	case StatusesDisc, StatusesLim, StatusesNew, StatusesReg, StatusesAlc, StatusesCurrent:
		return true
	}
	return false
}

type Flavor struct {
	FlavorID    int      `json:"flavor_id"`
	Title       string   `json:"title"`
	Lineup      string   `json:"lineup"`
	Description *string  `json:"description"`
	Rare        Rarity   `json:"rare"`
	Region      string   `json:"region"`
	Color       string   `json:"color"`
	Status      Statuses `json:"status"`
	Photo       *string  `json:"photo"`
}
