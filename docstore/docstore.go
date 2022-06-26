package docstore

import "github.com/silverswords/devcloud/docstore/driver"

type DriverMap struct {
	m map[string]interface{}
}

var driverMap = &DriverMap{}

func DefaultDriverMap() *DriverMap {
	return driverMap
}

func (d *DriverMap) Register(name string, db CollectionURLOpener) {
	// todo(abserari): add lock
	d.m[name] = db
}

type Docstore struct {
	driver driver.DB
}

var NewDocstore = newDocstore

func newDocstore(driver driver.DB) *Docstore {
	ds := &Docstore{
		driver: driver,
	}
	return ds
}
