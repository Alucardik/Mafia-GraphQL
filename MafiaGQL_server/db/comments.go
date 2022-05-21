package db

import (
	"MafiaGQL_server/graph/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (dh *mongoDbHandle) AddComment(ctx context.Context, comment *model.Comment) (*mongo.InsertOneResult, error) {

	sessionId, err := primitive.ObjectIDFromHex(comment.SessionID)
	if err != nil {
		return nil, errors.New("invalid session id")
	}

	if asscSession, err := dh.GetSessionById(ctx, sessionId); err != nil || asscSession == nil {
		return nil, errors.New("no session found")
	}

	res, err := dh.comments.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (dh *mongoDbHandle) GetSessionComments(ctx context.Context, sessionId primitive.ObjectID) ([]*model.Comment, error) {
	var res []*model.Comment

	iterator, err := dh.comments.Find(ctx, bson.D{{"sessionId", sessionId.Hex()}})
	if err != nil {
		return nil, err
	}

	if err = iterator.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, err
}
