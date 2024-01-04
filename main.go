package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	botToken, _ := os.LookupEnv("TOKEN")

	pref := tele.Settings{
		Token: botToken,

		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	db := InitDb()
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		// Universal markup builders.
		menu     = &tele.ReplyMarkup{ResizeKeyboard: true}
		selector = &tele.ReplyMarkup{}

		// Reply buttons.
		btnShowMovementTypes = menu.Text("➕ Добавить/➖ Убрать")

		btnShowMovements = menu.Text("📈 Показать передвижения")

		btnShowBalance = menu.Text("💰 Показать текущий баланс")

		// Inline buttons.
		btnPrev = selector.Data("⬅", "prev")
		btnNext = selector.Data("➡", "next")
	)

	menu.Reply(
		menu.Row(btnShowMovementTypes),
		menu.Row(btnShowMovements),
		menu.Row(btnShowBalance),
	)
	selector.Inline(
		selector.Row(btnPrev, btnNext),
	)

	b.Handle("/start", func(c tele.Context) error {
		result := CreateUser(db, c.Message().Sender.ID, 0, c.Message().Sender.Username) // передаем указатель на данные в Create
		if result.Error == nil {
			log.Println(result.Error)
			return c.Send(
				"Добрый день! Это бот для учета финансов. Подскажите, что вы хотите сделать?",
				menu)
		}
		answer := fmt.Sprintf("С возвращением! Что вам нужно?")
		return c.Send(answer, menu)
	})

	// TODO добавить FSM
	b.Handle(btnShowMovementTypes.Text, func(c tele.Context) error {
		return c.Send(
			"Чтобы добавить доход, напишите:\n"+
				"ДОХОД <категория> <количество>\n\n"+
				"Чтобы добавить расход, напишите:\n"+
				"РАСХОД <категория> <количество>\n\n"+
				"Чтобы добавить накопление, напишите:\n"+
				"КОПИЛКА <категория> <количество>\n\n"+
				"Чтобы убрать накопление, напишите:\n"+
				"-КОПИЛКА <категория> <количество>\n\n",
			menu,
		)
	})

	// TODO добавить FSM
	// TODO клавиатура выбора промежутка времени
	// TODO выбор по типу передвижений
	// TODO выбор по категории
	b.Handle(btnShowMovements.Text, func(c tele.Context) error {
		return c.Send(
			"Ваши денежные передвижения за:\n\n" +
				"- Сегодня\n" +
				"- Месяц\n" +
				"- Год\n",
		)
	})

	// TODO добавить FSM
	// TODO клавиатура выбора промежутка времени
	b.Handle(btnShowBalance.Text, func(c tele.Context) error {
		user := GetUserByTelegramId(db, c.Sender().ID)
		answer := fmt.Sprintf("Ваш баланс: %d руб.\n\n"+
			"Выйти в меню: /start", user.Balance)
		return c.Send(
			answer,
		)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		// All the text messages that weren't
		// captured by existing handlers.
		var (
			user    = GetUserByTelegramId(db, c.Message().Sender.ID)
			textArr = strings.Split(c.Text(), " ")
		)
		moneyMovementTypeString := textArr[0]
		switch moneyMovementTypeString {
		case "ДОХОД":
			quantity, err := strconv.Atoi(textArr[2])
			if err != nil {
				log.Println(err)
				return c.Send("Вы ошиблись в параметре <количество>! Отправьте, пожалуйста, число.")
			}

			moneyMovementType := GetMoneyMovementTypeByName(db, moneyMovementTypeString)

			var category Category
			categoryName := textArr[1]
			result := db.First(&category, "name = ?", categoryName)
			if result.Error != nil {
				//return c.Send("Данной категории нет. Создать?")
				result, category = CreateCategory(db, user, textArr[1], moneyMovementType)
			}

			CreateMoneyMovement(db, user, int64(quantity), category)
			UpdateUserBalance(db, user, quantity)
			return c.Send("ДОХОД успешно добавлен!")
		case "РАСХОД":
			quantity, err := strconv.Atoi(textArr[2])
			if err != nil {
				log.Println(err)
				return c.Send("Вы ошиблись в параметре <количество>! Отправьте, пожалуйста, число.")
			}

			moneyMovementType := GetMoneyMovementTypeByName(db, moneyMovementTypeString)

			var category Category
			categoryName := textArr[1]
			result := db.First(&category, "name = ?", categoryName)
			if result.Error != nil {
				//return c.Send("Данной категории нет. Создать?")
				result, category = CreateCategory(db, user, textArr[1], moneyMovementType)
			}

			CreateMoneyMovement(db, user, int64(quantity), category)
			UpdateUserBalance(db, user, -quantity)
			return c.Send("РАСХОД успешно добавлен!")
		case "КОПИЛКА":
			quantity, err := strconv.Atoi(textArr[2])
			if err != nil {
				log.Println(err)
				return c.Send("Вы ошиблись в параметре <количество>! Отправьте, пожалуйста, число.")
			}

			moneyMovementType := GetMoneyMovementTypeByName(db, moneyMovementTypeString)

			var category Category
			categoryName := textArr[1]
			result := db.First(&category, "name = ?", categoryName)
			if result.Error != nil {
				//return c.Send("Данной категории нет. Создать?")
				result, category = CreateCategory(db, user, textArr[1], moneyMovementType)
			}

			CreateMoneyMovement(db, user, int64(quantity), category)
			UpdateUserBalance(db, user, quantity)
			return c.Send("КОПИЛКА успешно добавлен!")
		case "-КОПИЛКА":
			quantity, err := strconv.Atoi(textArr[2])
			if err != nil {
				log.Println(err)
				return c.Send("Вы ошиблись в параметре <количество>! Отправьте, пожалуйста, число.")
			}

			moneyMovementType := GetMoneyMovementTypeByName(db, moneyMovementTypeString)

			var category Category
			categoryName := textArr[1]
			result := db.First(&category, "name = ?", categoryName)
			if result.Error != nil {
				//return c.Send("Данной категории нет. Создать?")
				result, category = CreateCategory(db, user, textArr[1], moneyMovementType)
			}

			CreateMoneyMovement(db, user, int64(quantity), category)
			UpdateUserBalance(db, user, -quantity)
			return c.Send("-КОПИЛКА успешно добавлен!")
		default:

			return c.Send("Я вас не очень понял! обратитесь в меню /help")
		}
	})

	// On inline button pressed (callback)
	b.Handle(&btnPrev, func(c tele.Context) error {
		return c.Respond()
	})

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
	log.Printf("BOT STARTED\n")
	b.Start()
}
