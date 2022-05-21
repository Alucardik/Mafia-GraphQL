package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MafiaGQL_server/graph/generated"
	"MafiaGQL_server/graph/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) StartSession(ctx context.Context, input model.NewGameSession) (*model.GameSession, error) {
	session := model.GameSession{
		ID:           primitive.NewObjectID().Hex(),
		Name:         input.Name,
		Ongoing:      true,
		Participants: []string{input.Initiator},
	}

	return r.DbHandle.CreateSession(ctx, &session)
}

func (r *mutationResolver) AddParticipant(ctx context.Context, input model.NewParticipant) (*model.GameSession, error) {
	parsedID, err := primitive.ObjectIDFromHex(input.SessionID)
	if err != nil {
		return nil, err
	}

	session, err := r.DbHandle.GetSessionById(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	if !session.Ongoing {
		return nil, errors.New("session has already been terminated, you can't add new participants")
	}

	session.Participants = append(session.Participants, input.UserID)
	err = r.DbHandle.UpdateSessionById(ctx, parsedID, session)
	if err != nil {
		return nil, err
	}

	return session, err

}

func (r *mutationResolver) AddComment(ctx context.Context, input model.NewComment) (string, error) {
	comment := model.Comment{
		SessionID: input.SessionID,
		Author:    input.Author,
		Contents:  input.Contents,
	}

	_, err := r.DbHandle.AddComment(ctx, &comment)
	if err != nil {
		return "", err
	}

	return "Successfully added comment", err
}

func (r *mutationResolver) EndSession(ctx context.Context, sessionID string) (string, error) {
	parsedID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return "", err
	}

	session, err := r.DbHandle.GetSessionById(ctx, parsedID)
	if err != nil {
		return "", err
	}

	if !session.Ongoing {
		return "", errors.New("session has already been terminated")
	}

	session.Ongoing = false
	err = r.DbHandle.UpdateSessionById(ctx, parsedID, session)
	if err != nil {
		return "", err
	}

	return "Session terminated", err
}

func (r *queryResolver) Sessions(ctx context.Context, ongoing *bool, sessionID *string) ([]*model.GameSession, error) {
	if sessionID != nil {
		parsedID, err := primitive.ObjectIDFromHex(*sessionID)
		if err != nil {
			return nil, err
		}
		res, err := r.DbHandle.GetSessionById(ctx, parsedID)
		return []*model.GameSession{res}, err
	}

	return r.DbHandle.GetSessionsByStatus(ctx, *ongoing)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
