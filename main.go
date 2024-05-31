package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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

	// Открываем файл для чтения входных данных
	inputFile, err := os.OpenFile(arg1, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("ОШИБКА: вероятно файл с именем input.txt не найден.\n")
		return
	}
	defer inputFile.Close() // Отложенное закрытие файла после выполнения всех операций.

	// Если файл для записи уже существует, то очищаем его содержимое
	if _, err := os.Stat(arg2); err == nil {
		err := ioutil.WriteFile(arg2, []byte{}, 0644)
		if err != nil {
			log.Fatalf("ошибка очистки файла: %v\n", err)
		}
	}

	// Создаем читатель для файла для буферизации чтения из файла
	fileReader := bufio.NewReader(inputFile)

	var result int // Для подсчета значений выражений вне цикла
	// Выражение для поиска всех последовательностей из цифр в строке
	nums := regexp.MustCompile(`[0-9]+`)
	// Выражение для поиска арифметических операций в строке
	actions := regexp.MustCompile(`[+-/*]`)

	var buffer []byte
	for {
		line, _, err := fileReader.ReadLine()
		if err == io.EOF {
			break
		}
		// Проверка соответствия строки основному шаблону
		if re.Match(line) {
			numbers := nums.FindAllStringSubmatch(string(line), -1)
			action := actions.FindStringSubmatch(string(line))

			number1, err := strconv.Atoi(numbers[0][0])
			if err != nil {
				log.Fatalf("error converting %v to int: %v\n", numbers[0][0], err)
			}
			number2, err := strconv.Atoi(numbers[1][0])
			if err != nil {
				log.Fatalf("error converting %v to int: %v\n", numbers[1][0], err)
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
			data := strconv.Itoa(number1) + action[0] + strconv.Itoa(number2) + "=" + strconv.Itoa(result) + "\n"
			buffer = append(buffer, []byte(data)...)

			// Если условие выполняется, то вызывается функция writeToFile для записи данных из буфера в файл с именем, указанным в аргументе arg2
			// После записи данных в файл, буфер обнуляется, чтобы подготовить его для следующей порции данных
			if len(buffer) >= 1 {
				writeToFile(arg2, buffer)
				buffer = nil
			}
		}
	}
	finish, _ := os.Lstat(arg2)
	fmt.Printf("\nЗавершена запись результатов в файл: %v, %v байт", finish.Name(), finish.Size())
}

// writeToFile открывает файл с именем fileName для записи
// Если файл не существует, он будет создан. Если файл существует,
// данные будут добавлены в конец файла
// Функция записывает данные data в открытый файл и закрывает его после записи
// В случае ошибки при открытии файла или записи, функция сообщает об ошибке через логирование
func writeToFile(fileName string, data []byte) {
	outputFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("ошибка открытия файла для записи: %vn", err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(data)
	if err != nil {
		log.Fatalf("ошибка записи в файл: %vn", err)
	}
}
