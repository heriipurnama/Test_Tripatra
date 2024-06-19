package handlers

import (
	"context"
	"log"
	"procurement-app-backend/models"
	"procurement-app-backend/utils"
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

func GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

	// Fetch the user's email
	user, err := GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	go utils.OrderCreatedEmail(user.Email, purchaseOrder.OrderID, purchaseOrder.TotalAmount)

	return purchaseOrder, nil
}

func GetPurchaseOrders(ctx context.Context) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder

	// Define the options for the query, sorting by CreatedAt in descending order
	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	log.Println("Fetching purchase orders with options:", findOptions)

	cursor, err := purchaseOrderCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Println("Error finding documents:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	for cursor.Next(ctx) {
		var po models.PurchaseOrder
		if err = cursor.Decode(&po); err != nil {
			log.Println("Error decoding document:", err)
			return nil, err
		}
		// Format the CreatedAt field
		createdAtString := po.CreatedAt.Format("2006-01-02 15:04:05")
		log.Println("Formatted CreatedAt:", createdAtString)

		purchaseOrders = append(purchaseOrders, &po)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	log.Println("Fetched purchase orders successfully:", purchaseOrders)
	return purchaseOrders, nil
}
