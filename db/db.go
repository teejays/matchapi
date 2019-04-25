package db

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/teejays/gofiledb"
	"github.com/teejays/matchapi/lib/pk"
)

var documentRoot = ".data"
var client *gofiledb.Client

var UserCollection string = "user"
var LikeCollection string = "like"

// InitDB initializes the database connection
func InitDB() error {

	// Initialize the database
	o := gofiledb.ClientInitOptions{
		DocumentRoot:          documentRoot,
		OverwritePreviousData: false,
	}
	err := gofiledb.Initialize(o)
	if err != nil {
		return fmt.Errorf("could not initialize the gofiledb client: %v", err)
	}

	// Get the initialized client
	client = gofiledb.GetClient()

	// Create the user collections
	err = client.AddCollection(gofiledb.CollectionProps{Name: UserCollection, EncodingType: gofiledb.ENCODING_JSON})
	if err != nil {
		return fmt.Errorf("could not create the gofiledb '%s' collection: %v", UserCollection, err)
	}

	// Create the like collections
	err = client.AddCollection(gofiledb.CollectionProps{Name: LikeCollection, EncodingType: gofiledb.ENCODING_JSON})
	if err != nil {
		return fmt.Errorf("could not create the gofiledb '%s' collection: %v", LikeCollection, err)
	}

	// Add an index on ReceiverID and GiverID for like
	err = client.AddIndex(LikeCollection, "GiverID")
	if err != nil {
		return fmt.Errorf("could not create the index 'GiverID' on '%s' collection: %v", LikeCollection, err)
	}
	err = client.AddIndex(LikeCollection, "ReceiverID")
	if err != nil {
		return fmt.Errorf("could not create the index 'ReceiverID' on '%s' collection: %v", LikeCollection, err)
	}

	return nil
}

// GetClient provided the client object that can be used to interact with the database
func GetClient() *gofiledb.Client {
	if client == nil {
		panic("db.client fetched before it is initialized")
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

// SaveNewEntity saves a new enity
func SaveNewEntity(collection string, entity interface{}) (pk.ID, error) {

	// Get a new ID
	id, err := GetNewEntityID(collection)
	if err != nil {
		return id, err
	}

	// Create a lock on the collection
	lock(collection)
	defer unlock(collection)

	// Add the ID to the entity using reflect package

	// - get the reflect.Value of the entity
	v := reflect.ValueOf(entity)
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return -1, fmt.Errorf("Cannot set the value of the ID (%d) of the new %s entity: entity is not a struct", id, collection)
	}
	fv := v.FieldByName("ID")
	if !fv.IsValid() {
		return -1, fmt.Errorf("Cannot set the value of the ID (%d) of the new %s entity: field value is not valid", id, collection)
	}
	if !fv.CanSet() {
		return -1, fmt.Errorf("Cannot set the value of the ID (%d) of the new %s entity: cannot set the field value", id, collection)
	}
	// - get the reflect.Value of the ID
	vID := reflect.ValueOf(id)
	fv.Set(vID)

	// Save the new entity
	entity = v.Interface()
	cl := GetClient()
	err = cl.SetStruct(collection, gofiledb.Key(id), entity)
	if err != nil {
		return id, err
	}

	return id, nil
}

// GetNewEntityID generates a unique pk.ID for the given collection
func GetNewEntityID(collection string) (pk.ID, error) {
	// Get an random ID
	id := GetNewID()
	// Check if it already exists
	cl := GetClient()
	_, err := cl.GetFile(collection, gofiledb.Key(id))
	if gofiledb.IsNotExist(err) { // If the file doesn't exist, we're good to go
		return id, nil
	}
	if err != nil {
		return id, fmt.Errorf("generated the new id %d but could not verify that it is unique: %v", id, err)
	}
	return GetNewEntityID(collection)
}

// GetNewID generates a new unique ID for an entity
func GetNewID() pk.ID {
	minID := 100000
	rng := 100000
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	r := rand.New(src)
	id := r.Intn(rng)
	return pk.ID(id + minID)
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
