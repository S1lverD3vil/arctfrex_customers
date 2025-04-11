package grouprole

import (
	"arctfrex-customers/internal/base"
	"log"

	"gorm.io/gorm"
)

func SeedGroupRoles(db *gorm.DB) {
	roleGroups := []GroupRole{
		{ID: "Marketing", Name: "Marketing", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "CRM", Name: "Customer Relationship Management", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "Director", Name: "Director", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "Finance", Name: "Finance", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "Settlement", Name: "Settlement", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "WPB Verifikator", Name: "WPB Verifikator", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "UKK APUPPT", Name: "UKK APUPPT", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
	}

	for _, rg := range roleGroups {
		if err := db.FirstOrCreate(&rg, GroupRole{ID: rg.ID}).Error; err != nil {
			log.Printf("Failed to insert RoleGroup %s: %v", rg.ID, err)
		}
	}
}

func strPtr(s string) *string {
	return &s
}
