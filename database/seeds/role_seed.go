package seeds

import (
	"go-article/internal/domain/model"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []model.Role{
		{Name: "Admin"},
		{Name: "User"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Fatalf("%s: %v", err.Error(), err)
		} else {
			log.Printf("[SeedRoles] Seeded role: %s", role.Name)
		}
	}
}
