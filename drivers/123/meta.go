package _123

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	driver.RootID
	//OrderBy        string `json:"order_by" type:"select" options:"file_id,file_name,size,update_at" default:"file_name"`
	//OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" default:"asc"`
	AccessToken  string `json:"accesstoken" type:"text"`
	PlatformType string `json:"platformType" type:"select" options:"android,tv" default:"android" required:"true"`
	DeviceName   string `json:"devicename" default:"XiaoMi"`
	DeiveType    string `json:"devicetype" default:"houji"`
	OsVersion    string `json:"osversion" default:"14"`
	LoginUuid    string `json:"loginuuid" default:"1fce20b2428d30899fd537f4cf231dfb"`
}

var config = driver.Config{
	Name:        "123Pan",
	DefaultRoot: "0",
	LocalSort:   true,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Pan123{}
	})
}
