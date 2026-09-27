package main

import (
	"bufio"
	"edtech/internal/infrastructure/config"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	"edtech/internal/infrastructure/migrator"
	"fmt"
	"log/slog"
	"os"
)

const MigDir = "./migrators"
const configPath = "./.env"

func main() {

	log := logger.Init()

	cfg, err := config.Load(configPath)
	if err != nil {

	}

	log.Info("Конфиг иннициализирован")

	postgres, err := db.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	log.Info("Подключение к базе данных выполнено")

	defer postgres.Pool.Close()
	migrator, err := migrator.NewMigrator(postgres.Pool, MigDir)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			if err := migrator.Up(); err != nil {
				log.Error("migration UP failed", sl.Error(err))
				os.Exit(1)
			}
			log.Info("UP Migration applied")
			return
		case "down":
			if err := migrator.Down(); err != nil {
				log.Error("migration DOWN failed", sl.Error(err))
				os.Exit(1)
			}
			log.Info("DOWN Migration applied")
			return
		}
	}

	Migration(log, migrator)

}

func Migration(log *slog.Logger, migrator *migrator.Migrator) {

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
