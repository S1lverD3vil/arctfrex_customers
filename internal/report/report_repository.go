package report

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (rr *reportRepository) GetActiveReports() (*[]Report, error) {
	var reports []Report
	queryParams := Report{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := rr.db.Find(&reports, &queryParams).Error; err != nil {
		return nil, err
	}

	return &reports, nil
}

func (rr *reportRepository) SaveGroupUserLogins(reportGroupUserLogins []ReportGroupUserLogins) error {
	// return rr.db.Save(reportGroupUserLogins).Error
	if len(reportGroupUserLogins) == 0 {
		return nil // No data to insert
	}

	// Use ON CONFLICT DO NOTHING to skip duplicates
	result := rr.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&reportGroupUserLogins)

	return result.Error
}
