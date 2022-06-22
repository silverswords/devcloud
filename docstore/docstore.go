package docstore

import (
	"github.com/silverswords/devcloud/docstore/driver"
)

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
