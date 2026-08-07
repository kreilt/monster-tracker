package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Flavor struct {
	Name        string `json:"n"`
	Lineup      string `json:"l"`
	Description string `json:"t"`
	Rare        int    `json:"r"`
	Region      string `json:"g"`
	Color       string `json:"c"`
	Status      string `json:"s"`
}

var rarityMap = map[int]string{
	1: "ordinary",
	2: "unusual",
	3: "rare",
	4: "very rare",
	5: "legendary",
}

func esc(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func main() {
	var flavors []Flavor
	data, err := os.ReadFile("flavors-data.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &flavors); err != nil {
		log.Fatal(err)
	}

	var b strings.Builder
	b.WriteString("-- +goose Up\n\n")
	b.WriteString("INSERT INTO flavors(title, lineup, description, rare, region, color, status)\n VALUES\n")

	rows := make([]string, 0, len(flavors))

	for _, f := range flavors {
		rarity, ok := rarityMap[f.Rare]
		if !ok {
			log.Fatalf("invalid rarity = %d name = %q", f.Rare, f.Name)
		}

		status := f.Status
		if status == "" {
			status = "current"
		}

		rows = append(rows, fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			esc(f.Name),
			esc(f.Lineup),
			esc(f.Description),
			rarity,
			esc(f.Region),
			esc(f.Color),
			esc(status)))
	}

	b.WriteString(strings.Join(rows, ",\n"))
	b.WriteString(";\n")

	b.WriteString("\n-- +goose Down\n\n")
	b.WriteString("DELETE FROM flavors\n WHERE title IN (")

	titles := make([]string, 0, len(flavors))

	for _, f := range flavors {
		titles = append(titles, fmt.Sprintf("'%s'", esc(f.Name)))
	}

	b.WriteString(strings.Join(titles, ",\n"))
	b.WriteString(");\n")

	if err := os.WriteFile("../migrations/00002_seed_flavors.sql", []byte(b.String()), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK, обработано строк:", len(rows))
}
