package model

import (
	"next-terminal/pkg/global"
	"next-terminal/pkg/utils"
)

// 密码
const Custom = "custom"

// 密钥
const PrivateKey = "private-key"

type Credential struct {
	ID         string         `gorm:"primary_key" json:"id"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	PrivateKey string         `json:"privateKey"`
	Passphrase string         `json:"passphrase"`
	Created    utils.JsonTime `json:"created"`
	Owner      string         `json:"owner"`
}

func (r *Credential) TableName() string {
	return "credentials"
}

type CredentialVo struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Username  string         `json:"username"`
	Created   utils.JsonTime `json:"created"`
	Owner     string         `json:"owner"`
	OwnerName string         `json:"ownerName"`
}

func FindAllCredential() (o []Credential, err error) {
	err = global.DB.Find(&o).Error
	return
}

func FindPageCredential(pageIndex, pageSize int, name, owner string) (o []CredentialVo, total int64, err error) {
	db := global.DB
	db = db.Table("credentials").Select("credentials.id,credentials.name,credentials.type,credentials.username,credentials.owner,credentials.created,users.nickname as owner_name").Joins("left join users on credentials.owner = users.id")
	if len(name) > 0 {
		db = db.Where("credentials.name like ?", "%"+name+"%")
	}
	if len(owner) > 0 {
		db = db.Where("credentials.owner = ?", owner)
	}

	err = db.Order("credentials.created desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&o).Count(&total).Error
	if o == nil {
		o = make([]CredentialVo, 0)
	}
	return
}

func CreateNewCredential(o *Credential) (err error) {
	if err = global.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindCredentialById(id string) (o Credential, err error) {
	err = global.DB.Where("id = ?", id).First(&o).Error
	return
}

func UpdateCredentialById(o *Credential, id string) {
	o.ID = id
	global.DB.Updates(o)
}

func DeleteCredentialById(id string) {
	global.DB.Where("id = ?", id).Delete(&Credential{})
}

func CountCredential() (total int64, err error) {
	err = global.DB.Find(&Credential{}).Count(&total).Error
	return
}
