package main

//import "database/sql"

//
//func getNumberOfUsers() (int64, error) {
//
//	var count int64
//
//	//Подключаемся к БД
//	db, err := sql.Open("postgres", dbInfo)
//	if err != nil {
//		return 0, err
//	}
//	defer db.Close()
//
//	//Отправляем запрос в БД для подсчета числа уникальных пользователей
//	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
//	err = row.Scan(&count)
//	if err != nil {
//		return 0, err
//	}
//
//	return count, nil
//}
