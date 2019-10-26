package sys

import (
	"time"

	"github.com/it234/goapp/internal/pkg/models/db"

	"github.com/jinzhu/gorm"
)

// 后台用户
type Admins struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment;" json:"id" form:"id"`                                             // 主键
	Memo      string    `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                                                          //备注
	UserName  string    `gorm:"column:user_name;size:32;unique_index:uk_admins_user_name;not null;" json:"user_name" form:"user_name"` // 用户名
	Password  string    `gorm:"column:password;type:char(32);not null;" json:"password" form:"password"`                               // 密码(sha1(md5(明文))加密)
	Salt      string    `gorm:"column:salt;type:char(8);not null;" json:"salt" form:"salt"`                                            // 密码(盐)
	Status    uint8     `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`                                   // 状态(1:正常 2:未激活 3:暂停使用)
	OrgID     string    `gorm:"column:org_id;type:tinyint(1);" json:"org_id" form:"org_id"`                                            // 机构号
	RealName  string    `gorm:"column:real_name;size:32;" json:"real_name" form:"real_name"`                                           // 真实姓名
	CertNum   string    `gorm:"column:cert_num;type:char(20);" json:"cert_num" form:"cert_num"`                                        // 证件号
	Sex       string    `gorm:"column:sex;type:tinyint(1);" json:"sex" form:"sex"`                                                     // 性别
	Phone     string    `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`                                                 // 手机号
	Email     string    `gorm:"column:email;size:64;" json:"email" form:"email"`                                                       // 邮箱
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;" json:"created_at" form:"created_at"`                         // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;" json:"updated_at" form:"updated_at"`                         // 更新时间
}

// 表名
func (Admins) TableName() string {
	return TableName("admins")
}

// 添加前
func (m *Admins) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *Admins) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 删除用户及关联数据
func (Admins) Delete(adminsids []uint64) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", adminsids).Delete(&Admins{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("admins_id in (?)", adminsids).Delete(&AdminsRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
