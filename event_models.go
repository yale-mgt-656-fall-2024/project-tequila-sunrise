package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event - encapsulates information about an event
type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Location  string             `bson:"location" json:"location"`
	Image     string             `bson:"image" json:"image"`
	Date      time.Time          `bson:"date" json:"date"`
	Attending []string           `bson:"attending" json:"attending"`
}

// getEventByID - returns the event with the specified id
func getEventByID(id string) (Event, bool) {
	collection := getCollection("events")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Event{}, false
	}

	var event Event
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		return Event{}, false
	}
	return event, true
}

// getAllEvents - returns all events
func getAllEvents() ([]Event, error) {
	collection := getCollection("events")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []Event
	for cursor.Next(context.TODO()) {
		var event Event
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// addAttendee - Adds an attendee to an event
func addAttendee(id string, email string) error {
	collection := getCollection("events")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$push": bson.M{"attending": email}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// addEvent - Adds an event to the collection
func addEvent(event Event) error {
	collection := getCollection("events")
	_, err := collection.InsertOne(context.TODO(), event)
	return err
}

// searchEvents - returns events matching the search term in their title
func searchEvents(searchTerm string) ([]Event, error) {
    collection := getCollection("events")
    
    var filter bson.M
    if searchTerm == "" {
        filter = bson.M{} // empty filter to match all documents
    } else {
        filter = bson.M{
            "title": bson.M{
                "$regex":   searchTerm,
                "$options": "i",
            },
        }
    }
    
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    var events []Event
    for cursor.Next(context.TODO()) {
        var event Event
        if err := cursor.Decode(&event); err != nil {
            return nil, err
        }
        events = append(events, event)
    }
    
    return events, nil
}
