package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Provider represents a telecomm provider
type Provider struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"size:90;not null"`
}

// NumberingType mobile or fixed
type NumberingType struct {
	ID   uint   `json:"id"`
	Code string `json:"name" gorm:"size:5;not null"`
}

// Numbering from ift
type Numbering struct {
	ID           uint   `json:"id"`
	Region       uint   `json:"region" gorm:"not null"`
	Area         uint   `json:"area" gorm:"not null"`
	Lada         uint   `json:"lada" gorm:"not null"`
	Series       uint   `json:"series" gorm:"not null"`
	Start        uint   `json:"start" gorm:"not null"`
	End          uint   `json:"end" gorm:"not null"`
	State        string `json:"state" gorm:"size:5;not null"`
	Municipality string `json:"municipality" gorm:"size:80;not null"`
	ProviderID   uint   `gorm:"not null"`
	TypeID       uint   `gorm:"not null"`
}

func clean(val string) string {
	return strings.Trim(val, " ")
}

func parseInt(val string) uint {
	a, err := strconv.Atoi(val)

	if err != nil {
		panic(err)
	}

	return uint(a)
}

func main() {
	db, err := gorm.Open("postgres", "user=postgres dbname=ift sslmode=disable password=postgres")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SingularTable(true)
	db.AutoMigrate(&Provider{})
	db.AutoMigrate(&NumberingType{})
	db.AutoMigrate(&Numbering{})

	f, err := os.Open("pnn_Publico_04_02_2018.csv")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()

	if err != nil {
		panic(err)
	}

	for _, row := range rows[1:] {
		provider := new(Provider)
		t := new(NumberingType)

		// Create Provider

		name := clean(row[14])

		if db.First(provider, "name = ?", name).RecordNotFound() {
			provider = &Provider{
				Name: name,
			}

			db.Create(provider)
		}

		// Create NumberingType

		code := clean(row[12])

		if db.First(t, "code = ?", code).RecordNotFound() {
			t = &NumberingType{
				Code: code,
			}

			db.Create(t)
		}

		db.Create(&Numbering{
			Region:       parseInt(row[5]),
			Area:         parseInt(row[6]),
			Lada:         parseInt(row[7]),
			Series:       parseInt(row[8]),
			Start:        parseInt(row[9]),
			End:          parseInt(row[10]),
			State:        row[3],
			Municipality: row[2],
			ProviderID:   provider.ID,
			TypeID:       t.ID,
		})
	}
}
