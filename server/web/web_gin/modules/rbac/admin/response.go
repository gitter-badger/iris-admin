package admin

import (
	"regexp"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Response struct {
	orm.Model
	BaseAdmin
	AuthorityIds []string `gorm:"-" json:"authorityIds"`
}

func (res *Response) ToString() {
	if res.Avatar.HeaderImg == "" {
		return
	}
	re := regexp.MustCompile("^http")
	if !re.MatchString(res.Avatar.HeaderImg) {
		res.Avatar.HeaderImg = str.Join("http://127.0.0.1:8085/upload/", res.Avatar.HeaderImg)
	}
}

type LoginResponse struct {
	orm.ReqId
	Password string `json:"password"`
}

func (res *Response) First(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(&Admin{}).Scopes(scopes...).First(res).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取失败", zap.String("First()", err.Error()))
		return err
	}
	return nil
}

// Paginate 分页
type PageResponse struct {
	Item []*Response
}

func (res *PageResponse) Paginate(db *gorm.DB, pageScope func(db *gorm.DB) *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	db = db.Model(&Admin{})
	var count int64
	err := db.Scopes(scopes...).Count(&count).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取总数失败", zap.String("Count()", err.Error()))
		return count, err
	}
	err = db.Scopes(pageScope).Find(&res.Item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取分页数据失败", zap.String("Find()", err.Error()))
		return count, err
	}

	return count, nil
}

func (res *PageResponse) Find(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	db = db.Model(&Admin{})
	err := db.Scopes(scopes...).Find(&res.Item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取数据失败", zap.String("Find()", err.Error()))
		return err
	}

	return nil
}