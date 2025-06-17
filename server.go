package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Адрес сервера
	serverAddress := "127.0.0.1:20080"

	// Подключение к серверу
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v\n", err)
	}
	defer conn.Close()

	fmt.Printf("Подключение к серверу %s установлено.\n", serverAddress)

	// Канал для завершения программы
	done := make(chan bool)

	// Горутина для получения данных с сервера
	go func() {
		for {
			// Буфер для получения данных
			response := make([]byte, 1024)
			n, err := conn.Read(response)
			if err != nil {
				log.Printf("Соединение с сервером разорвано: %v\n", err)
				done <- true
				return
			}

			// Вывод данных от сервера
			fmt.Printf("\nОтвет сервера: %s\n", string(response[:n]))
			fmt.Print(">> ") // Чтобы было удобно продолжать ввод
		}
	}()

	// Горутина для отправки пользовательских данных на сервер
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(">> ")
			// Считываем строку от пользователя
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Ошибка ввода: %v\n", err)
				done <- true
				return
			}

			// Отправляем данные на сервер
			_, err = conn.Write([]byte(input))
			if err != nil {
				log.Printf("Ошибка отправки данных: %v\n", err)
				done <- true
				return
			}
		}
	}()

	// Ожидание завершения (по сигналу из канала `done`)
	<-done
	fmt.Println("Клиент завершил работу.")
}
