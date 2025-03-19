package telegram

import (
	"gopkg.in/telebot.v3"
)


func (b *BotHandlers) HelpHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		msg := "<b>Доступные команды:</b>\n" +
			"🔹 <code>/help</code> – показать список команд\n" +
			"🔹 <code>/start</code> – начать работу с ботом\n" +
			"🔹 <code>/list</code> – показать список\n" +
			"🔹 <code>/track</code> – начать отслеживание\n" +
			"🔹 <code>/untrack</code> – прекратить отслеживание"
		
		return c.Send(msg, telebot.ModeHTML)
	}
}
