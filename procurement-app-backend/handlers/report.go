package handlers

import (
	"context"
	"procurement-app-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var reportCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	reportCollection = client.Database("procurementdb").Collection("reports")
}

func GenerateReport(ctx context.Context, orderId string) (*models.Report, error) {
	var po models.PurchaseOrder
	err := purchaseOrderCollection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&po)
	if err != nil {
		return nil, err
	}

	report := &models.Report{
		ReportID:    primitive.NewObjectID().Hex(),
		OrderID:     orderId,
		GeneratedAt: time.Now(),
		ReportData:  "Report Data", // Placeholder for actual report data
	}

	_, err = reportCollection.InsertOne(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func GetReport(ctx context.Context, orderId string) (*models.Report, error) {
	var report models.Report
	err := reportCollection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}
