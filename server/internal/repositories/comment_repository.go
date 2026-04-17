package repositories

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var CommentRepository = newCommentRepository()

func newCommentRepository() *commentRepository {
	return &commentRepository{}
}

type commentRepository struct {
}

func (r *commentRepository) Get(db *gorm.DB, id int64) *models.Comment {
	ret := &models.Comment{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *commentRepository) Take(db *gorm.DB, where ...interface{}) *models.Comment {
	ret := &models.Comment{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *commentRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Comment) {
	cnd.Find(db, &list)
	return
}

func (r *commentRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Comment {
	ret := &models.Comment{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *commentRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Comment, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *commentRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Comment, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Comment{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *commentRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Comment{})
}

func (r *commentRepository) Create(db *gorm.DB, t *models.Comment) (err error) {
	err = db.Create(t).Error
	return
}

func (r *commentRepository) Update(db *gorm.DB, t *models.Comment) (err error) {
	err = db.Save(t).Error
	return
}

func (r *commentRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Comment{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *commentRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *commentRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Comment{}, "id = ?", id)
}

// GetUnsummarizedComments 获取未总结的评论
func (r *commentRepository) GetUnsummarizedComments(db *gorm.DB, topicId int64) ([]models.Comment, int) {
	var comments []models.Comment
	err := db.Where("entity_type = ? AND entity_id = ? AND ai_sum IS NULL",
		constants.EntityTopic, topicId).
		Preload("User").
		Order("id ASC").
		Find(&comments).Error

	if err != nil {
		return nil, 0
	}

	// 计算总字符数
	totalChars := 0
	for _, comment := range comments {
		totalChars += len(comment.Content)
	}

	return comments, totalChars
}

// MarkCommentsAsSummarized 标记评论为已总结
func (r *commentRepository) MarkCommentsAsSummarized(db *gorm.DB, topicId int64) error {
	return db.Model(&models.Comment{}).
		Where("entity_type = ? AND entity_id = ? AND ai_sum IS NULL",
			constants.EntityTopic, topicId).
		Update("ai_sum", time.Now()).Error
}
