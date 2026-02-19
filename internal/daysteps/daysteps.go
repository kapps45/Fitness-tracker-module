package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разбирает строку формата "678,0h50m".
//
// Возвращает:
//   - int — количество шагов;
//   - time.Duration — продолжительность прогулки;
//   - error — ошибку, если формат строки неверный
//     или значения не являются положительными.
func parsePackage(data string) (int, time.Duration, error) {
	// Разбиваем строку по запятой на слайс подстрок.
	parts := strings.Split(data, ",")

	// Проверяем, что получили ровно два элемента.
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	// Преобразуем количество шагов в int.
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("некорректное количество шагов")
	}

	// Преобразуем продолжительность в time.Duration.
	dur, err := time.ParseDuration(parts[1])
	if err != nil || dur <= 0 {
		return 0, 0, fmt.Errorf("некорректная длительность")
	}

	// Возвращаем шаги, продолжительность и nil.
	return steps, dur, nil
}

// DayActionInfo формирует строку с информацией о дневной активности.
//
// Возвращает:
//   - string — отформатированную строку с результатами;
//   - пустую строку, если входные данные некорректны
//     или произошла ошибка при расчёте калорий.
//
// Ошибки логируются внутри функции.
func DayActionInfo(data string, weight, height float64) string {
	// Парсим входную строку.
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Вычисляем дистанцию:
	// количество шагов * длина шага / 1000.
	distanceKm := float64(steps) * stepLength / mInKm

	// Вычисляем потраченные калории через пакет spentcalories.
	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Формируем и возвращаем строку с результатами.
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		cal,
	)
}
