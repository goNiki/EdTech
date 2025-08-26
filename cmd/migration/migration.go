package main

import (
	"bufio"
	"edtech/internal/infrastructure/config"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	"fmt"
	"log/slog"
	"os"
)

const MigDir = "./migrators"

func main() {

	log := logger.Init()

	cfg := config.InitConfig()

	log.Info("Конфиг иннициализирован")
	fmt.Println(cfg)
	fmt.Println(cfg.DBConfig)

	postgres, err := db.New(cfg)
	if err != nil {
		log.Error("Ошибка в иннициализации База данных: %w", sl.Error(err))
		return
	}

	log.Info("Подключение к базе данных выполнено")

	defer postgres.Close()
	migrator := db.NewMigrator(postgres.Pool, MigDir)
	Migration(log, migrator)

}

func Migration(log *slog.Logger, migrator *db.Migrator) {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		var choice int64

		fmt.Println("Выберите действие: \n1 - Создать файл для миграции \n2 - Проверить статус миграций \n3 - Миграций UP \n4- Миграций UpTo \n5 - Миграция Down \n6 - Миграция DownTo \n7 - выйти")
		if scanner.Scan(); scanner.Err() != nil {
			fmt.Println("Ошибка ввода")
			continue
		}
		if _, err := fmt.Sscanf(scanner.Text(), "%d", &choice); err != nil || choice < 1 || choice > 7 {
			fmt.Println("Ошибка ввода. Не соответствует условиям")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("Введите название файла ")
			scanner.Scan()
			name := scanner.Text()
			if err := migrator.Create(name, "sql"); err != nil {
				log.Error("Ошибка: ", sl.Error(err))
			}
			log.Info("migration created successfully")

		case 2:
			if err := migrator.Status(); err != nil {
				log.Error("Ошибка", sl.Error(err))
			}
			log.Info("Status migrations apply")
		case 3:
			if err := migrator.Up(); err != nil {
				log.Error("Ошибка", sl.Error(err))
			} else {
				log.Info("UP Migration apply")
			}
		case 4:
			fmt.Println("Введите версию миграции:")
			scanner.Scan()
			var name int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d", &name); err != nil {
				fmt.Println("Ошибка при вводе версии миграции")
				continue
			}
			if err := migrator.UpTo(name); err != nil {
				log.Error("", sl.Error(err))
			}
			log.Info("UpTo migration apply")
		case 5:
			if err := migrator.Down(); err != nil {
				log.Error("", sl.Error(err))
			}
			log.Info("DOWN migration apply")
		case 6:
			fmt.Println("Введите версию миграции:")
			scanner.Scan()
			var name int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d", &name); err != nil {
				fmt.Println("Ошибка при вводе версии миграции")
				continue
			}

			if err := migrator.DownTo(name); err != nil {
				log.Error("", sl.Error(err))
			}
			log.Info("DownTo migration apply")
		case 7:
			return
		}
	}
}
