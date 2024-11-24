package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/NagaHarshita/receipt-processor-challenge/models"
    "github.com/google/uuid"
    "math"
	"strings"
	"time"
    "strconv"
)

var receiptStore = make(map[string]models.Receipt)
var pointsStore = make(map[string]int)

// SubmitReceipt processes the receipt submission
func SubmitReceipt(c *fiber.Ctx) error {
    var receipt models.Receipt
    if err := c.BodyParser(&receipt); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receipt format"})
    }

    // Generate a unique ID for the receipt (for simplicity, using length here)
    id := uuid.New().String()
    receiptStore[id] = receipt
    pointsStore[id] = calculatePoints(receipt)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

// GetReceiptPoints retrieves points for a given receipt ID
func GetReceiptPoints(c *fiber.Ctx) error {
    id := c.Params("id")
    points, exists := pointsStore[id]
    if !exists {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receipt not found"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"points": points})
}

func calculatePoints(receipt models.Receipt) int {
    points := 0

	// 1. One point for every alphanumeric character in the retailer name
	for _, char := range receipt.Retailer {
		if isAlphanumeric(char) {
			points++
		}
	}

	// 2. 50 points if the total is a round dollar amount with no cents
	if strings.HasSuffix(receipt.Total, ".00") {
		points += 50
	}

	// 3. 25 points if the total is a multiple of 0.25
	if totalFloat, err := parseToFloat(receipt.Total); err == nil && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// 5. Points based on item description length
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			if priceFloat, err := parseToFloat(item.Price); err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}

	// 6. 6 points if the day in the purchase date is odd
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if date.Day()%2 != 0 {
			points += 6
		}
	}

	// 7. 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if purchaseTime.Hour() == 14 || (purchaseTime.Hour() == 15 && purchaseTime.Minute() < 60) {
			points += 10
		}
	}

	return points
}

// Helper function to check if a character is alphanumeric
func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

// Helper function to parse string to float
func parseToFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
