package girm

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var (
	maxOprCount = 500
)

func bulkOperate[T any](es []*T, f func(any) *gorm.DB) error {
	if len(es) < maxOprCount {
		return f(es).Error
	}

	segs := lo.Chunk(es, maxOprCount)
	for _, seg := range segs {
		if err := f(seg).Error; err != nil {
			return err
		}
	}
	return nil
}
