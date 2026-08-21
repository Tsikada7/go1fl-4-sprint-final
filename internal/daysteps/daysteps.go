package daysteps

import (
	"fmt"
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

func parsePackage(data string) (int, time.Duration, error) {

	splitData := strings.Split(data, ",")
	if len(splitData) != 2 {
		return 0, 0, fmt.Errorf("Invalid data format: expected 2 elements, got %d", len(splitData))
	}

	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Conversion error: %v", err)
	}
	if steps < 1 {
		return 0, 0, fmt.Errorf("Number of steps cannot be zero, or less than zero")
	}

	duration, err := time.ParseDuration(splitData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Conversion error: %v", err)
	} else if duration < 1 {
		return 0, 0, fmt.Errorf("Duration cannot be less than zero, or be zero")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	} else if steps < 1 {
		return ""
	}

	distance := float64(steps) * stepLength
	distanceInKm := distance / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
}
