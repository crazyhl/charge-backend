package dto

import "charge/models"

// CategoryGroup 账户模型
type CategoryGroup struct {
	Type       uint8             `json:"type"`
	Categories []models.Category `json:"categories"`
}
