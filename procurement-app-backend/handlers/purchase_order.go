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

var purchaseOrderCollection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	purchaseOrderCollection = client.Database("procurementdb").Collection("purchaseorders")
}

func CreatePurchaseOrder(ctx context.Context, userId string, items []models.Item) (*models.PurchaseOrder, error) {
	var itemList []models.Item
	totalAmount := 0.0
	for _, item := range items {
		totalAmount += float64(item.Quantity) * item.Price
		itemList = append(itemList, models.Item{
			ItemID:   item.ItemID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	purchaseOrder := &models.PurchaseOrder{
		OrderID:     primitive.NewObjectID().Hex(),
		UserID:      userId,
		Items:       itemList,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
	}

	_, err := purchaseOrderCollection.InsertOne(ctx, purchaseOrder)
	if err != nil {
		return nil, err
	}

	// go utils.SendEmail(userId, "Purchase Order Created", "Your purchase order has been successfully created.")

	return purchaseOrder, nil
}

func GetPurchaseOrders(ctx context.Context) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder
	cursor, err := purchaseOrderCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var po models.PurchaseOrder
		if err = cursor.Decode(&po); err != nil {
			return nil, err
		}
		purchaseOrders = append(purchaseOrders, &po)
	}
	return purchaseOrders, nil
}
