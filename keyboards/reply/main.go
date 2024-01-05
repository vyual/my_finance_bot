package reply

import tele "gopkg.in/telebot.v3"

var (
	mainMenu = &tele.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.

	BtnShowMovementTypes = mainMenu.Text("➕ Добавить/➖ Убрать")
	BtnShowMovements     = mainMenu.Text("📈 Показать передвижения")
	BtnShowBalance       = mainMenu.Text("💰 Показать текущий баланс")
)

func BuildMainMenu() *tele.ReplyMarkup {
	mainMenu.Reply(
		mainMenu.Row(BtnShowMovementTypes),
		mainMenu.Row(BtnShowMovements),
		mainMenu.Row(BtnShowBalance),
	)
	return mainMenu
}
