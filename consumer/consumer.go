package main

import (
	"fmt"
	"main/redis"

	"github.com/streadway/amqp"
)

func main() {
	//Подключаемся с Реббиту
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	} //Проверка на ошибку при подключении
	defer conn.Close() //Гарант, что файл будет закрыт, когда main завершит выполнение, независимо от того, будет ли это нормальный выход или завершение из-за ошибки
	fmt.Println("Удалось подключиться к RabbitMq")

	//Делаем канал для сообщений
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	} //Проверка на ошибку при создании канала

	defer ch.Close() //гарант закрытия

	//Подключаем наш канал к Реббит с помощью этих параметров
	msgs, err := ch.Consume(
		"TestQueue", // имя очереди, из которой мы хотим получать сообщения
		"",          // consumerTag: уникальный идентификатор для данного потребителя (можно оставить пустым)
		true,        // autoAck: если true, сообщения будут автоматически подтверждаться после их получения
		false,       // exclusive: если true, очередь будет доступна только для этого потребителя
		false,       // noLocal: если true, потребитель не будет получать сообщения, отправленные им самим
		false,       // noWait: если true, метод не будет ждать подтверждения
		nil,         // аргументы: дополнительные параметры
	)

	//Делаем канал
	forever := make(chan bool)
	//Запускаем Гоурутину с выводом всех сообщений из Нашего канала msgs
	go func() {
		for d := range msgs {
			if msgType, ok := d.Headers["type"].(string); ok && msgType == "hello" {
				fmt.Println("Полученно сообщение с типом: hello")
				redis.Increment("test")
			}
			fmt.Printf("Полученное сообщение: %s\n", d.Body)
		}
	}()

	//Пишем в консоль всю инфу
	fmt.Println("Успешное подключение к rabbitMq")
	fmt.Println(" [*] - Ожидание сообщений")
	//Блокируем основную гоурутину
	<-forever

}
