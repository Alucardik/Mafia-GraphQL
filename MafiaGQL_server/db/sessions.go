package db

import (
	"MafiaGQL_server/graph/model"
	"MafiaGQL_server/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (dh *mongoDbHandle) CreateSession(ctx context.Context, session *model.GameSession) (*model.GameSession, error) {
	insertionRes, err := dh.sessions.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	var res model.GameSession
	err = dh.sessions.FindOne(ctx, bson.D{{"_id", insertionRes.InsertedID}}).Decode(&res)
	return &res, err
}

func (dh *mongoDbHandle) GetAllSessions(ctx context.Context) ([]*model.GameSession, error) {
	res := make([]*model.GameSession, 0)

	iterator, err := dh.sessions.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for iterator.Next(context.TODO()) {
		// A new session variable should be declared for each document.
		var session model.GameSession
		err := iterator.Decode(&session)
		utils.PanicOnError("Failed to parse session iterator", err)

		parsedId, err := primitive.ObjectIDFromHex(session.ID)
		utils.PanicOnError("Failed to parse session id", err)

		comments, err := dh.GetSessionComments(ctx, parsedId)
		if err != nil {
			return nil, err
		}

		session.Comments = comments
		res = append(res, &session)
	}

	return res, err
}

func (dh *mongoDbHandle) GetSessionsByStatus(ctx context.Context, ongoing bool) ([]*model.GameSession, error) {
	res := make([]*model.GameSession, 0)

	iterator, err := dh.sessions.Find(ctx, bson.D{{"ongoing", ongoing}})
	if err != nil {
		return nil, err
	}

	for iterator.Next(context.TODO()) {
		// A new session variable should be declared for each document.
		var session model.GameSession
		err := iterator.Decode(&session)
		utils.PanicOnError("Failed to parse session iterator", err)

		parsedId, err := primitive.ObjectIDFromHex(session.ID)
		utils.PanicOnError("Failed to parse session id", err)

		comments, err := dh.GetSessionComments(ctx, parsedId)
		if err != nil {
			return nil, err
		}

		session.Comments = comments
		res = append(res, &session)
	}

	return res, err
}

func (dh *mongoDbHandle) GetSessionById(ctx context.Context, id primitive.ObjectID) (*model.GameSession, error) {
	res := model.GameSession{}

	err := dh.sessions.FindOne(ctx, bson.D{{"_id", id.Hex()}}).Decode(&res)
	if err != nil {
		return nil, err
	}

	comments, err := dh.GetSessionComments(ctx, id)
	res.Comments = comments

	return &res, err
}

func (dh *mongoDbHandle) UpdateSessionById(ctx context.Context, id primitive.ObjectID, updated *model.GameSession) error {
	_, err := dh.sessions.ReplaceOne(ctx, bson.D{{"_id", id.Hex()}}, updated)
	return err
}
