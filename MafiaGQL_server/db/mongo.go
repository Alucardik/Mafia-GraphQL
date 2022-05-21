package db

import (
	"MafiaGQL_server/graph/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbHandle interface {
	InitConnection(username, password, host string, port int) error

	CreateSession(ctx context.Context, session *model.GameSession) (*model.GameSession, error)
	GetAllSessions(ctx context.Context) ([]*model.GameSession, error)
	GetSessionsByStatus(ctx context.Context, ongoing bool) ([]*model.GameSession, error)
	GetSessionById(ctx context.Context, id primitive.ObjectID) (*model.GameSession, error)
	UpdateSessionById(ctx context.Context, id primitive.ObjectID, updated *model.GameSession) error

	AddComment(ctx context.Context, comment *model.Comment) (*mongo.InsertOneResult, error)
	GetSessionComments(ctx context.Context, sessionId primitive.ObjectID) ([]*model.Comment, error)
}

type mongoDbHandle struct {
	ctx      context.Context
	client   *mongo.Client
	mafiaDb  *mongo.Database
	users    *mongo.Collection
	comments *mongo.Collection
	sessions *mongo.Collection
}

func (dh *mongoDbHandle) InitConnection(username, password, host string, port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), _CONNECTION_TM)
	defer cancel()

	dbAddress := fmt.Sprintf("mongodb://%s:%s@%s:%d/", username, password, host, port)
	clientOptions := options.Client().ApplyURI(dbAddress)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	dh.client = client

	err = dh.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	dh.mafiaDb = dh.client.Database("mafiaGraphQL")
	dh.users = dh.mafiaDb.Collection("users")
	_, err = dh.users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	dh.sessions = dh.mafiaDb.Collection("sessions")
	dh.comments = dh.mafiaDb.Collection("comments")
	return nil
}

func CreateMongoDBHandle() MongoDbHandle {
	return &mongoDbHandle{}
}
