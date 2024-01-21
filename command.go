package main

func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		err = SendMenu(chatId)
		break
	}

	return err
}
