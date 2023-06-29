package models

import (
	"context"
	"log"

	//arango "github.com/joselitofilho/gorm-arango/pkg"
	// "github.com/arangodb/go-driver/http"
	// driver "github.com/arangodb/go-driver"
	"os"

	//arango "github.com/arangodb/go-driver"
	// "github.com/arangodb/go-driver/http"
	"fmt"
	// "github.com/fatih/structs"
	// "strings"
	// "os"

	// "reflect"

	driver "github.com/arangodb/go-driver"
	gohttp "github.com/arangodb/go-driver/http"
)

var conDBArango driver.Database
var conDBArangoErr error

func ArangoDb() (driver.Database, error) {
	if conDBArango == nil || conDBArangoErr != nil {
		conn, err := gohttp.NewConnection(gohttp.ConnectionConfig{
			Endpoints: []string{os.Getenv("golangURI")},
		})
		if err != nil {
			return nil, err
		}
		client, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(os.Getenv("golangUser"), os.Getenv("golangPassword")),
		})

		db, err := client.Database(nil, os.Getenv("golangDatabase"))

		if err != nil {
			conDBArangoErr = err
			return nil, err
		}
		conDBArango = db
		conDBArangoErr = err
	} else if conDBArango != nil && conDBArangoErr == nil {
		return conDBArango, conDBArangoErr
	}
	return conDBArango, conDBArangoErr
}

func ArangoDbInsert(collection string, model interface{}) {
	db, err := ArangoDb()
	if err != nil {
		log.Panic(err.Error())
	}
	col, err := db.Collection(nil, collection)
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("arango Insert " + collection)
	_, errs, err := col.CreateDocuments(nil, model)

	if err != nil {
		log.Panic(err.Error())
	} else if err := errs.FirstNonNil(); err != nil {
		log.Panic(err.Error())
	}

}

func ArangoDbUpdate(collection string, conditions string, model interface{}) {
	db, err := ArangoDb()
	if err != nil {
		log.Panic(err.Error())
	}

	ctx := context.Background()
	query := "FOR u IN " + collection + " FILTER " + conditions + "  RETURN u"
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		log.Panic(err.Error())
	}
	defer cursor.Close()
	for {
		var doc interface{}
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}
		col, err := db.Collection(ctx, collection)
		if err != nil {
			log.Panic(err.Error())
		}
		_, err = col.UpdateDocument(ctx, meta.Key, model)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	//
	fmt.Println("arango Update " + collection)
	//getting key
}

func ArangoDbRemove(collection string, conditions string, bindVars map[string]interface{}) {
	db, err := ArangoDb()
	if err != nil {
		log.Panic(err.Error())
	}
	ctx := context.Background()
	query := "FOR u IN " + collection + " FILTER " + conditions + " RETURN u"
	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		log.Panic(err.Error())
	}
	defer cursor.Close()
	for {
		var doc interface{}
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}
		col, err := db.Collection(ctx, collection)
		if err != nil {
			log.Panic(err.Error())
		}
		//getting key

		_, err = col.RemoveDocument(ctx, meta.Key)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	fmt.Println("arango Delete " + collection)
}
func ArangoSelectCollection(collection string, conditions string, bindVar map[string]interface{}) []map[string]interface{} {
	db, err := ArangoDb()
	if err != nil {
		log.Panic(err.Error())
	}
	ctx := context.Background()
	query := "FOR u IN " + collection + " FILTER " + conditions + " RETURN u"
	fmt.Println(query)
	cursor, err := db.Query(ctx, query, bindVar)
	if err != nil {
		log.Panic(err.Error())
	}
	defer cursor.Close()
	output := make([]map[string]interface{}, 0)
	for {
		var doc interface{}
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Panic(err.Error())
		}
		mymap := doc.(map[string]interface{})
		output = append(output, mymap)
	}
	return output
}
