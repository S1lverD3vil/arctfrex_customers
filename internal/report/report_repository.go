package report

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
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
