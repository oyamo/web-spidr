package downloader

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func checkChanges(sc *scan) {

	defer sc.wg.Done()
	*sc.filter = append(* sc.filter, bson.M{"url": * sc.res.url})
	*sc.docsFrmNet = append(*sc.docsFrmNet, sc.res)
	/**


	*sc.docsFrmNet = append(*sc.docsFrmNet, operation)
	// Cancel When the function completes
	**/

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	*sc.counter++

	// Read from DB THE URLS
	// Get the changed
	// Get the New Urls to be added
	// Add the urls

	if sc.maxLen == *sc.counter && len( * sc.filter ) != 0 {

		collection := sc.db.Collection("ChangeColl")
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
		defer cancelFunc()

		// Find All the Urls from the DB
		filter := bson.M{"$or": * sc.filter}
		query, err := collection.Find(ctx, filter)

		if err != nil {
			log.Fatal(err)
		} else {
			if err = query.All(ctx, sc.docsFrmDB); err != nil {
				log.Fatal(err)
			}

		}


		// Get the changed URLS
		if len(* sc.docsFrmDB) != 0 {
			vfrmDBMap := make(map[string] *Response)
			vfrmNetMap := make(map[string] *rawResponse)

			for i := 0; i < len(*sc.docsFrmDB); i++ {
				vfrmDBMap[(*sc.docsFrmDB)[i].Url] = &(*sc.docsFrmDB)[i]
			}

			for i := 0; i < len(*sc.docsFrmNet); i++ {
				vfrmNetMap[*((*sc.docsFrmNet)[i].url)] = &(*sc.docsFrmNet)[i]
			}


			bulkOperations := make([]mongo.WriteModel, 0)

			for k, v := range vfrmNetMap {

				if val, inDb := vfrmDBMap[k]; inDb {
					if * v.hash != val.Hash {
						* sc.changedPages = append(* sc.changedPages, *val)
					}
				}

				operation := mongo.NewUpdateOneModel()
				operation.SetUpsert(true)
				operation.SetFilter(bson.M{"url": *v.url})
				operation.SetUpdate(bson.M{
					"$set": bson.M{
						"hash":     *v.hash,
						"url":      *v.url,
						"modified": time.Now().Unix(),
						"title":    *v.title,
					},
				})

				bulkOperations = append(bulkOperations, operation)
			}


			ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second * 10)
			defer cancelFunc()

			_, err := collection.BulkWrite(ctx, bulkOperations, &bulkOption)
			if err != nil {
				panic(err)
			}
		} else {

			bulkOperations := make([]mongo.WriteModel, 0)

			for _, v := range *sc.docsFrmNet {

				* sc.changedPages = append(* sc.changedPages, Response{
					Title:        *v.title,
					Url:          *v.url,
					Hash:         *v.hash,
					DateModified: time.Now().Unix(),
				})
				operation := mongo.NewUpdateOneModel()
				operation.SetUpsert(true)
				operation.SetFilter(bson.M{"url": *v.url})
				operation.SetUpdate(bson.M{
					"$set": bson.M{
						"hash":     *v.hash,
						"url":      *v.url,
						"modified": time.Now().Unix(),
						"title":    *v.title,
					},
				})

				bulkOperations = append(bulkOperations, operation)
			}


			ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second * 10)
			defer cancelFunc()

			_, err := collection.BulkWrite(ctx, bulkOperations, &bulkOption)
			if err != nil {
				panic(err)
			}
		}
	}
	log.Println(len(*sc.docsFrmNet), "Pending operations")

}

