package application

import (
	"errors"

	"github.com/quynhanh/internship-tracker/internal/model"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Create
func (s *Service) Create(userID uint, req CreateRequest) error {
	app := model.Application{
		UserID:      userID,
		CompanyName: req.CompanyName,
		Position:    req.Position,
		Status:      req.Status,
		Notes:       req.Notes,
	}
	return s.DB.Create(&app).Error
}

// Update
func (s *Service) UpdateStatus(userID, appID uint, status string) error {
	var app model.Application
	err := s.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error
	if err != nil {
		return errors.New("Application not found")
	}

	app.Status = status
	return s.DB.Save(&app).Error
}

// Delete
func (s *Service) Delete(userID, appID uint) error {
	res := s.DB.Where("id = ? AND user_id = ?", appID, userID).Delete(&model.Application{})
	if res.RowsAffected == 0 {
		return errors.New("Application not found")
	}
	return res.Error
}

// List returns a paginated list of applications for a specific user
func (s *Service) List(
	userID uint,
	page int,
	limit int,
	status string,
) ([]model.Application, int64, error) {
	// Slice to store the queried applications
	var apps []model.Application

	// Total number of records matching the query (without pagination)
	var total int64

	// Ensure page number is at least 1
	if page < 1 {
		page = 1
	}
	// Validate number to prevent zero, negative, or excessively large queries
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Calculate the SQL offset based on the current page and limit
	offset := (page - 1) * limit

	// Build base query scoped to the authenticated user
	query := s.DB.Model(&model.Application{}).Where("user_id = ?", userID)

	// Apply status filer if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records matching the query (without limit/offset)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated application records
	// Ordered by creation of time in descending order (newest first)
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}
