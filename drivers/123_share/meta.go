package _123Share

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	ShareKey string `json:"sharekey" required:"true"`
	SharePwd string `json:"sharepassword"`
	driver.RootID
	//OrderBy        string `json:"order_by" type:"select" options:"file_name,size,update_at" default:"file_name"`
	//OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" default:"asc"`
	AccessToken  string `json:"accesstoken" type:"text" required:"true"`
	PlatformType string `json:"platformType" type:"select" options:"android,tv" default:"android" required:"true"`
	DeviceName   string `json:"devicename" default:"XiaoMi"`
	DeiveType    string `json:"devicetype" default:"houji"`
	OsVersion    string `json:"osversion" default:"14"`
	LoginUuid    string `json:"loginuuid" default:"1fce20b2428d30899fd537f4cf231dfb"`
}

var config = driver.Config{
	Name:              "123PanShare",
	LocalSort:         true,
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          true,
	NeedMs:            false,
	DefaultRoot:       "0",
	CheckStatus:       false,
	Alert:             "",
	NoOverwriteUpload: false,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Pan123Share{}
	})
}
