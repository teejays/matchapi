package db

import (
	"fmt"

	"github.com/teejays/clog"
	"github.com/teejays/gofiledb"

	"github.com/teejays/matchapi/lib/pk"
)

var documentRoot = ".data"
var client *gofiledb.Client
var isInitialized bool

var UserCollection string = "user"
var LikeCollection string = "like"

// InitDB initializes the database connection
func InitDB() error {
	var err error
	client, err = initClient(documentRoot)
	return err
}

func initClient(dir string) (*gofiledb.Client, error) {
	// Initialize the database
	o := gofiledb.ClientInitOptions{
		DocumentRoot:          dir,
		OverwritePreviousData: false,
	}
	err := gofiledb.Initialize(o)
	if err != nil {
		return nil, fmt.Errorf("could not initialize the gofiledb client: %v", err)
	}

	// Get the initialized client
	cl := gofiledb.GetClient()

	// Create the user collections
	err = cl.AddCollection(gofiledb.CollectionProps{Name: UserCollection, EncodingType: gofiledb.ENCODING_JSON})
	if err != nil {
		return nil, fmt.Errorf("could not create the gofiledb '%s' collection: %v", UserCollection, err)
	}

	// Create the like collections
	err = cl.AddCollection(gofiledb.CollectionProps{Name: LikeCollection, EncodingType: gofiledb.ENCODING_JSON})
	if err != nil {
		return nil, fmt.Errorf("could not create the gofiledb '%s' collection: %v", LikeCollection, err)
	}

	// Add an index on ReceiverID and GiverID for like
	err = cl.AddIndex(LikeCollection, "GiverID")
	if err != nil {
		return nil, fmt.Errorf("could not create the index 'GiverID' on '%s' collection: %v", LikeCollection, err)
	}
	err = cl.AddIndex(LikeCollection, "ReceiverID")
	if err != nil {
		return nil, fmt.Errorf("could not create the index 'ReceiverID' on '%s' collection: %v", LikeCollection, err)
	}

	return cl, nil
}

// GetClient provided the client object that can be used to interact with the database
func GetClient() *gofiledb.Client {

	if client == nil {
		panic("db.client is initialized but detected as nil")
	}

	return client
}

// GetEntityByID saves an entity by it's ID in the persistent storage
func GetEntityByID(collection string, key pk.ID, addr interface{}) error {

	// Create a read lock on the collection
	rlock(collection)
	defer runlock(collection)

	cl := GetClient()
	return cl.GetStruct(collection, gofiledb.Key(key), addr)

}

// SaveEntityByID saves an entity by it's ID in the persistent storage
func SaveEntityByID(collection string, key pk.ID, entity interface{}) error {

	// Create a lock on the collection
	lock(collection)
	defer unlock(collection)

	cl := GetClient()
	return cl.SetStruct(collection, gofiledb.Key(key), entity)

}

// SaveNewEntity ...
func SaveNewEntity(collection string, entity interface{}) (pk.ID, error) {

	// Create a lock on the collection
	lock(collection)
	defer unlock(collection)

	clog.Debugf("DB | SaveNewEntity: Saving new %s entity...", collection)
	cl := GetClient()
	id, err := cl.SaveNewEntity(collection, entity)

	return pk.ID(id), err
}

// Query ..
func Query(collection string, query string) ([]interface{}, error) {

	// Get the client
	cl := GetClient()

	// Run the query
	resp, err := cl.Search(collection, query)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}
