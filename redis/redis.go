package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client //Создаю переменную для информации о сервере Redis

func init() {
	// Иициализирую клиент Redis.
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес
		Password: "",               // Пароль (Из-за того что я использую докер, пароля нет)
		DB:       0,                // Номер БД
	})
}

func Main() {
	ping, err := client.Ping(context.Background()).Result() //Пингуем сервер, чтобы убедиться что он запущен
	if err != nil {                                         //Проверяем ошибку - в случае чего выводим ее
		fmt.Println(err.Error())
		return
	}
	if ping == "PONG" {
		fmt.Println("Сервер Redis подключён") //Сообщение, что мы создали подключение с Redis
	}

	err = client.Set(context.Background(), "received_hello", 0, 0).Err()
	if err != nil {
		fmt.Printf("Не удалось задать значение received_hello. Вот ошибка: %s", err.Error())
		return
	} // В случае неудачи - выводим ошибку

	err = client.Set(context.Background(), "sent_hello", 0, 0).Err()
	if err != nil {
		fmt.Printf("Не удалось задать значение sent_hello. Вот ошибка: %s", err.Error())
		return
	} // В случае неудачи - выводим ошибку

}

func GetReceivedHelloCount() (int64, error) {
	val, err := client.Get(context.Background(), "received_hello").Result()
	if err != nil {
		fmt.Printf("Не удалось достать значение received_hello. Вот ошибка: %s\n", err.Error())
		return 0, err // Возвращаем 0 и ошибку при неудаче
	} // В случае неудачи - выводим ошибку

	// Преобразуем значение val к типу int64
	count, convErr := strconv.ParseInt(val, 10, 64)
	if convErr != nil {
		fmt.Printf("Не удалось преобразовать значение received_hello к int64. Вот ошибка: %s\n", convErr.Error())
		return 0, convErr // Возвращаем 0 и ошибку при неудаче преобразования
	}

	fmt.Printf("Полученное значение по ключу received_hello: %d\n", count)
	return count, nil // Возвращаем значение и nil, если все прошло успешно
}

// Инкремент строки в Redis	 Пример использования  ->   redis.Increment("Строка")
func Increment(key string) error {
	_, err := client.Incr(context.Background(), key).Result()
	fmt.Println("+1 к строке ", key)
	if err != nil {
		fmt.Printf("Не удалось увеличить строку %s. Вот ошибка: %s", key, err.Error())
		return err
	}
	return err
}

func GetSentHelloCount() (int64, error) {
	val, err := client.Get(context.Background(), "sent_hello").Result()
	if err != nil {
		fmt.Printf("Не удалось достать значение sent_hello. Вот ошибка: %s\n", err.Error())
		return 0, err // Возвращаем 0 и ошибку при неудаче
	} // В случае неудачи - выводим ошибку

	// Преобразуем значение val к типу int64
	count, convErr := strconv.ParseInt(val, 10, 64)
	if convErr != nil {
		fmt.Printf("Не удалось преобразовать значение sent_hello к int64. Вот ошибка: %s\n", convErr.Error())
		return 0, convErr // Возвращаем 0 и ошибку при неудаче преобразования
	}

	fmt.Printf("Полученное значение по ключу sent_hello: %d\n", count)
	return count, nil // Возвращаем значение и nil, если все прошло успешно
}
