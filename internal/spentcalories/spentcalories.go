package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	// lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining разбирает строку формата "шаги,тип,продолжительность".
//
// Возвращает:
//   - int — количество шагов;
//   - string — тип тренировки;
//   - time.Duration — продолжительность тренировки;
//   - error — ошибку, если формат строки неверный
//     или данные некорректны.
func parseTraining(data string) (int, string, time.Duration, error) {
	// Разбиваем строку по запятой на слайс подстрок.
	parts := strings.Split(data, ",")

	// Проверяем количество элементов.
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат тренировки")
	}

	// Парсим шаги.
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("некорректное количество шагов")
	}

	// Тип тренировки.
	activity := parts[1]

	// Парсим продолжительность.
	dur, err := time.ParseDuration(parts[2])
	if err != nil || dur <= 0 {
		return 0, "", 0, fmt.Errorf("некорректная длительность")
	}

	return steps, activity, dur, nil
}

// distance вычисляет пройденную дистанцию в километрах.
//
// Возвращает:
//   - float64 — дистанцию в километрах.
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient

	return float64(steps) * stepLen / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч.
//
// Возвращает:
//   - float64 — среднюю скорость в километрах в час.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

// RunningSpentCalories рассчитывает количество калорий,
// потраченных при беге.
//
// Возвращает:
//   - float64 — количество сожжённых калорий;
//   - error — ошибку, если входные параметры некорректны.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	cal := weight * speed * duration.Minutes() / minInH

	return cal, nil
}

// WalkingSpentCalories рассчитывает количество калорий,
// потраченных при ходьбе.
//
// Возвращает:
//   - float64 — количество сожжённых калорий
//     с учётом коэффициента ходьбы;
//   - error — ошибку, если входные параметры некорректны.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	cal, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return cal * walkingCaloriesCoefficient, nil
}

// TrainingInfo формирует строку отчёта о тренировке.
//
// Возвращает:
//   - string — отформатированную строку с результатами;
//   - error — ошибку, если данные некорректны
//     или тип тренировки неизвестен.
func TrainingInfo(data string, weight, height float64) (string, error) {
	// Парсим входную строку.
	steps, activity, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Общие расчёты.
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, dur)

	var cal float64

	switch activity {
	case "Бег":
		cal, err = RunningSpentCalories(steps, weight, height, dur)
	case "Ходьба":
		cal, err = WalkingSpentCalories(steps, weight, height, dur)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	// Формируем и возвращаем строку с результатами.
	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		dur.Hours(),
		dist,
		speed,
		cal,
	), nil
}
