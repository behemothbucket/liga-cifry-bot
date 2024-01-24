package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot() *Bot {
	var err error

	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Panic(err)
	}

	return &Bot{bot: bot}
}

func getToken() string {
	token, exists := os.LookupEnv("TOKEN")

	if !exists {
		log.Print("Токен не обнаружен")
	}

	return token
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.handleMessage(update)
		break

	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery)
		break
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	if update.Message.From == nil {
		return
	}

	user := update.Message.From
	text := update.Message.Text

	logMessage(user, text)

	var err error

	if isValidMessageText(text) {
		switch {
		case update.Message.IsCommand():
			b.handleCommand(update)
		case searchMode:
			err = b.sendAcceptMessage(update.Message.Chat.ID)
		default:
			err = b.sendMainMenu(update.Message.Chat.ID)
		}
	} else {
		err = b.sendMediaErrorMessage(update.Message.Chat.ID)
	}

	if err != nil {
		log.Println("Ошибка:", err.Error())
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	command := update.Message.Text
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	botName := b.bot.Self.UserName

	switch command {
	case "/start", "/start@" + botName:
		b.sendMainMenu(chatID)
		break
	case "/checkmember", "/checkmember@" + botName:
		member, _ := b.bot.GetChatMember(
			tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					ChatID: chatID,
					UserID: userID,
				},
			})

		b.sendMessage(chatID, fmt.Sprintf("Вы являетесь участником группы: %t\nИмя пользователя: %s", member.IsMember, member.User.UserName))
		break
	case "/user", "/user@" + botName:
		showSearchResultsMode = true
		text := `
🧑‍💼 <b>ФИО</b>
Рахматуллина Елена Сергеевна

🏛 <b>Организация</b>
АНО ВО "Университет Иннополис"

💼 <b>Должность</b>
Доцент, руководитель направления по работе с партнерами

📝 <b>Компетенции</b>
Доцент, к.э.н., Консультации по направлениям: анализ финансово-хозяйственной деятельности предприятий, стратегическое управление , управление персоналом, бизнес-аналитика, редакторско-издательская деятельность научных журналов

🤝 <b>Направления сотрудничества</b>
Организация и проведение лекций, вебинаров в качестве организатора и спикера. Консультации по вопросам принципов работы и приобретения программы ИСУ РИД (информационная система управления результатами интеллектуальной деятельности). Аналитические исследования.

📱 <b>Контакты</b> @woman_21`
		b.sendMarkupMessage(chatID, text)
	case "/university", "/university@" + botName:
		showSearchResultsMode = true
		text := `
🏛 <b>Организация</b>
АНО ВО "Университет Иннополис"

🏢 <b>Структурное подразделение</b>
Межотраслевой центр трансфера технологий

🤝 <b>Направления сотрудничества</b>
Управление результатами интеллектуальной деятельности на всем цикле от планирования и создания РИД до процессов коммерциализации и трансфера технологий

🚀 <b>«Приоритет-2030»</b>
Нет

🎓 <b>Кампус мирового уровня</b>
Нет

🌐 <b>Членство в консорциуме</b>
Да

⚛️ <b>Разработки отвечественного ПО</b>
Программа ИСУ РИД (Комплексное автоматизированное решение, предназначенное для автоматизации процесса сбора, учета и обработки данных о результатах интеллектуальной деятельности организаций). Запись в реестре от 01.03.2023 №16846)

🔬 <b>Лабораторные площадки и НОЦ</b>
Нет`
		b.sendMessage(chatID, text)

		textSkills := `📝 <b>Компетенции</b>
- Регистрация результатов интеллектуальной деятельности
- Подготовка материалов для заявки на изобретение и ведение делопроизводства
- Подготовка материалов для заявки на полезную модель и ведение делопроизводства
- Подготовка материалов для заявки на промышленный образец и ведение делопроизводства
- Подготовка материалов для заявки на товарный знак и ведение делопроизводства
- Подготовка материалов для заявки на товарный знак и ведение делопроизводства «под ключ» (подача заявки и делопроизводство, проверка юридической чистоты товарного знака – оформление необходимых договоров передачи прав с авторами логотипа, запросы писем-согласий)
- Подготовка материалов и регистрация программы для ЭВМ
- Подготовка материалов и регистрация базы данных
- Проверка и подготовка к подаче комплекта документации для внесения IT-решения в Единый реестр российских программ для электронных вычислительных машин и баз данных Патентные и маркетинговые исследования
- Проведение патентно-информационного поиска на определение аналогов в отношении одного технического решения
- Проведение патентных исследований по ГОСТ на определение уровня техники
- Проведение патентных исследований по ГОСТ на патентную чистоту
- Подготовка патентного ландшафта
- Проведение маркетинговых исследований Техническая документация
- Подготовка и оформление программы и методики испытаний по ГОСТ
- Подготовка и оформление руководства пользователя/оператора по ГОСТ
- Подготовка и оформление руководства администратора/системного программиста по ГОСТ
- Подготовка и оформление описания жизненного цикла программного обеспечения Международные заявки
- Подготовка материалов и подача международной заявки на изобретение
- Подготовка материалов и подача международной заявки на товарный знак
- Подготовка материалов для евразийской заявки на изобретение и ведение делопроизводства
- Подготовка материалов для евразийской заявки на промышленный образец и ведение делопроизводства
- Подготовка материалов международной заявки на промышленный образец и ведение делопроизводства
- Подготовка материалов и подача международной заявки на товарный знак Аудит и консалтинг
- Подготовка отчета об оценке уровня готовности технологий по ГОСТ
- Консультации по вопросам охраны, управления и коммерциализации РИД, трансфера технологий Аудит системы управления интеллектуальной собственностью Технологический аудит Обучение по вопросам охраны, управления и коммерциализации РИД, трансфера технологий
`
		b.sendMarkupMessage(chatID, textSkills)
	}
}

func (b *Bot) handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := mainMenuMarkup
	message := query.Message

	if query.Data == searchUserButton {
		text = searchMenuDescription
		markup = getUserSearchMenuMarkup()
		currentSearchScreen = "user"
	} else if query.Data == searchUniversityButton {
		text = searchMenuDescription
		markup = getUniversitySearchMenuMarkup()
		currentSearchScreen = "university"
	} else if query.Data == backButton || query.Data == backToMainMenuButton {
		text = mainMenuDescription
		markup = mainMenuMarkup
		for k := range searchCriterias {
			delete(searchCriterias, k)
		}
	} else if query.Data == applyButton {
		if len(searchCriterias) == 0 {
			text = "️❗️Пожалуйста, выберите хотя-бы один критерий поиска"
			markup = getUserSearchMenuMarkup()
		} else {
			text = getCriteria()
			searchMode = true
			cancelMenuMarkup = getCancelMenuMarkup()
			markup = cancelMenuMarkup
		}
	} else if query.Data == cancelButton {
		removeAllSearchCriterias()
		resetCriteriaButtons()
		searchMode = false
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		b.bot.Send(callbackCfg)
		b.sendMainMenu(message.Chat.ID)
		return
	} else if criteriaButtonIsClicked(query.Data) {
		toggleCriteriaButton(query.Data)
		text = searchMenuDescription
		if currentSearchScreen == "user" {
			markup = getUserSearchMenuMarkup()
		} else {
			markup = getUniversitySearchMenuMarkup()
		}
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.bot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.bot.Send(msg)
}
