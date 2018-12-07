package main

// resources:
// - https://www.compose.com/articles/mongodb-and-go-moving-on-from-mgo/
// - https://gitlab.com/wemgl/todocli/blob/master/main.go
// - https://godoc.org/github.com/mongodb/mongo-go-driver/bson

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/brianvoe/gofakeit"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

// grab a var for the connection string from the command line.
var uri = flag.String("uri", "mongodb://localhost", "The URI of the MongoDB instance you want to connect to")

// define database
var database = flag.String("db", "bank", "The database to work in")

// define destination collection
var collName = flag.String("coll", "customers", "the collection to write into")

// drop or append?
var drop = flag.Bool("drop", true, "Drop existing collection? false to append.")

// How many docs to read and write at once as part of a bulk insert
var batchSize = flag.Int("batchSize", 1000, "the number of documents to process in one batch")

// How many docs to read and write at once as part of a bulk insert
var docsLeft = flag.Int("docsToCreate", 100000, "the number of documents to create in total")

func main() {
	fmt.Println("\n---------------------------------\n--Generate Fake Bank Customer Records\n---------------------------------")

	// parse the flags
	flag.Parse()

	fmt.Printf("\n--uri: %v -- db is: %v, collection is: %v", *uri, *database, *collName)
	fmt.Printf("\n--drop collection:%v\n", *drop)

	// create a context. (note to self, learn what a context is...)
	ctx := context.Background()

	// create a client for the DB
	client, err := mongo.NewClient(*uri)
	if err != nil {
		log.Fatal(err)
	}

	// Connect the client to the DB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// select the database to use
	db := client.Database(*database)

	if *drop {
		// drop the collection (naughty! must deal with errors!)
		fmt.Printf("\n--Dropping Collection:%v\n", *collName)
		_ = db.Collection(*collName).Drop(ctx)
	} else {
		fmt.Printf("\n--NOT Dropping Collection\n")
	}

	createDocs(ctx, db, *batchSize, *docsLeft)

}

// A little function for making a percentage, constraines to 2 decimal places
func pct(min, max float64) float64 {
	return math.Round((rand.Float64()*(max-min)+min)*1000) / 1000
}

// Another little function for creating random balances,
// esentialy using a StdDev model, but with min+max clamps
func bal(seed float64) float64 {

	nf := rand.NormFloat64()

	switch {
	case nf > 0.001:
		nf = rand.Float64()
	case nf < -0.001:
		nf = rand.Float64()
	}
	v := math.Round((seed*nf+seed)*100) / 100
	return v
}

//create the documents
func createDocs(ctx context.Context, db *mongo.Database, batchSize int, docsLeft int) error {

	fmt.Printf("--starting to create documents in batches of %v\n", batchSize)

	// some arrays of various options
	countries := []string{"EN", "EN", "EN", "EN", "FR", "FR", "DE", "DE", "IT", "IT", "ES", "PT", "GR", "DN", "SE"}
	branch := []string{"EC-1", "EC-2", "EC-3", "EC-4"}

	for docsLeft > 0 {
		//create a slice to hold the batch of documents in.
		var docs []interface{}
		coll := db.Collection(*collName)
		for i := 0; i < batchSize; i++ {

			// create a new document with a bson.NewDocument object
			doc := bson.D{
				{"name", gofakeit.Name()},
				{"branch", gofakeit.StreetName() + " " + gofakeit.StreetSuffix()},
				{"branch_id", branch[rand.Intn(len(branch))]},
				{"manager", gofakeit.Name()},
				{"country", countries[rand.Intn(len(countries))]},
				{"rankLevel", rand.Intn(10)},
			}

			// which bank products should this user have?
			i := rand.Intn(6)

			switch i {
			case 0:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "Special"}, {"overdraftLimit", 1000}, {"balance", bal(1234)}}}}}
				doc = append(doc, docAccounts...)
			case 1:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "CurrentPlus"}, {"overdraftLimit", 1000}, {"balance", bal(1932)}},
					bson.D{{"accountType", "Savings"}, {"accountSubType", "SuperSaver"}, {"interestRate", pct(1, 4)}, {"balance", bal(32145)}}}}}
				doc = append(doc, docAccounts...)
			case 2:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "4U"}, {"overdraftLimit", 1000}, {"balance", bal(456)}},
					bson.D{{"accountType", "ISA"}, {"accountSubType", "SuperTaxFreeISA"}, {"interestRate", pct(1, 4)}, {"balance", bal(33456)}}}}}
				doc = append(doc, docAccounts...)
			case 3:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "Reserved"}, {"overdraftLimit", 1000}, {"balance", bal(3200)}},
					bson.D{{"accountType", "Mortgage"}, {"accountSubType", "BuildingDeluxe"}, {"interestRate", pct(2, 8)}, {"balance", math.Round((bal(123456)-123456*2)*100) / 100}}}}}
				doc = append(doc, docAccounts...)
			case 4:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "SuperSpecial"}, {"overdraftLimit", 1000}, {"balance", bal(1456)}},
					bson.D{{"accountType", "Savings"}, {"accountSubType", "SuperSaver"}, {"interestRate", pct(1, 4)}, {"balance", bal(12435)}},
					bson.D{{"accountType", "Mortgage"}, {"accountSubType", "BuildingDeluxe"}, {"interestRate", pct(2, 8)}, {"balance", math.Round((bal(142346)-142346*2)*100) / 100}}}}}
				doc = append(doc, docAccounts...)
			case 5:
				docAccounts := bson.D{{"accounts", bson.A{
					bson.D{{"accountType", "Current"}, {"accountSubType", gofakeit.HackerAdjective() + "CurrentAccount"}, {"overdraftLimit", bsonx.Double(1000.0)}, {"balance", bal(2673)}},
					bson.D{{"accountType", "Savings"}, {"accountSubType", "SuperSaver"}, {"interestRate", pct(1, 4)}, {"balance", bal(12345)}},
					bson.D{{"accountType", "ISA"}, {"accountSubType", "SuperTaxFreeISA"}, {"interestRate", pct(1, 4)}, {"balance", bal(23456)}},
					bson.D{{"accountType", "Mortgage"}, {"accountSubType", "BuildingDeluxe"}, {"interestRate", pct(2, 8)}, {"balance", math.Round((bal(234567)-234567*2)*100) / 100}}}}}
				doc = append(doc, docAccounts...)
			}

			docs = append(docs, doc)
		}

		// insert the docs into the DB.
		docsLeft = docsLeft - batchSize
		fmt.Printf("%v Docs left to process\n", docsLeft)

		_, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("could not insert: %v", err)
		}
	}
	return nil
}
