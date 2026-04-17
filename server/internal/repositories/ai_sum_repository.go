package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var AISumRepository = newAISumRepository()

type aiSumRepository struct{}

func newAISumRepository() *aiSumRepository {
	return &aiSumRepository{}
}

func (r *aiSumRepository) GetByID(db *gorm.DB, id int64) *models.AISum {
	if id <= 0 {
		return nil
	}
	var result models.AISum
	err := db.First(&result, id).Error
	if err != nil {
		return nil
	}
	return &result
}

func (r *aiSumRepository) GetValidSummary(db *gorm.DB, topicId int64) *models.AISum {
	var result models.AISum
	err := db.Where("topic_id = ? AND sum_valid = ?", topicId, "Y").
		Order("sum_time DESC").
		First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

func (r *aiSumRepository) Create(db *gorm.DB, entity *models.AISum) error {
	return db.Create(entity).Error
}

func (r *aiSumRepository) InvalidatePreviousSummaries(db *gorm.DB, topicId int64) error {
	return db.Model(&models.AISum{}).
		Where("topic_id = ? AND sum_valid = ?", topicId, "Y").
		Update("sum_valid", "N").Error
}

func (r *aiSumRepository) FindByTopicId(db *gorm.DB, topicId int64) ([]models.AISum, error) {
	var list []models.AISum
	err := db.Where("topic_id = ?", topicId).
		Order("sum_time DESC").
		Find(&list).Error
	return list, err
}

func (r *aiSumRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) ([]models.AISum, *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *aiSumRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) ([]models.AISum, *sqls.Paging) {
	var list []models.AISum
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.AISum{})

	paging := &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return list, paging
}

func (r *aiSumRepository) Update(db *gorm.DB, entity *models.AISum) error {
	return db.Save(entity).Error
}

func (r *aiSumRepository) Delete(db *gorm.DB, id int64) error {
	return db.Delete(&models.AISum{}, id).Error
}
