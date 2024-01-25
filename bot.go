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

	chatID := update.Message.Chat.ID
	user := update.Message.From
	message := update.Message

	logMessage(user, message.Text, chatID)

	var err error

	if b.isValidMessageText(update) {
		switch {
		case update.Message.IsCommand():
			b.handleCommand(update)
		case searchMode:
			err = b.sendAcceptMessage(chatID)
		default:
			err = b.sendMainMenu(chatID)
		}
	} else if len(message.NewChatMembers) != 0 {
		userName := message.NewChatMembers[0].UserName
		firstName := message.NewChatMembers[0].FirstName
		lastName := message.NewChatMembers[0].LastName
		userID := message.NewChatMembers[0].ID
		isBot := message.NewChatMembers[0].IsBot
		AddUser(userID, userName, firstName, lastName, isBot)
	} else if update.Message.LeftChatMember.UserName != "" {
		userID := update.Message.LeftChatMember.ID
		DeleteUser(userID)
	} else {
		err = b.sendMediaErrorMessage(update.Message.Chat.ID)
	}

	if err != nil {
		log.Println("Ошибка:", err.Error())
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	command := update.Message.Text
	//chatType := update.Message.Chat.Type
	chatID := update.Message.Chat.ID
	//userID := update.Message.From.ID
	botName := fmt.Sprintf("@%s", b.bot.Self.UserName)

	//if chatType == "group" || chatType == "private" || chatType == "supergroup" || chatType == "channel" {
	//	log.Printf("Обнаружен чат с типом %s, блокирую комманду", chatType)
	//	b.SendMessage(chatID, "Вы не являетесь администратором чтобы использовать данную комманду")
	//	return
	//}

	switch command {
	case "/start", "/start" + botName:
		b.sendMainMenu(chatID)
		break
	case "/checkmember", "/checkmember" + botName:
		member, _ := b.bot.GetChatMember(
			tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					ChatID: -1001992755263,
					UserID: 155401792,
				},
			})
		b.SendMessage(chatID, fmt.Sprintf("Является участником группы: %s\nВышел из группы: %t\nИмя пользователя: %s\nchatID: %d", member.Status, member.HasLeft(), member.User.UserName, chatID))
		break
	case "/user", "/user" + botName:
		showSearchResultsMode = true
		card1 := `
🧑‍💼 <b>ФИО</b>
Абдрахманова Марина Васильевна

🏛 <b>Организация</b>
АНО ВО "Университет Иннополис"

💼 <b>Должность</b>
Руководитель межотраслевого центра трансфера технологий

📝 <b>Компетенции</b>
🔹Управление интеллектуальной собственностью
🔹Защита и коммерциализация IT решений
🔹Стратеческое управление нематериальными активами
🔹Регистрация в реестре отечественного ПО
🔹Создание цифровых решений в области управления интеллектуальной собственностью и баз знаний
🔹Управление инновационной деятельностью
🔹Нормативно-методическое сопровождение управеления инновационной деятельность
🔹Оценка уровня готовности проектов и технологий
🔹Подготовка маркетинго-аналитических исследований сложных технологических проектов
🔹Создание и организация подготовки кастомизированных образовательных программ в рамках экспертных компетенций

🤝 <b>Направления сотрудничества</b>
🔹обмен компетенциями в области технического, правового, экономического аудита технологических проектов;
🔹совместно формирование и проведение мероприятий и образовательных программ в области трансфера технологий;
🔹привлечение спикеров на мероприятия в области трансфера технологий;
🔹совместная работа и создание технологических проектов и цифровых решений, в том числе выход на совместное грантовое или частное финансирование;
🔹усиление компетенций в области привлечения частного финансирования проектов (привлечение инвестирования);
🔹усиление компетенций в области создания системы принятия решений по инвестированию во внутренние инициативные проекты органзиации;
🔹кандидаты, усиление и расширение штата патентных поверенных и IP юристов

📧 <b>Контакты</b> m.abdrakhmanova@innopolis.ru`
		card2 := `
🧑‍💼 <b>ФИО</b>
Рахматуллина Елена Сергеевна
		
🏛 <b>Организация</b>
АНО ВО "Университет Иннополис"
		
💼 <b>Должность</b>
Доцент, руководитель направления по работе с партнерами
		
📝 <b>Компетенции</b>
Консультации по направлениям: 
🔹анализ финансово-хозяйственной деятельности предприятий
🔹стратегическое управление
🔹управление персоналом
🔹бизнес-аналитика
🔹редакторско-издательская деятельность научных журналов
		
🤝 <b>Направления сотрудничества</b>
🔹Организация и проведение лекций, вебинаров в качестве организатора и спикера
🔹Консультации по вопросам принципов работы и приобретения программы ИСУ РИД (информационная система управления результатами интеллектуальной деятельности)
🔹Аналитические исследования.
		
📱 <b>Контакты</b> @woman_21`
		b.sendMarkupMessage(chatID, card1)
		showSearchResultsMode = true
		b.sendMarkupMessage(chatID, card2)
	case "/university", "/university@" + botName:
		showSearchResultsMode = true
		text := `
🏛 <b>Организация</b>
АНО ВО "Университет Иннополис"

🏢 <b>Структурное подразделение</b>
Межотраслевой центр трансфера технологий

📍 <b>Город</b>
г. Иннополис

🤝 <b>Направления сотрудничества</b>
Управление результатами интеллектуальной деятельности на всем цикле от планирования и создания РИД до процессов коммерциализации и трансфера технологий

🚀 <b>«Приоритет-2030»</b>
Нет

🌐 <b>Членство в консорциуме</b>
Да

⚛️ <b>Разработки отвечественного ПО</b>
Программа ИСУ РИД (Комплексное автоматизированное решение, предназначенное для автоматизации процесса сбора, учета и обработки данных о результатах интеллектуальной деятельности организаций). Запись в реестре от 01.03.2023 №16846)

🔬 <b>Лабораторные площадки и НОЦ</b>
Нет`
		b.SendMessage(chatID, text)

		textSkills := `📝 <b>Компетенции</b>
🔹 Регистрация результатов интеллектуальной деятельности
🔹 Подготовка материалов для заявки на изобретение и ведение делопроизводства
🔹 Подготовка материалов для заявки на полезную модель и ведение делопроизводства
🔹 Подготовка материалов для заявки на промышленный образец и ведение делопроизводства
🔹 Подготовка материалов для заявки на товарный знак и ведение делопроизводства
🔹 Подготовка материалов для заявки на товарный знак и ведение делопроизводства «под ключ» (подача заявки и делопроизводство, проверка юридической чистоты товарного знака – оформление необходимых договоров передачи прав с авторами логотипа, запросы писем-согласий)
🔹 Подготовка материалов и регистрация программы для ЭВМ
🔹 Подготовка материалов и регистрация базы данных
🔹 Проверка и подготовка к подаче комплекта документации для внесения IT-решения в Единый реестр российских программ для электронных вычислительных машин и баз данных Патентные и маркетинговые исследования
🔹 Проведение патентно-информационного поиска на определение аналогов в отношении одного технического решения
🔹 Проведение патентных исследований по ГОСТ на определение уровня техники
🔹 Проведение патентных исследований по ГОСТ на патентную чистоту
🔹 Подготовка патентного ландшафта
🔹 Проведение маркетинговых исследований Техническая документация
🔹 Подготовка и оформление программы и методики испытаний по ГОСТ
🔹 Подготовка и оформление руководства пользователя/оператора по ГОСТ
🔹 Подготовка и оформление руководства администратора/системного программиста по ГОСТ
🔹 Подготовка и оформление описания жизненного цикла программного обеспечения Международные заявки
🔹 Подготовка материалов и подача международной заявки на изобретение
🔹 Подготовка материалов и подача международной заявки на товарный знак
🔹 Подготовка материалов для евразийской заявки на изобретение и ведение делопроизводства
🔹 Подготовка материалов для евразийской заявки на промышленный образец и ведение делопроизводства
🔹 Подготовка материалов международной заявки на промышленный образец и ведение делопроизводства
🔹 Подготовка материалов и подача международной заявки на товарный знак Аудит и консалтинг
🔹 Подготовка отчета об оценке уровня готовности технологий по ГОСТ
🔹 Консультации по вопросам охраны, управления и коммерциализации РИД, трансфера технологий Аудит системы управления интеллектуальной собственностью
🔹Технологический аудит Обучение по вопросам охраны, управления и коммерциализации РИД, трансфера технологий
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
