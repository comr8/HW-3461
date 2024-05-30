package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	arg1 := "input.txt"  // Дефолтное имя входного файла
	arg2 := "output.txt" // Дефолтное имя для вывода результатов

	if len(os.Args) < 3 {
		fmt.Printf("Необходимо указать 2 аргумента: имя входного файла и имя файла для вывода\nВыполнение программы продолжено со стандартными именами файлов: ./input.txt - для входных данных и ./output.txt - для вывода.\n")
	} else {
		// Значения аргументов
		arg1 = fmt.Sprintf("%s.txt", os.Args[1]) // имя входного файла
		arg2 = fmt.Sprintf("%s.txt", os.Args[2]) // имя файла для вывода результатов
	}

	// Регулярное выражение для проверки соответствия строки шаблону математического выражения
	re := regexp.MustCompile(`^[0-9]+[-+*/][0-9]+[=][?]$`)

	data := "" // Переменная для хранения конечной информации для записи

	// Открываем файл для чтения входных данных
	inputFile, err := os.OpenFile(arg1, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("ОШИБКА: вероятно файл с именем input.txt не найден.\n")
		return
	}
	defer inputFile.Close() // Отложенное закрытие файла после выполнения всех операций.

	// Создаем читатель для файла для буферизации чтения из файла
	fileReader := bufio.NewReader(inputFile)

	var result int // Для подсчета значений выражений вне цикла
	for {
		line, _, err := fileReader.ReadLine()
		if err == io.EOF {
			break
		} else {
			// Проверка соответствия строки основному шаблону
			if re.Match(line) {
				// Выражение для поиска всех последовательностей из цифр в строке
				nums := regexp.MustCompile(`[0-9]+`).FindAllStringSubmatch(string(line), -1)
				// Выражение для поиска арифметических операций в строке
				action := regexp.MustCompile(`[+-/*]`).FindStringSubmatch(string(line))

				number1, err := strconv.Atoi(nums[0][0])
				if err != nil {
					log.Fatalf("error converting %v to int: %v\n", nums[0][0], err)
				}
				number2, err := strconv.Atoi(nums[1][0])
				if err != nil {
					log.Fatalf("error converting %v to int: %v\n", nums[1][0], err)
				}

				switch action[0] {
				case "-":
					result = number1 - number2
				case "+":
					result = number1 + number2
				case "*":
					result = number1 * number2
				case "/":
					result = number1 / number2
				}
				// Выполняем математическую операцию в зависимости от найденной операции
				// Обновляем данные, добавляя рассчитанное значение в конечный вид строки
				data = data + strconv.Itoa(number1) + action[0] + strconv.Itoa(number2) + "=" + strconv.Itoa(result) + "\n"
			}
		}
	}
	if len(data) == 0 {
		fmt.Println("Результатов для записи не обнаружено\nПрограмма завершена")
		return
	}
	// Открываем файл output.txt для только для записи и создание файла, если он не существует
	// с полными правами
	outputFile, err := os.OpenFile(arg2, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close() // Отложенное закрытие файла после выполнения всех операций.

	// Создаем writer для файла для буферизации записи
	writer := bufio.NewWriter(outputFile)
	// Записываем данные в файл
	writer.Write([]byte(data))
	// Сбрасываем данные из буфера на диск
	writer.Flush()
	finish, _ := os.Lstat(arg2)
	fmt.Printf("\nЗавершена запись результатов в файл: %v, %v байт", finish.Name(), finish.Size())

}
