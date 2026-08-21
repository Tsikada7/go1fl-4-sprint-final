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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	splitData := strings.Split(data, ",")
	if len(splitData) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: expected 3 elements, got %d", len(splitData))
	}

	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Conversion error: %v", err)
	}
	if steps < 1 {
		return 0, "", 0, fmt.Errorf("number of steps cannot be zero, or less than zero")
	}

	duration, err := time.ParseDuration(splitData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Conversion error: %v", err)
	}
	if duration < 1 {
		return 0, "", 0, fmt.Errorf("Duration cannot be less than 1")
	}

	return steps, splitData[1], duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient
	distanceInMetre := stepLength * float64(steps)
	distanceInKm := distanceInMetre / float64(mInKm)

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration < 1 || steps < 1 || height < 1 {
		return 0
	}

	distanceInKm := distance(steps, height)
	meanSpeed := distanceInKm / duration.Hours()

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)

	switch activity {
	case "Бег":
		runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", duration.Hours(), distance, meanSpeed, runningSpentCalories), nil
	case "Ходьба":
		walkingSpentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", duration.Hours(), distance, meanSpeed, walkingSpentCalories), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps < 1 || weight < 1 || height < 1 || duration < 1 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	caloriesSpent := (weight * meanSpeed * durationInMinutes) / minInH

	return caloriesSpent, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps < 1 || weight < 1 || height < 1 || duration < 1 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	caloriesSpent := (weight * meanSpeed * durationInMinutes) / minInH
	caloriesPercent := caloriesSpent * walkingCaloriesCoefficient

	return caloriesPercent, nil
}
