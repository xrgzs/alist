package _123

import (
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/op"
)

type Addition struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	UseQrCodeLogin bool   `json:"use_qr_code_login"`
	UniID          string `json:"uni_id"`
	driver.RootID
	//OrderBy        string `json:"order_by" type:"select" options:"file_id,file_name,size,update_at" default:"file_name"`
	//OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" default:"asc"`
	AccessToken  string `json:"accesstoken" type:"text"`
	PlatformType string `json:"platformType" type:"select" options:"android,tv" default:"android" required:"true"`
	DeviceName   string `json:"devicename" default:"Xiaomi"`
	DeiveType    string `json:"devicetype" default:"M1810E5A"`
	OsVersion    string `json:"osversion" default:"Android_8.1.0"`
	LoginUuid    string `json:"loginuuid" default:""`
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
