package main

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestSearchEvents(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	// Define test cases
	tests := []struct {
		name       string
		searchTerm string
		mockResp   []bson.D
		wantCount  int
		wantTitle  string
	}{
		{
			name:       "Search for 'Party'",
			searchTerm: "Party",
			mockResp: []bson.D{
				{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Birthday Party"}},
			},
			wantCount: 1,
			wantTitle: "Birthday Party",
		},
		{
			name:       "Case insensitive search",
			searchTerm: "party",
			mockResp: []bson.D{
				{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Birthday Party"}},
			},
			wantCount: 1,
			wantTitle: "Birthday Party",
		},
		{
			name:       "Empty search returns all events",
			searchTerm: "",
			mockResp: []bson.D{
				{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Birthday Party"}},
				{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Tech Conference"}},
				{{Key: "_id", Value: primitive.NewObjectID()}, {Key: "title", Value: "Wedding"}},
			},
			wantCount: 3,
			wantTitle: "Birthday Party", // Just check first result
		},
		{
			name:       "Non-existent event",
			searchTerm: "NonExistent",
			mockResp:   []bson.D{},
			wantCount:  0,
			wantTitle:  "",
		},
	}

	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			dbClient = mt.Client

			cursor := mtest.CreateCursorResponse(1, "eventbrite_clone.events", mtest.FirstBatch, tt.mockResp...)
			killCursor := mtest.CreateCursorResponse(0, "eventbrite_clone.events", mtest.NextBatch)
			mt.AddMockResponses(cursor, killCursor)

			events, err := searchEvents(tt.searchTerm)
			if err != nil {
				t.Errorf("searchEvents() error = %v", err)
				return
			}

			if len(events) != tt.wantCount {
				t.Errorf("searchEvents() got %d events, want %d", len(events), tt.wantCount)
			}

			if tt.wantCount > 0 && events[0].Title != tt.wantTitle {
				t.Errorf("searchEvents() first event title = %s, want %s", events[0].Title, tt.wantTitle)
			}
		})
	}
}
