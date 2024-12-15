package main

import (
	"fmt"

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

	//Объявление очереди
	q, err := ch.QueueDeclare(
		"TestQueue",
		false, // durable: очередь не будет сохраняться при перезапуске сервера
		false, // autoDelete: очередь не будет удалена, когда все подписчики отключатся
		false, // exclusive: очередь не будет эксклюзивной для данного соединения
		false, // noWait: не ждать подтверждения создания очереди
		nil,   // аргументы (можно передать дополнительные параметры)
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(q)

	//Публикация сообщения в очередь
	err = ch.Publish(
		"",
		"TestQueue", // имя очереди, в которую отправляется сообщение
		false,       // mandatory: если `true`, и сообщение не может быть доставлено, оно будет возвращено отправителю
		false,       // immediate: если `true`, сообщение будет отправлено немедленно
		amqp.Publishing{
			ContentType: "text/plain", // тип содержимого сообщения
			Headers: amqp.Table{
				"type": "hello",
			},
			Body: []byte("Это сообщение"), // само сообщение
		},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	} // Обработка ошибки при публикации
	fmt.Println("Сообщение успешно отпрвленно в очередь")
}
