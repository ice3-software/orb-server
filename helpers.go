
package main

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func BSONOrb(t *testing.T) (bsonOrb []byte) {

	bsonOrb, err := bson.Marshal(&Orb{
		X: 1.0,
		Y: 2.0,
		ID: "4aa4ec0e-cf3f-4f4d-827f-f1415b51d26d",
	})
	if err != nil {
		t.Error("Error marshalling the test data: ", err)
	} else {
		t.Log("Marshalled test data: ", bsonOrb)
	}

	return
}
