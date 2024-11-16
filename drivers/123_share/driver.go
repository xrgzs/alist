package _123Share

import (
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type Pan123Share struct {
	model.Storage
	Addition
	apiRateLimit sync.Map
	params       Params
}

func (d *Pan123Share) Config() driver.Config {
	return config
}

func (d *Pan123Share) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Pan123Share) Init(ctx context.Context) error {
	// TODO  refresh token
	// 拼接UserAgent
	if d.PlatformType == "android" {
		d.params.UserAgent = AndroidUserAgentPrefix + "(" + d.OsVersion + ";" + d.DeviceName + " " + d.DeiveType + ")"
		d.params.Platform = AndroidPlatformParam
		d.params.AppVersion = AndroidAppVer
		d.params.XChannel = AndroidXChannel
		d.params.XAppVersion = AndroidXAppVer

	} else if d.PlatformType == "tv" {
		d.params.UserAgent = TVUserAgentPrefix + "(" + d.OsVersion + ";" + d.DeviceName + " " + d.DeiveType + ")"
		d.params.Platform = TVPlatformParam
		d.params.AppVersion = TVAndroidAppVer
	}

	d.params.OsVersion = d.OsVersion
	d.params.LoginUuid = d.LoginUuid
	d.params.DeviceName = d.DeviceName
	d.params.DeviceType = d.DeiveType

	_, err := d.request(UserInfo, http.MethodGet, nil, nil)
	return err
}

func (d *Pan123Share) Drop(ctx context.Context) error {
	_, _ = d.request(Logout, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{})
	}, nil)
	return nil
}

func (d *Pan123Share) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	// TODO return the files list, required
	files, err := d.getFiles(ctx, dir.GetID())
	if err != nil {
		return nil, err
	}
	return utils.SliceConvert(files, func(src File) (model.Obj, error) {
		return src, nil
	})
}

func (d *Pan123Share) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	// TODO return link of file, required
	if f, ok := file.(File); ok {
		//var resp DownResp
		data := base.Json{
			"driveId":   "0",
			"shareKey":  d.ShareKey,
			"SharePwd":  d.SharePwd,
			"etag":      f.Etag,
			"fileId":    f.FileId,
			"s3keyFlag": f.S3KeyFlag,
			"FileName":  f.FileName,
			"size":      f.Size,
		}
		resp, err := d.request(DownloadInfo, http.MethodPost, func(req *resty.Request) {
			req.SetBody(data)
		}, nil)
		if err != nil {
			return nil, err
		}
		downloadUrl := utils.Json.Get(resp, "data", "DownloadURL").ToString()
		u, err := url.Parse(downloadUrl)
		if err != nil {
			return nil, err
		}
		nu := u.Query().Get("params")
		if nu != "" {
			du, _ := base64.StdEncoding.DecodeString(nu)
			u, err = url.Parse(string(du))
			if err != nil {
				return nil, err
			}
		}
		u_ := u.String()
		log.Debug("download url: ", u_)
		res, err := base.NoRedirectClient.R().SetHeader("Referer", "https://www.123pan.com/").Get(u_)
		if err != nil {
			return nil, err
		}
		log.Debug(res.String())
		link := model.Link{
			URL: u_,
		}
		log.Debugln("res code: ", res.StatusCode())
		if res.StatusCode() == 302 {
			link.URL = res.Header().Get("location")
		} else if res.StatusCode() < 300 {
			link.URL = utils.Json.Get(res.Body(), "data", "redirect_url").ToString()
		}
		link.Header = http.Header{
			"Referer": []string{"https://www.123pan.com/"},
		}
		return &link, nil
	}
	return nil, fmt.Errorf("can't convert obj")
}

func (d *Pan123Share) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	// TODO create folder, optional
	return errs.NotSupport
}

func (d *Pan123Share) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	// TODO move obj, optional
	return errs.NotSupport
}

func (d *Pan123Share) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	// TODO rename obj, optional
	return errs.NotSupport
}

func (d *Pan123Share) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	// TODO copy obj, optional
	return errs.NotSupport
}

func (d *Pan123Share) Remove(ctx context.Context, obj model.Obj) error {
	// TODO remove obj, optional
	return errs.NotSupport
}

func (d *Pan123Share) Put(ctx context.Context, dstDir model.Obj, stream model.FileStreamer, up driver.UpdateProgress) error {
	// TODO upload file, optional
	return errs.NotSupport
}

//func (d *Pan123Share) Other(ctx context.Context, args model.OtherArgs) (interface{}, error) {
//	return nil, errs.NotSupport
//}

func (d *Pan123Share) APIRateLimit(ctx context.Context, api string) error {
	value, _ := d.apiRateLimit.LoadOrStore(api,
		rate.NewLimiter(rate.Every(800*time.Millisecond), 1))
	limiter := value.(*rate.Limiter)

	return limiter.Wait(ctx)
}

var _ driver.Driver = (*Pan123Share)(nil)
