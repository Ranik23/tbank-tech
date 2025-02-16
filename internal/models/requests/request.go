package requests

type AddLinkRequest struct {
	Link    string   `json:"link"`    // Ссылка
	Tags    []string `json:"tags"`    // Теги
	Filters []string `json:"filters"` // Фильтры
}

type RemoveLinkRequest struct {
	Link string `json:"link"` // Ссылка
}
