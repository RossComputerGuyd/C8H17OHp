package model

import (
	"next-terminal/pkg/config"
	"next-terminal/pkg/utils"
)

const (
	Connected    = "connected"
	Disconnected = "disconnected"
	NoConnect    = "no_connect"
)

type Session struct {
	ID               string         `gorm:"primary_key" json:"id"`
	Protocol         string         `json:"protocol"`
	IP               string         `json:"ip"`
	Port             int            `json:"port"`
	ConnectionId     string         `json:"connectionId"`
	AssetId          string         `json:"assetId"`
	Username         string         `json:"username"`
	Password         string         `json:"password"`
	Creator          string         `json:"creator"`
	ClientIP         string         `json:"clientIp"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	Status           string         `json:"status"`
	ConnectedTime    utils.JsonTime `json:"connectedTime"`
	DisconnectedTime utils.JsonTime `json:"disconnectedTime"`
}

func (r *Session) TableName() string {
	return "sessions"
}

type SessionVo struct {
	ID               string         `json:"id"`
	Protocol         string         `json:"protocol"`
	IP               string         `json:"ip"`
	Port             int            `json:"port"`
	Username         string         `json:"username"`
	ConnectionId     string         `json:"connectionId"`
	AssetId          string         `json:"assetId"`
	Creator          string         `json:"creator"`
	ClientIP         string         `json:"clientIp"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	Status           string         `json:"status"`
	ConnectedTime    utils.JsonTime `json:"connectedTime"`
	DisconnectedTime utils.JsonTime `json:"disconnectedTime"`
	AssetName        string         `json:"assetName"`
	CreatorName      string         `json:"creatorName"`
}

func FindPageSession(pageIndex, pageSize int, status, userId, clientIp, assetId, protocol string) (results []SessionVo, total int64, err error) {

	db := config.DB
	var params []interface{}

	params = append(params, status)

	itemSql := "SELECT s.id, s.protocol, s.connection_id, s.asset_id, s.creator, s.client_ip, s.width, s.height, s.ip, s.port, s.username, s.status, s.connected_time, s.disconnected_time, a.name AS asset_name, u.nickname AS creator_name FROM sessions s LEFT JOIN assets a ON s.asset_id = a.id LEFT JOIN users u ON s.creator = u.id WHERE s.STATUS = ? "
	countSql := "select count(*) from sessions as s where s.status = ? "

	if len(userId) > 0 {
		itemSql += " and s.creator = ?"
		countSql += " and s.creator = ?"
		params = append(params, userId)
	}

	if len(clientIp) > 0 {
		itemSql += " and s.client_ip like ?"
		countSql += " and s.client_ip like ?"
		params = append(params, "%"+clientIp+"%")
	}

	if len(assetId) > 0 {
		itemSql += " and s.asset_id = ?"
		countSql += " and s.asset_id = ?"
		params = append(params, assetId)
	}

	if len(protocol) > 0 {
		itemSql += " and s.protocol = ?"
		countSql += " and s.protocol = ?"
		params = append(params, protocol)
	}

	params = append(params, (pageIndex-1)*pageSize, pageSize)
	itemSql += " order by s.connected_time desc LIMIT ?, ?"

	db.Raw(countSql, params...).Scan(&total)

	err = db.Raw(itemSql, params...).Scan(&results).Error

	if results == nil {
		results = make([]SessionVo, 0)
	}
	return
}

func FindSessionByStatus(status string) (o []Session, err error) {
	err = config.DB.Where("status = ?", status).Find(&o).Error
	return
}

func CreateNewSession(o *Session) (err error) {
	err = config.DB.Create(o).Error
	return
}

func FindSessionById(id string) (o Session, err error) {
	err = config.DB.Where("id = ?", id).First(&o).Error
	return
}

func FindSessionByConnectionId(connectionId string) (o Session, err error) {
	err = config.DB.Where("connection_id = ?", connectionId).First(&o).Error
	return
}

func UpdateSessionById(o *Session, id string) {
	o.ID = id
	config.DB.Updates(o)
}

func DeleteSessionById(id string) {
	config.DB.Where("id = ?", id).Delete(&Session{})
}

func DeleteSessionByStatus(status string) {
	config.DB.Where("status = ?", status).Delete(&Session{})
}

func CountOnlineSession() (total int64, err error) {
	err = config.DB.Where("status = ?", Connected).Find(&Session{}).Count(&total).Error
	return
}
