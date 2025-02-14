package girm

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	invalidIdMsg = "invalid id"
)

type DB[T any] struct {
	db *gorm.DB
}

func NewDb[T any](db *gorm.DB) *DB[T] {
	return &DB[T]{
		db: db,
	}
}

func (d *DB[T]) Insert(c *gin.Context) {
	operation(c, func(es ...*T) error {
		return bulkOperate(es, d.db.Create)
	})
}

func (d *DB[T]) SelectAll(c *gin.Context) {
	var es []*T
	if err := d.db.Find(&es).Error; err != nil {
		JsonFail(c, err.Error())
		return
	}
	JsonOK(c, es)
}

func (d *DB[T]) SelectById(c *gin.Context) {
	var err error
	var e *T
	defer func() {
		if err != nil {
			JsonFail(c, err.Error())
		} else {
			JsonOK(c, e)
		}
	}()

	id, ok := c.Params.Get("id")
	if !ok {
		err = errors.New(invalidIdMsg)
		return
	}
	err = d.db.First(e, "id = ?", id).Error
}

//conditionFields: key:field in query;value:field in db. like key=="nodeName" value=="node_name"
func (d *DB[T]) SelectByConditions(c *gin.Context, conditionFields map[string]string) {
	var err error
	var es []*T
	defer func() {
		if err != nil {
			JsonFail(c, err.Error())
		} else {
			JsonOK(c, es)
		}
	}()

	conditions := make(map[string]any)
	for queryField, dbField := range conditionFields {
		queryValue := c.Query(queryField)
		conditions[dbField] = queryValue
	}
	err = d.db.Where(conditions).Find(es).Error
}

func (d *DB[T]) SelectByPage(c *gin.Context, where func(db *gorm.DB) *gorm.DB) {
	var err error
	page := struct {
		PageNum  int
		PageSize int
		Total    int
		Data     []*T
	}{
		PageNum: 1,
	}

	defer func() {
		if err != nil {
			JsonFail(c, err.Error())
		} else {
			JsonOK(c, &page)
		}
	}()

	sPageSize := c.Query("pageSize")
	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil || pageSize <= 0 {
		return
	}
	page.PageSize = pageSize

	sPageNum := c.Query("pageNum")
	pageNum, err := strconv.Atoi(sPageNum)
	if err != nil || pageNum <= 0 {
		return
	}
	page.PageNum = pageNum

	var count int64
	offset := (pageNum - 1) * pageSize
	if where != nil {
		where(d.db.Model(new(T))).Count(&count)
		where(d.db).Offset(offset).Limit(pageSize).Find(&page.Data)
	} else {
		d.db.Model(new(T)).Count(&count)
		d.db.Offset(offset).Limit(pageSize).Find(&page.Data)
	}

	page.Total = int(count)
}

func (d *DB[T]) Save(c *gin.Context) {
	operation(c, func(es ...*T) error {
		return bulkOperate(es, d.db.Save)
	})
}

func (d *DB[T]) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		JsonFail(c, err.Error())
		return
	}

	var es []*T
	if err := d.db.Where("id IN ?", ids).Find(&es).Error; err != nil {
		JsonFail(c, err.Error())
		return
	}
	if err := d.db.Delete(&es).Error; err != nil {
		JsonFail(c, err.Error())
		return
	}
	JsonOK(c, nil)
}
