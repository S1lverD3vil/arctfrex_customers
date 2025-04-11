package role

import (
	"arctfrex-customers/internal/base"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []Role{
		// Marketing Hierarchy
		{ID: "HM", Name: "Head of Marketing", CommissionRate: 0.01, GroupRoleID: "marketing", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "SBM", Name: "Senior Business Manager", CommissionRate: 0.02, GroupRoleID: "marketing", ParentRoleID: strPtr("hm"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "BM", Name: "Business Manager", CommissionRate: 0.03, GroupRoleID: "marketing", ParentRoleID: strPtr("sbm"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "ABM", Name: "Assistant Business Manager", CommissionRate: 0.05, GroupRoleID: "marketing", ParentRoleID: strPtr("bm"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "MKT", Name: "Marketing", CommissionRate: 0.07, GroupRoleID: "marketing", ParentRoleID: strPtr("abm"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "IB", Name: "Freelance", CommissionRate: 0.10, GroupRoleID: "marketing", ParentRoleID: strPtr("mkt"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},

		// CRM Hierarchy
		{ID: "CRM-Manager", Name: "CRM Manager", GroupRoleID: "crm", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "CRM-Support", Name: "CRM Support", GroupRoleID: "crm", ParentRoleID: strPtr("crm-manager"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "CRM-Agent", Name: "CRM Agent", GroupRoleID: "crm", ParentRoleID: strPtr("crm-support"), BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},

		// Independent Roles
		{ID: "DIR", Name: "Director", GroupRoleID: "director", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "FIN", Name: "Finance", GroupRoleID: "finance", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "SETT", Name: "Settlement", GroupRoleID: "settlement", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "WPB", Name: "WPB Verifikator", GroupRoleID: "wpb_verifikator", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
		{ID: "UKK", Name: "UKK APUPPT", GroupRoleID: "ukk_apuppt", BaseModel: base.BaseModel{IsActive: true, CreatedBy: "system"}},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, Role{ID: role.ID}).Error; err != nil {
			log.Printf("Failed to insert Role %s: %v", role.ID, err)
		}
	}
}

func strPtr(s string) *string {
	return &s
}
