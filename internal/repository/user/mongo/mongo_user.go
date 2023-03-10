package mongo

import (
	"context"
	"time"

	"github.com/piatoss3612/tx-noti-bot/internal/models"
	"github.com/piatoss3612/tx-noti-bot/internal/repository/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoUserRepository struct {
	name   string
	client *mongo.Client
}

func New(ctx context.Context, name, uri string) (user.UserRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &mongoUserRepository{client: client}, nil
}

func (m *mongoUserRepository) getCollection() *mongo.Collection {
	return m.client.Database(m.name).Collection("user")
}

func (m *mongoUserRepository) CreateUser(ctx context.Context, u *models.User) error {
	coll := m.getCollection()

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := coll.InsertOne(ctx, u)

	return err
}

func (m *mongoUserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	coll := m.getCollection()

	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}})

	cursor, err := coll.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var users []*models.User

	for cursor.Next(ctx) {
		var user models.User

		err = cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (m *mongoUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	coll := m.getCollection()

	var user models.User

	err := coll.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *mongoUserRepository) UpdateUser(ctx context.Context, u *models.User) error {
	coll := m.getCollection()

	_, err := coll.UpdateOne(ctx, bson.M{"id": u.ID}, bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "email", Value: u.Email},
				{Key: "discord_id", Value: u.DiscordID},
				{Key: "otp_enabled", Value: u.OtpEnabled},
				{Key: "otp_verified", Value: u.OtpVerified},
				{Key: "otp_secret", Value: u.OtpSecret},
				{Key: "otp_url", Value: u.OtpUrl},
				{Key: "updated_at", Value: time.Now()},
			},
		},
	})
	return err
}
func (m *mongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	coll := m.getCollection()

	_, err := coll.DeleteOne(ctx, bson.M{"id": id})

	return err
}

func (m *mongoUserRepository) Drop(ctx context.Context) error {
	coll := m.getCollection()

	return coll.Drop(ctx)
}
