/*******************************************************************************
 * Copyright 2018 Cavium
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

package memory

import (
	"errors"
	"fmt"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	contract "github.com/edgexfoundry/edgex-go/pkg/models"
	"github.com/globalsign/mgo/bson"
)

// Schedule event
func (m *MemDB) GetAllScheduleEvents(se *[]contract.ScheduleEvent) error {
	cpy := make([]contract.ScheduleEvent, len(m.scheduleEvents))
	copy(cpy, m.scheduleEvents)
	*se = cpy
	return nil
}

func (m *MemDB) AddScheduleEvent(se *contract.ScheduleEvent) error {
	currentTime := db.MakeTimestamp()
	se.Created = currentTime
	se.Modified = currentTime
	se.Id = bson.NewObjectId().Hex()

	for _, s := range m.scheduleEvents {
		if s.Name == se.Name {
			return db.ErrNotUnique
		}
	}

	validAddressable := false
	// Test addressable id or name exists
	for _, a := range m.addressables {
		if a.Name == se.Addressable.Name {
			validAddressable = true
			break
		}
		if a.Id == se.Addressable.Id {
			validAddressable = true
			break
		}
	}

	if !validAddressable {
		return errors.New("Invalid addressable")
	}

	m.scheduleEvents = append(m.scheduleEvents, *se)
	return nil
}

func (m *MemDB) GetScheduleEventByName(se *contract.ScheduleEvent, n string) error {
	for _, s := range m.scheduleEvents {
		if s.Name == n {
			a, err := m.GetAddressableById(s.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					se.Addressable.Id, se.Id)
			}
			s.Addressable = a
			*se = s
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) UpdateScheduleEvent(se contract.ScheduleEvent) error {
	for i, s := range m.scheduleEvents {
		if s.Id == se.Id {
			m.scheduleEvents[i] = se
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetScheduleEventById(se *contract.ScheduleEvent, id string) error {
	for _, s := range m.scheduleEvents {
		if s.Id == id {
			a, err := m.GetAddressableById(s.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					se.Addressable.Id, se.Id)
			}
			s.Addressable = a
			*se = s
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetScheduleEventsByScheduleName(ses *[]contract.ScheduleEvent, n string) error {
	l := []contract.ScheduleEvent{}
	for _, se := range m.scheduleEvents {
		if se.Schedule == n {
			a, err := m.GetAddressableById(se.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for se %s",
					se.Addressable.Id, se.Id)
			}
			se.Addressable = a
			l = append(l, se)
		}
	}
	*ses = l
	return nil
}

func (m *MemDB) GetScheduleEventsByAddressableId(ses *[]contract.ScheduleEvent, id string) error {
	l := []contract.ScheduleEvent{}
	for _, se := range m.scheduleEvents {
		if se.Addressable.Id == id {
			a, err := m.GetAddressableById(se.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for se %s",
					se.Addressable.Id, se.Id)
			}
			se.Addressable = a
			l = append(l, se)
		}
	}
	*ses = l
	return nil
}

func (m *MemDB) GetScheduleEventsByServiceName(ses *[]contract.ScheduleEvent, n string) error {
	l := []contract.ScheduleEvent{}
	for _, se := range m.scheduleEvents {
		if se.Service == n {
			a, err := m.GetAddressableById(se.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for se %s",
					se.Addressable.Id, se.Id)
			}
			se.Addressable = a
			l = append(l, se)
		}
	}
	*ses = l
	return nil
}

func (m *MemDB) DeleteScheduleEventById(id string) error {
	for i, s := range m.scheduleEvents {
		if s.Id == id {
			m.scheduleEvents = append(m.scheduleEvents[:i], m.scheduleEvents[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

// Schedule
func (m *MemDB) GetAllSchedules(s *[]contract.Schedule) error {
	cpy := make([]contract.Schedule, len(m.schedules))
	copy(cpy, m.schedules)
	*s = cpy
	return nil
}

func (m *MemDB) AddSchedule(s *contract.Schedule) error {
	currentTime := db.MakeTimestamp()
	s.Created = currentTime
	s.Modified = currentTime
	s.Id = bson.NewObjectId().Hex()

	for _, ss := range m.schedules {
		if ss.Name == s.Name {
			return db.ErrNotUnique
		}
	}

	m.schedules = append(m.schedules, *s)
	return nil
}

func (m *MemDB) GetScheduleByName(s *contract.Schedule, n string) error {
	for _, ss := range m.schedules {
		if ss.Name == n {
			*s = ss
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) UpdateSchedule(s contract.Schedule) error {
	s.Modified = db.MakeTimestamp()
	for i, ss := range m.schedules {
		if ss.Id == s.Id {
			m.schedules[i] = s
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MemDB) GetScheduleById(s *contract.Schedule, id string) error {
	for _, ss := range m.schedules {
		if ss.Id == id {
			*s = ss
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) DeleteScheduleById(id string) error {
	for i, ss := range m.schedules {
		if ss.Id == id {
			m.schedules = append(m.schedules[:i], m.schedules[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

// Device Report
func (m *MemDB) GetAllDeviceReports(drs *[]contract.DeviceReport) error {
	cpy := make([]contract.DeviceReport, len(m.deviceReports))
	copy(cpy, m.deviceReports)
	*drs = cpy
	return nil
}

func (m *MemDB) GetDeviceReportByDeviceName(drs *[]contract.DeviceReport, n string) error {
	l := []contract.DeviceReport{}
	for _, dr := range m.deviceReports {
		if dr.Name == n {
			l = append(l, dr)
		}
	}
	*drs = l
	return nil
}

func (m *MemDB) GetDeviceReportByName(dr *contract.DeviceReport, n string) error {
	for _, d := range m.deviceReports {
		if d.Name == n {
			*dr = d
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceReportById(dr *contract.DeviceReport, id string) error {
	for _, d := range m.deviceReports {
		if d.Id == id {
			*dr = d
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) AddDeviceReport(dr *contract.DeviceReport) error {
	currentTime := db.MakeTimestamp()
	dr.Created = currentTime
	dr.Modified = currentTime
	dr.Id = bson.NewObjectId().Hex()

	dummy := contract.DeviceReport{}
	if m.GetDeviceReportByName(&dummy, dr.Name) == nil {
		return db.ErrNotUnique
	}

	m.deviceReports = append(m.deviceReports, *dr)
	return nil

}

func (m *MemDB) UpdateDeviceReport(dr *contract.DeviceReport) error {
	for i, d := range m.deviceReports {
		if d.Id == dr.Id {
			m.deviceReports[i] = *dr
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceReportsByScheduleEventName(drs *[]contract.DeviceReport, n string) error {
	l := []contract.DeviceReport{}
	for _, dr := range m.deviceReports {
		if dr.Event == n {
			l = append(l, dr)
		}
	}
	*drs = l
	return nil
}

func (m *MemDB) DeleteDeviceReportById(id string) error {
	for i, c := range m.deviceReports {
		if c.Id == id {
			m.deviceReports = append(m.deviceReports[:i], m.deviceReports[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

// Device
func (m *MemDB) updateDeviceValues(d *contract.Device) error {
	a, err := m.GetAddressableById(d.Addressable.Id)
	if err != nil {
		return fmt.Errorf("Could not find addressable %s for ds %s",
			d.Addressable.Id, d.Id)
	}
	d.Addressable = a

	err = m.GetDeviceServiceById(&d.Service, d.Service.Id)
	if err != nil {
		return fmt.Errorf("Could not find DeviceService %s for ds %s",
			d.Service.Id, d.Id)
	}
	err = m.GetDeviceProfileById(&d.Profile, d.Profile.Id)
	if err != nil {
		return fmt.Errorf("Could not find DeviceProfile %s for ds %s",
			d.Profile.Id, d.Id)
	}
	return nil
}

type deviceCmp func(contract.Device) bool

func (m *MemDB) getDeviceBy(d *contract.Device, f deviceCmp) error {
	for _, dd := range m.devices {
		if f(dd) {
			if err := m.updateDeviceValues(&dd); err != nil {
				return err
			}
			*d = dd
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) getDevicesBy(d *[]contract.Device, f deviceCmp) error {
	l := []contract.Device{}
	for _, dd := range m.devices {
		if f(dd) {
			if err := m.updateDeviceValues(&dd); err != nil {
				return err
			}
			l = append(l, dd)
		}
	}
	*d = l
	return nil
}

func (m *MemDB) UpdateDevice(d contract.Device) error {
	for i, dd := range m.devices {
		if dd.Id == d.Id {
			m.devices[i] = d
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceById(d *contract.Device, id string) error {
	return m.getDeviceBy(d,
		func(dd contract.Device) bool {
			return dd.Id == id
		})
}

func (m *MemDB) GetDeviceByName(d *contract.Device, n string) error {
	return m.getDeviceBy(d,
		func(dd contract.Device) bool {
			return dd.Name == n
		})
}

func (m *MemDB) GetAllDevices(d *[]contract.Device) error {
	cpy := make([]contract.Device, len(m.devices))
	copy(cpy, m.devices)
	*d = cpy
	return nil
}

func (m *MemDB) GetDevicesByProfileId(d *[]contract.Device, id string) error {
	return m.getDevicesBy(d,
		func(dd contract.Device) bool {
			return dd.Profile.Id == id
		})
}

func (m *MemDB) GetDevicesByServiceId(d *[]contract.Device, id string) error {
	return m.getDevicesBy(d,
		func(dd contract.Device) bool {
			return dd.Service.Id == id
		})
}

func (m *MemDB) GetDevicesByAddressableId(d *[]contract.Device, id string) error {
	return m.getDevicesBy(d,
		func(dd contract.Device) bool {
			return dd.Addressable.Id == id
		})
}

func (m *MemDB) GetDevicesWithLabel(d *[]contract.Device, l string) error {
	return m.getDevicesBy(d,
		func(dd contract.Device) bool {
			return stringInSlice(l, dd.Labels)
		})
}

func (m *MemDB) AddDevice(d *contract.Device) error {
	currentTime := db.MakeTimestamp()
	d.Created = currentTime
	d.Modified = currentTime
	d.Id = bson.NewObjectId().Hex()

	for _, dd := range m.devices {
		if dd.Name == d.Name {
			return db.ErrNotUnique
		}
	}

	validAddressable := false
	// Test addressable id or name exists
	for _, a := range m.addressables {
		if a.Name == d.Addressable.Name {
			validAddressable = true
			break
		}
		if a.Id == d.Addressable.Id {
			validAddressable = true
			break
		}
	}

	if !validAddressable {
		return errors.New("Invalid addressable")
	}

	m.devices = append(m.devices, *d)
	return nil
}

func (m *MemDB) DeleteDeviceById(id string) error {
	for i, dd := range m.devices {
		if dd.Id == id {
			m.devices = append(m.devices[:i], m.devices[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) UpdateDeviceProfile(dp *contract.DeviceProfile) error {
	for i, d := range m.deviceProfiles {
		if d.Id == dp.Id {
			m.deviceProfiles[i] = *dp
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) AddDeviceProfile(d *contract.DeviceProfile) error {
	currentTime := db.MakeTimestamp()
	d.Created = currentTime
	d.Modified = currentTime
	d.Id = bson.NewObjectId().Hex()

	for _, dd := range m.deviceProfiles {
		if dd.Name == d.Name {
			return db.ErrNotUnique
		}
	}

	m.deviceProfiles = append(m.deviceProfiles, *d)
	return nil
}

func (m *MemDB) GetAllDeviceProfiles(d *[]contract.DeviceProfile) error {
	cpy := make([]contract.DeviceProfile, len(m.deviceProfiles))
	copy(cpy, m.deviceProfiles)
	*d = cpy
	return nil
}

func (m *MemDB) GetDeviceProfileById(d *contract.DeviceProfile, id string) error {
	for _, dp := range m.deviceProfiles {
		if dp.Id == id {
			*d = dp
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) DeleteDeviceProfileById(id string) error {
	for i, d := range m.deviceProfiles {
		if d.Id == id {
			m.deviceProfiles = append(m.deviceProfiles[:i], m.deviceProfiles[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceProfilesByModel(dps *[]contract.DeviceProfile, model string) error {
	l := []contract.DeviceProfile{}
	for _, dp := range m.deviceProfiles {
		if dp.Model == model {
			l = append(l, dp)
		}
	}
	*dps = l
	return nil
}

func (m *MemDB) GetDeviceProfilesWithLabel(dps *[]contract.DeviceProfile, label string) error {
	l := []contract.DeviceProfile{}
	for _, dp := range m.deviceProfiles {
		if stringInSlice(label, dp.Labels) {
			l = append(l, dp)
		}
	}
	*dps = l
	return nil
}

func (m *MemDB) GetDeviceProfilesByManufacturerModel(dps *[]contract.DeviceProfile, man string, mod string) error {
	l := []contract.DeviceProfile{}
	for _, dp := range m.deviceProfiles {
		if dp.Manufacturer == man && dp.Model == mod {
			l = append(l, dp)
		}
	}
	*dps = l
	return nil
}

func (m *MemDB) GetDeviceProfilesByManufacturer(dps *[]contract.DeviceProfile, man string) error {
	l := []contract.DeviceProfile{}
	for _, dp := range m.deviceProfiles {
		if dp.Manufacturer == man {
			l = append(l, dp)
		}
	}
	*dps = l
	return nil
}

func (m *MemDB) GetDeviceProfileByName(d *contract.DeviceProfile, n string) error {
	for _, dp := range m.deviceProfiles {
		if dp.Name == n {
			*d = dp
			return nil
		}
	}
	return db.ErrNotFound
}

// Addressable
func (m *MemDB) UpdateAddressable(orig contract.Addressable) error {

	for i, aa := range m.addressables {
		if aa.Id == orig.Id {
			m.addressables[i] = orig
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MemDB) AddAddressable(a contract.Addressable) (string, error) {
	currentTime := db.MakeTimestamp()
	a.Created = currentTime
	a.Modified = currentTime

	for _, aa := range m.addressables {
		if aa.Name == a.Name {
			return a.Id, db.ErrNotUnique
		}
	}

	m.addressables = append(m.addressables, a)
	return a.Id, nil
}

type addressableCmp func(contract.Addressable) bool

func (m *MemDB) getAddressableBy(f addressableCmp) (contract.Addressable, error) {
	for _, aa := range m.addressables {
		if f(aa) {
			return aa, nil
		}
	}
	return contract.Addressable{}, db.ErrNotFound
}

func (m *MemDB) getAddressablesBy(f addressableCmp) []contract.Addressable {
	l := []contract.Addressable{}
	for _, aa := range m.addressables {
		if f(aa) {
			l = append(l, aa)
		}
	}
	return l
}

func (m *MemDB) GetAddressableById(id string) (contract.Addressable, error) {
	return m.getAddressableBy(
		func(aa contract.Addressable) bool {
			return aa.Id == id
		})
}

func (m *MemDB) GetAddressableByName(n string) (contract.Addressable, error) {
	return m.getAddressableBy(
		func(aa contract.Addressable) bool {
			return aa.Name == n
		})
}

func (m *MemDB) GetAddressablesByTopic(t string) ([]contract.Addressable, error) {
	return m.getAddressablesBy(
		func(aa contract.Addressable) bool {
			return aa.Topic == t
		}), nil
}

func (m *MemDB) GetAddressablesByPort(p int) ([]contract.Addressable, error) {
	return m.getAddressablesBy(
		func(aa contract.Addressable) bool {
			return aa.Port == p
		}), nil
}

func (m *MemDB) GetAddressablesByPublisher(p string) ([]contract.Addressable, error) {
	return m.getAddressablesBy(
		func(aa contract.Addressable) bool {
			return aa.Publisher == p
		}), nil
}

func (m *MemDB) GetAddressablesByAddress(add string) ([]contract.Addressable, error) {
	return m.getAddressablesBy(
		func(aa contract.Addressable) bool {
			return aa.Address == add
		}), nil
}

func (m *MemDB) GetAddressables() ([]contract.Addressable, error) {
	return m.addressables, nil
}

func (m *MemDB) DeleteAddressableById(id string) error {
	var found bool
	list := []contract.Addressable{}
	for i, aa := range m.addressables {
		if aa.Id != id {
			list = append(list, m.addressables[i])
		} else {
			found = true
		}
	}
	if !found {
		return db.ErrNotFound
	}
	m.addressables = list
	return nil
}

// Device service
func (m *MemDB) UpdateDeviceService(ds contract.DeviceService) error {
	for i, d := range m.deviceServices {
		if d.Id == ds.Id {
			m.deviceServices[i] = ds
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceServicesByAddressableId(d *[]contract.DeviceService, id string) error {
	l := []contract.DeviceService{}
	for _, ds := range m.deviceServices {
		if ds.Addressable.Id == id {
			_, err := m.GetAddressableById(ds.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					ds.Addressable.Id, ds.Id)
			}
			l = append(l, ds)
		}
	}
	*d = l
	return nil
}

func (m *MemDB) GetDeviceServicesWithLabel(d *[]contract.DeviceService, label string) error {
	l := []contract.DeviceService{}
	for _, ds := range m.deviceServices {
		if stringInSlice(label, ds.Labels) {
			_, err := m.GetAddressableById(ds.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					ds.Addressable.Id, ds.Id)
			}
			l = append(l, ds)
		}
	}
	*d = l
	return nil
}

func (m *MemDB) GetDeviceServiceById(d *contract.DeviceService, id string) error {
	for _, ds := range m.deviceServices {
		if ds.Id == id {
			_, err := m.GetAddressableById(ds.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					ds.Addressable.Id, ds.Id)
			}
			*d = ds
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceServiceByName(d *contract.DeviceService, n string) error {
	for _, ds := range m.deviceServices {
		if ds.Name == n {
			_, err := m.GetAddressableById(ds.Addressable.Id)
			if err != nil {
				return fmt.Errorf("Could not find addressable %s for ds %s",
					ds.Addressable.Id, ds.Id)
			}
			*d = ds
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetAllDeviceServices(d *[]contract.DeviceService) error {
	for _, ds := range m.deviceServices {
		_, err := m.GetAddressableById(ds.Addressable.Id)
		if err != nil {
			return fmt.Errorf("Could not find addressable %s for ds %s",
				ds.Addressable.Id, ds.Id)
		}
	}
	cpy := make([]contract.DeviceService, len(m.deviceServices))
	copy(cpy, m.deviceServices)
	*d = cpy
	return nil
}

func (m *MemDB) AddDeviceService(ds *contract.DeviceService) error {
	currentTime := db.MakeTimestamp()
	ds.Created = currentTime
	ds.Modified = currentTime
	ds.Id = bson.NewObjectId().Hex()

	for _, d := range m.deviceServices {
		if d.Name == ds.Name {
			return db.ErrNotUnique
		}
	}

	validAddressable := false
	// Test addressable id or name exists
	for _, a := range m.addressables {
		if a.Name == ds.Addressable.Name {
			validAddressable = true
			break
		}
		if a.Id == ds.Addressable.Id {
			validAddressable = true
			break
		}
	}

	if !validAddressable {
		return errors.New("Invalid addressable")
	}

	m.deviceServices = append(m.deviceServices, *ds)
	return nil
}

func (m *MemDB) DeleteDeviceServiceById(id string) error {
	for i, d := range m.deviceServices {
		if d.Id == id {
			m.deviceServices = append(m.deviceServices[:i], m.deviceServices[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

// Provision watcher
type provisionWatcherComp func(contract.ProvisionWatcher) bool

func (m *MemDB) getProvisionWatcherBy(pw *contract.ProvisionWatcher, f provisionWatcherComp) error {
	for _, p := range m.provisionWatchers {
		if f(p) {
			err := m.GetDeviceServiceById(&p.Service, p.Service.Id)
			if err != nil {
				return fmt.Errorf("Could not find DeviceService %s for ds %s",
					p.Service.Id, p.Id.Hex())
			}
			err = m.GetDeviceProfileById(&p.Profile, p.Profile.Id)
			if err != nil {
				return fmt.Errorf("Could not find DeviceProfile %s for ds %s",
					p.Profile.Id, p.Id.Hex())
			}
			*pw = p
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) getProvisionWatchersBy(pws *[]contract.ProvisionWatcher, f provisionWatcherComp) error {
	l := []contract.ProvisionWatcher{}
	for _, pw := range m.provisionWatchers {
		if f(pw) {
			err := m.GetDeviceServiceById(&pw.Service, pw.Service.Id)
			if err != nil {
				return fmt.Errorf("Could not find DeviceService %s for ds %s",
					pw.Service.Id, pw.Id.Hex())
			}
			err = m.GetDeviceProfileById(&pw.Profile, pw.Profile.Id)
			if err != nil {
				return fmt.Errorf("Could not find DeviceProfile %s for ds %s",
					pw.Profile.Id, pw.Id.Hex())
			}
			l = append(l, pw)
		}
	}
	*pws = l
	return nil
}

func (m *MemDB) GetProvisionWatcherById(pw *contract.ProvisionWatcher, id string) error {
	return m.getProvisionWatcherBy(pw,
		func(p contract.ProvisionWatcher) bool {
			return p.Id.Hex() == id
		})
}

func (m *MemDB) GetAllProvisionWatchers(pw *[]contract.ProvisionWatcher) error {
	*pw = m.provisionWatchers
	return nil
}

func (m *MemDB) GetProvisionWatcherByName(pw *contract.ProvisionWatcher, n string) error {
	return m.getProvisionWatcherBy(pw,
		func(p contract.ProvisionWatcher) bool {
			return p.Name == n
		})
}

func (m *MemDB) GetProvisionWatchersByProfileId(pw *[]contract.ProvisionWatcher, id string) error {
	return m.getProvisionWatchersBy(pw,
		func(p contract.ProvisionWatcher) bool {
			return p.Profile.Id == id
		})
}

func (m *MemDB) GetProvisionWatchersByServiceId(pw *[]contract.ProvisionWatcher, id string) error {
	return m.getProvisionWatchersBy(pw,
		func(p contract.ProvisionWatcher) bool {
			return p.Service.Id == id
		})
}

func (m *MemDB) GetProvisionWatchersByIdentifier(pw *[]contract.ProvisionWatcher, k string, v string) error {
	return m.getProvisionWatchersBy(pw,
		func(p contract.ProvisionWatcher) bool {
			return p.Identifiers[k] == v
		})
}

func (m *MemDB) updateProvisionWatcherValues(pw *contract.ProvisionWatcher) error {
	// get Device Service
	validDeviceService := false
	var dev contract.DeviceService
	var err error
	if pw.Service.Id != "" {
		if err = m.GetDeviceServiceById(&dev, pw.Service.Id); err == nil {
			validDeviceService = true
		}
	} else if pw.Service.Name != "" {
		if err = m.GetDeviceServiceByName(&dev, pw.Service.Name); err == nil {
			validDeviceService = true
		}
	} else {
		return errors.New("Device Service ID or Name is required")
	}
	if !validDeviceService {
		return fmt.Errorf("Invalid DeviceService: %v", err)
	}
	pw.Service = dev

	// get Device Profile
	validDeviceProfile := false
	var dp contract.DeviceProfile
	if pw.Profile.Id != "" {
		if err = m.GetDeviceProfileById(&dp, pw.Profile.Id); err == nil {
			validDeviceProfile = true
		}
	} else if pw.Profile.Name != "" {
		if err = m.GetDeviceProfileByName(&dp, pw.Profile.Name); err == nil {
			validDeviceProfile = true
		}
	} else {
		return errors.New("Device Profile ID or Name is required")
	}
	if !validDeviceProfile {
		return fmt.Errorf("Invalid DeviceProfile: %v", err)
	}
	pw.Profile = dp
	return nil
}

func (m *MemDB) AddProvisionWatcher(pw *contract.ProvisionWatcher) error {
	currentTime := db.MakeTimestamp()
	pw.Created = currentTime
	pw.Modified = currentTime
	pw.Id = bson.NewObjectId()

	p := contract.ProvisionWatcher{}
	if err := m.GetProvisionWatcherByName(&p, pw.Name); err == nil {
		return db.ErrNotUnique
	}

	if err := m.updateProvisionWatcherValues(pw); err != nil {
		return err
	}
	m.provisionWatchers = append(m.provisionWatchers, *pw)
	return nil
}

func (m *MemDB) UpdateProvisionWatcher(pw contract.ProvisionWatcher) error {
	pw.Modified = db.MakeTimestamp()

	if err := m.updateProvisionWatcherValues(&pw); err != nil {
		return err
	}
	for i, p := range m.provisionWatchers {
		if pw.Id == p.Id {
			m.provisionWatchers[i] = p
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) DeleteProvisionWatcherById(id string) error {
	for i, p := range m.provisionWatchers {
		if p.Id.Hex() == id {
			m.provisionWatchers = append(m.provisionWatchers[:i], m.provisionWatchers[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

// Command
func (m *MemDB) GetCommandById(c *contract.Command, id string) error {
	for _, cc := range m.commands {
		if cc.Id == id {
			*c = cc
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetCommandByName(d *[]contract.Command, name string) error {
	cmds := []contract.Command{}
	for _, cc := range m.commands {
		if cc.Name == name {
			cmds = append(cmds, cc)
		}
	}
	*d = cmds
	return nil
}

func (m *MemDB) AddCommand(c *contract.Command) error {
	currentTime := db.MakeTimestamp()
	c.Created = currentTime
	c.Modified = currentTime
	c.Id = bson.NewObjectId().Hex()

	m.commands = append(m.commands, *c)
	return nil
}

func (m *MemDB) GetAllCommands(d *[]contract.Command) error {
	cpy := make([]contract.Command, len(m.commands))
	copy(cpy, m.commands)
	*d = cpy
	return nil
}

func (m *MemDB) UpdateCommand(updated *contract.Command, orig *contract.Command) error {
	if updated == nil {
		return nil
	}
	if updated.Name != "" {
		orig.Name = updated.Name
	}
	if updated.Get != nil && (updated.Get.String() != contract.Get{}.String()) {
		orig.Get = updated.Get
	}
	if updated.Put != nil && (updated.Put.String() != contract.Put{}.String()) {
		orig.Put = updated.Put
	}
	if updated.Origin != 0 {
		orig.Origin = updated.Origin
	}

	for i, c := range m.commands {
		if c.Id == orig.Id {
			m.commands[i] = *orig
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MemDB) DeleteCommandById(id string) error {
	for i, c := range m.commands {
		if c.Id == id {
			m.commands = append(m.commands[:i], m.commands[i+1:]...)
			return nil
		}
	}
	return db.ErrNotFound
}

func (m *MemDB) GetDeviceProfilesUsingCommand(dps *[]contract.DeviceProfile, c contract.Command) error {
	l := []contract.DeviceProfile{}
	for _, dp := range m.deviceProfiles {
		for _, cc := range dp.Commands {
			if cc.Id == c.Id {
				l = append(l, dp)
				break
			}
		}
	}
	*dps = l
	return nil
}

func (m *MemDB) ScrubMetadata() error {
	m.addressables = nil
	m.commands = nil
	m.deviceServices = nil
	m.schedules = nil
	m.scheduleEvents = nil
	m.provisionWatchers = nil
	m.deviceReports = nil
	m.deviceProfiles = nil
	m.devices = nil
	return nil
}
