package models

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/gorm"
)

const (
	SeedFile = "data/seed.json"
)

type (
	Seeder struct {
		db *gorm.DB
	}
	SeedData struct {
		Plans []Plan `json:"plans"`
	}
)

func NewSeeder(db *gorm.DB) (seeder *Seeder) {
	seeder = &Seeder{db: db}
	return
}

func (seeder *Seeder) SeedPlans(seedData SeedData) (err error) {
	log.Println("Seeding plans...", len(seedData.Plans))
	for _, plan := range seedData.Plans {
		if err = seeder.db.FirstOrCreate(&plan, Plan{Name: plan.Name}).Error; err != nil {
			return
		}
	}
	return
}

func (seeder *Seeder) Seed() (err error) {
	seedData := SeedData{}
	seedDataBytes, err := os.ReadFile(SeedFile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(seedDataBytes, &seedData); err != nil {
		return
	}
	if err = seeder.SeedPlans(seedData); err != nil {
		return
	}
	return
}
