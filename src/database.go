// Responsible for saving to database.
package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Binary document schema.
type Binary struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string
	Version       string
	Platform      string
	Path          string
	Sig           string
	Parts         int
	Locale        string
	Creation_date time.Time
}

// Part document schema.
type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string
	Version       string
	Platform      string
	Path          string
	Part_no       int
	Belongs_to    primitive.ObjectID
	Creation_date time.Time
}

// Checks if a document already exists in the database.
func db_doc_exists(binary_col *mongo.Collection, name string, platform string, version string, locale string) bool {
	mongodb_context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"name":     name,
		"platform": platform,
		"version":  version,
		"locale":   locale,
		// "parts":    parts_no,
	}

	// SetLimit is set to 1 so it stops as soon as it finds one.
	count, err := binary_col.CountDocuments(mongodb_context, filter, options.Count().SetLimit(1))
	check(err)

	return count > 0
}

// Inserts an array of Binary and an array of Part
// to the database.
func db_insert(binary_col *mongo.Collection, part_col *mongo.Collection, binaries []Binary, parts []Part) {
	log.Print("Saving to database...")

	mongodb_context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parts are added before Binaries
	// in case they are not fully added,
	// a frontend wont request the
	// the binary they belong to.
	if len(parts) > 0 {
		new_parts := make([]interface{}, len(parts))

		for i := range parts {
			new_parts[i] = parts[i]
		}

		part_col.InsertMany(mongodb_context, new_parts)
	}

	if len(binaries) > 0 {
		new_binaries := make([]interface{}, len(binaries))

		for i := range binaries {
			new_binaries[i] = binaries[i]
		}

		binary_col.InsertMany(mongodb_context, new_binaries)
	}

	log.Print("Finished saving to database.")

}
