package fnos_share

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	// Usually one of two
	driver.RootPath
	// define other
	ShareId  string `json:"share_id" required:"false" help:"The part after the last / in the shared link"`
	SharePwd string `json:"share_pwd" required:"false" help:"The password of the shared link"`
	Host     string `json:"host" required:"true" help:"You can change it to your local area network"`
}

var config = driver.Config{
	Name:              "fnOS Share",
	LocalSort:         true,
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          true,
	NeedMs:            false,
	DefaultRoot:       "",
	CheckStatus:       false,
	Alert:             "",
	NoOverwriteUpload: false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &FnOSShare{}
	})
}
