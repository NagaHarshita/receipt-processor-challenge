package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/NagaHarshita/receipt-processor-challenge/handlers"
)

func SetupRoutes(app *fiber.App) {
    app.Post("/receipts/process", handlers.SubmitReceipt)
    app.Get("/receipts/:id/points", handlers.GetReceiptPoints)
}
