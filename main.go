package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type ReceiptProcessResponse struct {
	Id string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

// TODO: Change those variable names
func main() {
	points := make(map[string]int)

	http.HandleFunc("POST /receipts/process", func(w http.ResponseWriter, r *http.Request) {
		var receipt Receipt
		total_points := 0

		json.NewDecoder(r.Body).Decode(&receipt)

		if err := validateReceipt(&receipt); err != nil {
			fmt.Printf("Receipt not valid: %s\n", err.Error())
			http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
			return
		}

		retailer_points := getRetailerPoints(receipt.Retailer)
		total_points += retailer_points
		fmt.Printf("%d points - retailer_points\n", retailer_points)

		total_is_round_dollar_points := getTotalIsRoundDollarPoints(receipt.Total)
		total_points += total_is_round_dollar_points
		fmt.Printf("%d points - total_is_round_dollar_points\n", total_is_round_dollar_points)

		total_is_multiple_of_quarters_points := getTotalIsMultipleOfQuartersPoints(receipt.Total)
		total_points += total_is_multiple_of_quarters_points
		fmt.Printf("%d points - total_is_multiple_of_quarters_points\n", total_is_multiple_of_quarters_points)

		length_of_items := len(receipt.Items)
		two_items_at_time_points := getTwoItemsAtATimePoints(length_of_items)
		total_points += two_items_at_time_points
		fmt.Printf("%d points - two_items_at_time_points\n", two_items_at_time_points)

		for _, item := range receipt.Items {
			total_points += getTrimmedDescriptionLengthIsMultipleOfThreePoints(item)
		}

		purchase_day_odd_points := getPurchaseDateIsOddPoints(receipt.PurchaseDate)
		total_points += purchase_day_odd_points
		fmt.Printf("%d points - purchase_day_odd_points\n", purchase_day_odd_points)

		purchase_time_btwn_two_four_points := getPurchaseTimeBetweenTwoAndFourPoints(receipt.PurchaseTime)
		total_points += purchase_time_btwn_two_four_points
		fmt.Printf("%d points - purchase_time_btwn_two_four_points\n", purchase_time_btwn_two_four_points)
		fmt.Printf("Total Points - %d\n", total_points)

		rand_id := createRandomId()
		points[rand_id] = total_points

		response := ReceiptProcessResponse{
			Id: rand_id,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("GET /receipts/{id}/points", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		value, exists := points[id]
		if exists {
			response := PointsResponse{
				Points: value,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this route is working :)")
	})

	fmt.Println("Listening on http://localhost:3000/")
	http.ListenAndServe(":3000", nil)
}
