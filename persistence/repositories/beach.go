package repositories

import (
	"fmt"

	"github.com/mendelgusmao/cetesb-telegram-bot/persistence"
	"github.com/mendelgusmao/cetesb-telegram-bot/persistence/models"
	"gorm.io/gorm/clause"
)

type BeachRepository struct {
	persistence *persistence.Persistence
}

func NewBeachRepository(persistence persistence.Persistence) *BeachRepository {
	return &BeachRepository{
		persistence: &persistence,
	}
}

func (r *BeachRepository) All() ([]models.Beach, error) {
	var beaches []models.Beach

	result := r.persistence.Find(&beaches)

	if result.Error != nil {
		return beaches, fmt.Errorf("[persistence.repositories.BeachRepository.All] %v", result.Error)
	}

	return beaches, nil
}

func (r *BeachRepository) CreateOrUpdate(beaches []models.Beach) error {
	tx := r.persistence.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "city"},
			{Name: "name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"city", "name", "place", "quality"}),
	}).Create(&beaches)

	if tx.Error != nil {
		return fmt.Errorf("[persistence.repositories.BeachRepository.CreateOrUpdateBeach] %v", tx.Error)
	}

	return nil
}
