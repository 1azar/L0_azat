package main

import (
	"L0_azat/tests"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	url                      = "nats://localhost:4222"
	subject                  = "L0-subject"
	display_sending_messages = true
)

func main() {
	fmt.Println("CLI для генерации данных заказа и отправки сервису")
	fmt.Println("Параметры:")
	fmt.Println("\turl: " + url)
	fmt.Println("\tsubject: " + subject)
	fmt.Println("\tdisplay_sending_messages: " + strconv.FormatBool(display_sending_messages))
	fmt.Println("Доступные команды:")
	fmt.Println("\tspam=<интервал спама в сек.> - генерирует случайные заказы через указанный интервал времени:")
	fmt.Println("\tsend=<количество заказов> - отправляет указанное количество случайных заказов")

	scanner := bufio.NewScanner(os.Stdin)

MAIN_LOOP:
	for {
		fmt.Print("Введите команду (или 'выход' для завершения): ")
		scanner.Scan()
		command := strings.ToLower(scanner.Text())

		switch {
		case strings.HasPrefix(command, "send="):
			valStr := strings.TrimPrefix(command, "send=")
			valInt, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("неверное знчение аргумента")
				continue
			}
			for i := 0; i < valInt; i++ {
				msg := tests.GenerateMsg()
				tests.SendMsg(url, subject, msg)
				if display_sending_messages {
					msgByte, err := json.Marshal(msg)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println("Message sent: \n", string(msgByte))
				}
			}

		case strings.HasPrefix(command, "spam="):
			valStr := strings.TrimPrefix(command, "spam=")
			_, err := strconv.Atoi(valStr)
			if err != nil {
				fmt.Println("неверное знчение аргумента")
				continue
			}
			fmt.Println("Implement me!")

		case strings.HasPrefix(command, "выход"):
			break MAIN_LOOP

		default:
			fmt.Println("Неизвестная команда. Попробуйте еще раз.")
		}
	}
}

//func getCurrentTime() string {
//	return "Реальная реализация для получения времени здесь"
//}
