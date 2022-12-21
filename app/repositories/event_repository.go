package repositories

import (
	"context"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventRepositories struct {
	*mongo.Client
}

func (r *EventRepositories) CreateEvent(event models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.EventCollection(r.Client).InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepositories) GetEventById(eventId string) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event := &models.Event{}

	objId, _ := primitive.ObjectIDFromHex(eventId)

	err := helpers.EventCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepositories) GetEventsByClassroomId(classroomId string) (*[]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(classroomId)

	cursor, err := helpers.EventCollection(r.Client).Find(ctx, bson.M{"_classroom_id": objId})
	if err != nil {
		return nil, err
	}

	events := &[]models.Event{}
	cursor.All(ctx, events)

	return events, nil
}

func (r *EventRepositories) UpdateEvent(event models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.EventCollection(r.Client).UpdateByID(ctx, event.Id, bson.D{
		{Key: "$set", Value: event},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepositories) DeleteEvent(eventId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)
	_, err := helpers.EventCollection(r.Client).DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}
