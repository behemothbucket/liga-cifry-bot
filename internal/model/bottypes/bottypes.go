package bottypes

// Типы для описания состава кнопок телеграм сообщения.
// Кнопка сообщения.
type TgInlineButton struct {
	DisplayName string
	Value       string
}

// Строка с кнопками сообщения.
type TgRowButtons []TgInlineButton