package services

import (
	"errors"
	"time"

	"github.com/vhybZApp/api/database"
	"gorm.io/gorm"
)

type TokenQuotaService struct {
	db *gorm.DB
}

func NewTokenQuotaService(db *gorm.DB) *TokenQuotaService {
	return &TokenQuotaService{db: db}
}

// GetUserQuota returns the daily token quota for a user
func (s *TokenQuotaService) GetUserQuota(userID string) (*database.DBTokenQuota, error) {
	var quota database.DBTokenQuota
	result := s.db.Where("user_id = ?", userID).First(&quota)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create default quota if not exists
			quota = database.DBTokenQuota{
				UserID:     userID,
				DailyQuota: 100000, // Default 100k tokens per day
			}
			if err := s.db.Create(&quota).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}
	return &quota, nil
}

// GetDailyUsage returns the token usage for a user on a specific date
func (s *TokenQuotaService) GetDailyUsage(userID string, date time.Time) (*database.DBTokenUsage, error) {
	var usage database.DBTokenUsage
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	result := s.db.Where("user_id = ? AND date = ?", userID, startOfDay).First(&usage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new usage record if not exists
			usage = database.DBTokenUsage{
				UserID: userID,
				Date:   startOfDay,
				Tokens: 0,
			}
			if err := s.db.Create(&usage).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}
	return &usage, nil
}

// UpdateUsage updates the token usage for a user
func (s *TokenQuotaService) UpdateUsage(userID string, tokens int) error {
	now := time.Now()

	// Get or create daily usage
	usage, err := s.GetDailyUsage(userID, now)
	if err != nil {
		return err
	}

	// Get user quota
	quota, err := s.GetUserQuota(userID)
	if err != nil {
		return err
	}

	// Check if usage exceeds quota
	if usage.Tokens+tokens > quota.DailyQuota {
		return errors.New("daily token quota exceeded")
	}

	// Update usage
	usage.Tokens += tokens
	return s.db.Save(usage).Error
}

// ResetDailyUsage resets the token usage for all users at the start of a new day
func (s *TokenQuotaService) ResetDailyUsage() error {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Delete all usage records from previous days
	return s.db.Where("date < ?", startOfDay).Delete(&database.DBTokenUsage{}).Error
}
