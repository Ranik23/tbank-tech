package responses

type LinkResponse struct {
	ID      int64    `json:"id"`      // ID ссылки
	URL     string   `json:"url"`     // Ссылка
	Tags    []string `json:"tags"`    // Теги
	Filters []string `json:"filters"` // Фильтры
}

type ListLinksResponse struct {
	Links []LinkResponse `json:"links"` // Список ссылок
	Size  int            `json:"size"`  // Количество ссылок
}

type ApiErrorResponse struct {
	Description     string   `json:"description"`     // Описание ошибки
	Code           string   `json:"code"`           // Код ошибки
	ExceptionName  string   `json:"exceptionName"`  // Название исключения
	ExceptionMessage string `json:"exceptionMessage"` // Сообщение исключения
	Stacktrace     []string `json:"stacktrace"`     // Стектрейс
}
