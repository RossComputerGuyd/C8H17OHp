package model

import (
	"next-terminal/pkg/global"
	"next-terminal/pkg/utils"
	"strings"
)

type Asset struct {
	ID           string         `gorm:"primary_key " json:"id"`
	Name         string         `json:"name"`
	IP           string         `json:"ip"`
	Protocol     string         `json:"protocol"`
	Port         int            `json:"port"`
	AccountType  string         `json:"accountType"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	CredentialId string         `json:"credentialId"`
	PrivateKey   string         `json:"privateKey"`
	Passphrase   string         `json:"passphrase"`
	Description  string         `json:"description"`
	Active       bool           `json:"active"`
	Created      utils.JsonTime `json:"created"`
	Tags         string         `json:"tags"`
	Owner        string         `json:"owner"`
}

type AssetVo struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	IP        string         `json:"ip"`
	Protocol  string         `json:"protocol"`
	Port      int            `json:"port"`
	Active    bool           `json:"active"`
	Created   utils.JsonTime `json:"created"`
	Tags      string         `json:"tags"`
	Owner     string         `json:"owner"`
	OwnerName string         `json:"ownerName"`
}

func (r *Asset) TableName() string {
	return "assets"
}

func FindAllAsset() (o []Asset, err error) {
	err = global.DB.Find(&o).Error
	return
}

func FindAssetByConditions(protocol string) (o []Asset, err error) {
	db := global.DB

	if len(protocol) > 0 {
		db = db.Where("protocol = ?", protocol)
	}
	err = db.Find(&o).Error
	return
}

func FindPageAsset(pageIndex, pageSize int, name, protocol, tags string) (o []AssetVo, total int64, err error) {
	db := global.DB
	db = db.Table("assets").Select("assets.id,assets.name,assets.ip,assets.port,assets.protocol,assets.active,assets.owner,assets.created, users.nickname as creator_name").Joins("left join users on assets.owner = users.id")
	if len(name) > 0 {
		db = db.Where("assets.name like ?", "%"+name+"%")
	}

	if len(protocol) > 0 {
		db = db.Where("assets.protocol = ?", protocol)
	}

	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for i := range tagArr {
			db = db.Where("find_in_set(?, assets.tags)", tagArr[i])
		}
	}

	err = db.Order("assets.created desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&o).Count(&total).Error

	if o == nil {
		o = make([]AssetVo, 0)
	}
	return
}

func CreateNewAsset(o *Asset) (err error) {
	if err = global.DB.Create(o).Error; err != nil {
		return err
	}
	return nil
}

func FindAssetById(id string) (o Asset, err error) {
	err = global.DB.Where("id = ?", id).First(&o).Error
	return
}

func UpdateAssetById(o *Asset, id string) {
	o.ID = id
	global.DB.Updates(o)
}

func UpdateAssetActiveById(active bool, id string) {
	sql := "update assets set active = ? where id = ?"
	global.DB.Exec(sql, active, id)
}

func DeleteAssetById(id string) {
	global.DB.Where("id = ?", id).Delete(&Asset{})
}

func CountAsset() (total int64, err error) {
	err = global.DB.Find(&Asset{}).Count(&total).Error
	return
}

func FindAssetTags() (o []string, err error) {
	var assets []Asset
	err = global.DB.Not("tags = ?", "").Find(&assets).Error
	if err != nil {
		return nil, err
	}

	o = make([]string, 0)

	for i := range assets {
		if len(assets[i].Tags) == 0 {
			continue
		}
		split := strings.Split(assets[i].Tags, ",")

		o = append(o, split...)
	}

	return utils.Distinct(o), nil
}
