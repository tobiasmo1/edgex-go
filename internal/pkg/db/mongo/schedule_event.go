/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package mongo

import (
 	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// Internal version of the schedule event struct
// Use this to handle DBRef
type mongoScheduleEvent struct {
	models.ScheduleEvent
}

// Custom marshaling into mongo
func (mse mongoScheduleEvent) GetBSON() (interface{}, error) {
	se := struct {
		models.BaseObject `bson:",inline"`
		Id                string `bson:"_id,omitempty"`
		Name              string        `bson:"name"`        // non-database unique identifier for a schedule event
		Schedule          string        `bson:"schedule"`    // Name to associated owning schedule
		Addressable       mgo.DBRef     `bson:"addressable"` // address {MQTT topic, HTTP address, serial bus, etc.} for the action (can be empty)
		Parameters        string        `bson:"parameters"`  // json body for parameters
		Service           string        `bson:"service"`     // json body for parameters
	}{
		Id:          mse.Id,
		Name:        mse.Name,
		Schedule:    mse.Schedule,
		Parameters:  mse.Parameters,
		Service:     mse.Service,
		Addressable: mgo.DBRef{Collection: db.Addressable, Id: mse.Addressable.Id},
	}
	// What is current value for 
	// se.Addressable.HTTPMethod = mse.Addressable.HTTPMethod
	fmt.Println("Assigned DBREF FOR ADDRESSABLE ID %v", se)

/*RESULT looks like this:
level=WARN ts=2019-02-28T08:29:19.44079761Z app=edgex-core-metadata source=rest_scheduleevent.go:110 msg="STATUS: CALLED GetAddressableById for se.Addressable.Id"
Assigning DBREF FOR ADDRESSABLE ID %v {{0 0 0} 75f3c0ea-a29b-4eed-8cb4-8551f0561a46 
turnOnSwitch 10sec-schedule 
{addressable ba271fdb-a579-4d5b-a21d-13a00466e21a } 
{"Switch": "true"} device-simple}
*/
	return se, nil
}

// Custom unmarshaling out of mongo
func (mse *mongoScheduleEvent) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		models.BaseObject `bson:",inline"`
		Id                string `bson:"_id,omitempty"`
		Name              string        `bson:"name"`        // non-database unique identifier for a schedule event
		Schedule          string        `bson:"schedule"`    // Name to associated owning schedule
		Addressable       mgo.DBRef     `bson:"addressable"` // address {MQTT topic, HTTP address, serial bus, etc.} for the action (can be empty)
		Parameters        string        `bson:"parameters"`  // json body for parameters
		Service           string        `bson:"service"`     // json body for parameters
	})

	bsonErr := raw.Unmarshal(decoded)
	if bsonErr != nil {
		return bsonErr
	}
	fmt.Println("RAW DECODED BSON [%v]", decoded)

	// Copy over the non-DBRef fields
	mse.Id = decoded.Id
	mse.Name = decoded.Name
	mse.Schedule = decoded.Schedule
	mse.Parameters = decoded.Parameters
	mse.Service = decoded.Service

	// De-reference the DBRef fields
	m, err := getCurrentMongoClient()
	if err != nil {
		return err
	}
	s := m.session.Copy()
	defer s.Close()

	addCol := s.DB(m.database.Name).C(db.Addressable)

	var a models.Addressable

// This find does not fetch the complete models.Addressable.
// It is a direct query, but in Go it is missing Addressable.HTTPMethod...
//
// Local mongo shell (same DB) works!
// db.addressable.find({uuid: {$eq:"d1294a5f-7932-441d-b7da-756f3d8cc51a"}})
	err = addCol.Find(bson.D{{Name: "uuid", Value: decoded.Addressable.Id}}).One(&a)
	if err == mgo.ErrNotFound {
		fmt.Println("FAILED TO FIND DBREF BY ADDRESSABLE UUID, TRYING _ID")
		err = addCol.Find(bson.M{"_id": decoded.Addressable.Id}).One(&a)
		if err == mgo.ErrNotFound {
			fmt.Println("FAILED TO FIND DBREF BY ADDRESSABLE _ID")
		}
	}
	if err != nil {
		return err
	}
	fmt.Println("BSON Addressable DBREF RESOLUTION (By UUID) via GLOBALSIGN/MGO [%v]", a)

	mse.Addressable = a

	return nil
}
