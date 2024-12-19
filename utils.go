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
			return fmt.Errorf("Invalid item price - %s", item.Price)
		}
	}

	return nil
}

// One point for every alphanumeric character in the retailer name
func retailerPoints(retailer string) int {
	retailer_points := 0
	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			retailer_points += 1
		}
	}
	fmt.Printf("%d points - retailerPoints\n", retailer_points)
	return retailer_points
}

// 50 points if the total is a round dollar amount with no cents
func roundDollarPoints(total string) int {
	total_str := strings.Split(total, ".")
	if total_str[1] == "00" {
		fmt.Println("50 points - roundDollarPoints")
		return 50
	}
	return 0
}

// 25 points if the total is a multiple of 0.25
func totalQuartersPoints(total string) int {
	total_str := strings.Split(total, ".")
	cents_str := total_str[1]
	cents, _ := strconv.Atoi(cents_str)

	if cents%25 == 0 {
		fmt.Println("25 points - totalQuartersPoints")
		return 25
	}
	return 0
}

// 5 points for every two items on the receipt
func itemPairPoints(len_items int) int {
	two_items_at_time_points := 5 * (len_items / 2)
	if two_items_at_time_points > 0 {
		fmt.Printf("%d points - itemPairPoints\n",
			two_items_at_time_points)
	}
	return two_items_at_time_points
}

/*
If the trimmed length of the item description is a
multiple of 3,multiply the price by 0.2 and round up
to the nearest integer. The result is the number of points earned
*/
func descPoints(item Item) int {
	description := strings.TrimSpace(item.ShortDescription)
	price, _ := strconv.ParseFloat(item.Price, 64)

	if len(description)%3 == 0 {
		point := int(math.Ceil(price * 0.2))
		fmt.Printf("%d points - descPoints\n", point)
		return point
	}
	return 0
}

// 6 points if the day in the purchase date is odd.
func oddDayPoints(purchaseDate string) int {
	date_arr := strings.Split(purchaseDate, "-")
	day_str := date_arr[2]
	day, _ := strconv.Atoi(day_str)

	if day%2 != 0 {
		fmt.Println("6 points - oddDayPoints")
		return 6
	}
	return 0
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func purchaseTimePoints(purchaseTime string) int {
	time_arr := strings.Split(purchaseTime, ":")
	hour, _ := strconv.Atoi(time_arr[0])
	minutes, _ := strconv.Atoi(time_arr[1])

	if (hour == 14 && minutes > 0) || (hour == 15 && minutes < 60) {
		fmt.Println("10 points - purchaseTimePoints")
		return 10
	}
	return 0
}
