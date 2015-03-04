
package main

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func BSONFromOrb(t *testing.T, orb Orb) (bsonOrb []byte) {

	bsonOrb, err := bson.Marshal(&orb)
	if err != nil {
		t.Error("Error marshalling the test data: ", err)
	} else {
		t.Log("Marshalled test data: ", bsonOrb)
	}

	return
}

func OrbFromBSON(t *testing.T, bsonOrb []byte) (orb Orb) {

	if err := bson.Unmarshal(bsonOrb, &orb); err != nil {
		t.Error("Could not unmarshal bson: ", err)
	} else {
		t.Log("Unmarshalled Orb: ", orb)
	}

	return
}

func AssertOrbAt(t *testing.T, orb Orb, x float32, y float32) {

	if orb.X != x {
		t.Errorf("X not at %f, got %f", x, orb.X)
	}
	if orb.Y != y {
		t.Errorf("Y not at %f, got %f", y, orb.Y)
	}
}
