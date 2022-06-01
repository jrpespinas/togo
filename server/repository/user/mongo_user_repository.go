package repository

import (
	"context"
	"errors"
	"os"
	"time"
	"togo/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodb struct {
	conn *mongo.Client
}

func NewMongoRepository(conn *mongo.Client) UserRepository {
	return &mongodb{
		conn: conn,
	}
}

// Register new User
func (db *mongodb) Register(user *models.User) (*models.User, error) {
	collection := db.conn.Database(os.Getenv("DATABASE_NAME")).Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *mongodb) Login(user *models.User) error {
	collection := db.conn.Database(os.Getenv("DATABASE_NAME")).Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	update := bson.D{{"$set", bson.D{{Key: "token", Value: user.Token}}}}

	result, _ := collection.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return errors.New("unable to Login user")
	} else {
		return nil
	}
}

// Get existing user by Email
func (db *mongodb) GetUser(email string, user *models.User) (models.User, error) {
	collection := db.conn.Database(os.Getenv("DATABASE_NAME")).Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var found models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&found)
	if err != nil {
		return found, err
	}
	return found, nil
}
