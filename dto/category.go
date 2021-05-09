package dto

// CategoryGroup 账户模型
type CategoryGroup map[uint8][]Category

type Category struct {
	ID   uint   `json:"id"`
	Type uint8  `json:"type,omitempty"`
	Name string `json:"name"`
	Sort uint8  `json:"sort,omitempty"`
}
