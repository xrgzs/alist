package qingque

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	// Usually one of two
	driver.RootID
	// define other
	Cookie string `json:"cookie" type:"string" required:"true"`
}

var config = driver.Config{
	Name:              "Qingque",
	LocalSort:         true, // TODO: support cloud sort
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          true, // TODO: support upload
	NeedMs:            false,
	DefaultRoot:       "",
	CheckStatus:       false,
	Alert:             "",
	NoOverwriteUpload: false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Qingque{}
	})
}
