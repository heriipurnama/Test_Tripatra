package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"procurement-app-backend/graphql"
	"procurement-app-backend/middleware"
	"procurement-app-backend/models"
	"strings"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8080"

var collection *mongo.Collection

func init() {
	// Replace with your actual MongoDB URI
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	collection = client.Database("procurementdb").Collection("purchaseorders")
}

func fetchPurchaseOrders(ctx context.Context) ([]*models.PurchaseOrder, error) {
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var purchaseOrders []*models.PurchaseOrder
	for cur.Next(ctx) {
		var po models.PurchaseOrder
		if err := cur.Decode(&po); err != nil {
			return nil, err
		}
		purchaseOrders = append(purchaseOrders, &po)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return purchaseOrders, nil
}

func generatePDF(purchaseOrders []*models.PurchaseOrder) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Purchase Orders")

	pdf.SetFontSize(12)
	pdf.Ln(20)

	// Define headers
	headers := []string{"Order ID", "User ID", "Items", "Total Amount", "Created At"}

	// Determine max lengths for each column
	colWidths := make([]float64, len(headers))
	for i, header := range headers {
		colWidths[i] = float64(len(header))
	}
	for _, po := range purchaseOrders {
		colWidths[0] = max(colWidths[0], float64(len(po.OrderID)))
		colWidths[1] = max(colWidths[1], float64(len(po.UserID)))
		colWidths[2] = max(colWidths[2], float64(maxItemWidth(po.Items)))
		colWidths[3] = max(colWidths[3], float64(len(fmt.Sprintf("%.2f", po.TotalAmount))))
		colWidths[4] = max(colWidths[4], float64(len(po.CreatedAt.Format("2006-01-02 15:04:05"))))
	}

	// Convert lengths to widths
	for i := range colWidths {
		colWidths[i] *= 2.5 // Adjust this multiplier as needed for better fitting
	}

	// Print headers
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Print rows
	for _, po := range purchaseOrders {
		itemLines := formatItems(po.Items, colWidths[2])
		itemHeight := pdf.GetStringWidth(itemLines) / colWidths[2] * 10
		if pdf.GetY()+itemHeight > 250 { // Adjust this value as per your requirement
			pdf.AddPage()
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
		pdf.CellFormat(colWidths[0], 10, po.OrderID, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 10, po.UserID, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[2], 10, itemLines, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[3], 10, fmt.Sprintf("%.2f", po.TotalAmount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 10, po.CreatedAt.Format("2006-01-02 15:04:05"), "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func maxItemWidth(items []models.Item) int {
	maxWidth := 0
	for _, item := range items {
		width := len(fmt.Sprintf("Name: %s\nQuantity: %d\nPrice: %.2f\n\n", item.Name, item.Quantity, item.Price))
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func formatItems(items []models.Item, maxWidth float64) string {
	var formattedItems string
	for i, item := range items {
		itemInfo := fmt.Sprintf("Name: %s\nQuantity: %d\nPrice: %.2f", item.Name, item.Quantity, item.Price)
		formattedItems += itemInfo
		// Add a newline after each item except for the last one
		if i < len(items)-1 {
			formattedItems += "\n\n"
		}
	}
	return formattedItems
}

func splitLines(s string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(s)
	var currentLine string
	for _, word := range words {
		if getStringWidth(currentLine+word) > maxWidth {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func getStringWidth(s string) float64 {
	// Replace this with the actual implementation to calculate the width of the string
	// using the gofpdf package or any other method you prefer
	// For example, if you're using gofpdf, you can use pdf.GetStringWidth(s)
	return 0.0
}

func getStringHeight(s string, maxWidth float64) float64 {
	lines := splitLines(s, maxWidth)
	return float64(len(lines)) * 10 // Adjust the line height as needed
}

func downloadPurchaseOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Fetch purchase orders data
	purchaseOrders, err := fetchPurchaseOrders(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch purchase orders", http.StatusInternalServerError)
		return
	}

	// Generate PDF
	pdfBytes, err := generatePDF(purchaseOrders)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Encode PDF to base64
	base64EncodedPDF := base64.StdEncoding.EncodeToString(pdfBytes)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"base64PDF": "%s"}`, base64EncodedPDF)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Inisialisasi server GraphQL
	srv := graphql.NewGraphQLServer()

	// Konfigurasi CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// Membuat router baru
	mux := http.NewServeMux()

	// Menambahkan endpoint untuk GraphQL playground
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	// Menambahkan endpoint untuk GraphQL API
	mux.Handle("/graphql", middleware.AuthMiddleware(srv))

	// Endpoint untuk mengunduh purchase orders
	mux.HandleFunc("/download-purchase-orders", downloadPurchaseOrdersHandler)

	// Handler untuk semua permintaan
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Serve permintaan HTTP menggunakan router mux
		mux.ServeHTTP(w, r)
	}

	// Menetapkan handler CORS
	corsHandler := c.Handler(http.HandlerFunc(handler))

	// Mulai server HTTP
	log.Printf("Server berjalan di http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
