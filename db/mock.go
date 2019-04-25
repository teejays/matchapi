package db

import (
	"github.com/teejays/gofiledb/util"
)

var mockDocumentRoot = ".data_mock"

// InitMockClient ...
func InitMockClient() error {
	var err error

	// create the mock data director if it doesn't exist
	err = util.CreateDirIfNotExist(mockDocumentRoot)
	if err != nil {
		return err
	}

	// initiialize and overwrite the existing client
	client, err = initClient(mockDocumentRoot)
	return err
}

// DestoryMockClient ...
func DestoryMockClient() error {
	cl := GetClient()
	isInitialized = false
	return cl.Destroy()
}
