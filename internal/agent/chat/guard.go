package chat

import "strings"

const mutationRefused = "Пользователь не просил изменить или удалить ссылку."

// mutationAllowed разрешает delete_url и update_url только по прямой просьбе в последней реплике пользователя.
func mutationAllowed(history []Message, name string) bool {
	if name != "delete_url" && name != "update_url" {
		return true
	}

	text := strings.ToLower(lastUserText(history))
	for _, word := range []string{"удал", "измен", "обнов", "помен", "delete", "update", "change", "remove"} {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

// lastUserText возвращает текст последнего сообщения пользователя.
func lastUserText(history []Message) string {
	var text string
	for _, message := range history {
		if message.Role == "user" {
			text = message.Content
		}
	}
	return text
}

// visibleTools скрывает изменение и удаление, пока пользователь прямо об этом не попросил.
func visibleTools(history []Message, tools []Tool) []Tool {
	if mutationAllowed(history, "delete_url") {
		return tools
	}

	visible := make([]Tool, 0, len(tools))
	for _, tool := range tools {
		if tool.Name == "delete_url" || tool.Name == "update_url" {
			continue
		}
		visible = append(visible, tool)
	}
	return visible
}

// knownTool сообщает, есть ли инструмент в списке, который видит модель.
func knownTool(tools []Tool, name string) bool {
	for _, tool := range tools {
		if tool.Name == name {
			return true
		}
	}
	return false
}
