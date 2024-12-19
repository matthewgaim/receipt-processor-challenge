// TODO: Fix function names & maybe 1 function per parameter

package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func createRandomId() string {
	id := uuid.New()
	return id.String()
}

func validateReceipt(r *Receipt) error {
	if r.Retailer == "" || r.PurchaseDate == "" || r.PurchaseTime == "" ||
		r.Total == "" || len(r.Items) == 0 {
		return fmt.Errorf("Missing required fields")
	}

	if !regexp.MustCompile(`^[\w\s\-&]+$`).MatchString(r.Retailer) {
		return fmt.Errorf("Invalid retailer format - %s", r.Retailer)
	}

	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
		return fmt.Errorf("Invalid date - %s", r.PurchaseDate)
	}

	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("Invalid time - %s", r.PurchaseTime)
	}

	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(r.Total) {
		return fmt.Errorf("Invalid total price - %s", r.Total)
	}

	for _, item := range r.Items {
		if !regexp.MustCompile(`^[\w\s\-]+$`).MatchString(item.ShortDescription) {
			return fmt.Errorf("Invalid short description - %s", item.ShortDescription)
		}
		if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(item.Price) {
			return fmt.Errorf("Invalid price - %s", item.Price)
		}
	}

	return nil
}

func getRetailerPoints(retailer string) int {
	retailer_points := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			retailer_points += 1
		}
	}
	return retailer_points
}

func getTotalIsRoundDollarPoints(total string) int {
	total_str := strings.Split(total, ".")
	if total_str[1] == "00" {
		return 50
	}
	return 0
}

func getTotalIsMultipleOfQuartersPoints(total string) int {
	total_str := strings.Split(total, ".")
	cents_str := total_str[1]
	cents, _ := strconv.Atoi(cents_str)

	if cents%25 == 0 {
		return 25
	}
	return 0
}

func getTwoItemsAtATimePoints(len_items int) int {
	return 5 * (len_items / 2)
}

// Trim here, not in main
func getTrimmedDescriptionLengthIsMultipleOfThreePoints(item Item) int {
	description := strings.TrimSpace(item.ShortDescription)
	price, _ := strconv.ParseFloat(item.Price, 64)
	if len(description)%3 == 0 {
		return int(math.Ceil(price * 0.2))
	}
	return 0
}

func getPurchaseDateIsOddPoints(purchaseDate string) int {
	date_arr := strings.Split(purchaseDate, "-")
	day_str := date_arr[2]
	day, _ := strconv.Atoi(day_str)
	if day%2 != 0 {
		return 6
	}
	return 0
}

func getPurchaseTimeBetweenTwoAndFourPoints(purchaseTime string) int {
	time_arr := strings.Split(purchaseTime, ":")
	hour, _ := strconv.Atoi(time_arr[0])
	minutes, _ := strconv.Atoi(time_arr[1])

	if (hour == 14 && minutes > 0) || (hour == 15 && minutes < 60) {
		return 10
	}

	return 0
}
