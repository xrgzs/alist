package _123

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"hash/crc32"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

// do others that not defined in Driver interface

const (
	Api              = "https://www.123pan.com/api"
	AApi             = "https://www.123pan.com/a/api"
	BApi             = "https://www.123pan.com/b/api"
	MainApi          = Api
	SignIn           = MainApi + "/user/sign_in"
	Logout           = MainApi + "/user/logout"
	UserInfo         = MainApi + "/user/info"
	FileList         = MainApi + "/file/list/new"
	DownloadInfo     = MainApi + "/file/download_info"
	Mkdir            = MainApi + "/file/upload_request"
	Move             = MainApi + "/file/mod_pid"
	Rename           = MainApi + "/file/rename"
	Trash            = MainApi + "/file/trash"
	UploadRequest    = MainApi + "/file/upload_request"
	UploadComplete   = MainApi + "/file/upload_complete"
	S3PreSignedUrls  = MainApi + "/file/s3_repare_upload_parts_batch"
	S3Auth           = MainApi + "/file/s3_upload_object/auth"
	UploadCompleteV2 = MainApi + "/file/upload_complete/v2"
	S3Complete       = MainApi + "/file/s3_complete_multipart_upload"
	//AuthKeySalt      = "8-8D$sL8gPjom7bk#cY"
	QrcodeGenerate = MainApi + "/user/qr-code/generate"
	QrcodeResult   = MainApi + "/user/qr-code/result"
)

const (
	AndroidUserAgentPrefix = "123pan/v2.4.7" // 123pan/v2.4.7(Android_14;XiaoMi)
	AndroidPlatformParam   = "android"
	AndroidAppVer          = "69"
	AndroidXAppVer         = "2.4.7"
	AndroidXChannel        = "1002"
	TVUserAgentPrefix      = "123pan_android_tv/1.0.0" // 123pan_android_tv/1.0.0(14;samsung SM-X800)
	TVPlatformParam        = "android_tv"
	TVAndroidAppVer        = "100"
)

type Params struct {
	UserAgent   string
	Platform    string
	AppVersion  string
	OsVersion   string
	LoginUuid   string
	DeviceType  string
	DeviceName  string
	XChannel    string
	XAppVersion string
}

func signPath(path string, os string, version string) (k string, v string) {
	table := []byte{'a', 'd', 'e', 'f', 'g', 'h', 'l', 'm', 'y', 'i', 'j', 'n', 'o', 'p', 'k', 'q', 'r', 's', 't', 'u', 'b', 'c', 'v', 'w', 's', 'z'}
	random := fmt.Sprintf("%.f", math.Round(1e7*rand.Float64()))
	now := time.Now().In(time.FixedZone("CST", 8*3600))
	timestamp := fmt.Sprint(now.Unix())
	nowStr := []byte(now.Format("200601021504"))
	for i := 0; i < len(nowStr); i++ {
		nowStr[i] = table[nowStr[i]-48]
	}
	timeSign := fmt.Sprint(crc32.ChecksumIEEE(nowStr))
	data := strings.Join([]string{timestamp, random, path, os, version, timeSign}, "|")
	dataSign := fmt.Sprint(crc32.ChecksumIEEE([]byte(data)))
	return timeSign, strings.Join([]string{timestamp, random, dataSign}, "-")
}

func GetApi(rawUrl string) string {
	u, _ := url.Parse(rawUrl)
	query := u.Query()
	query.Add(signPath(u.Path, "web", "3"))
	u.RawQuery = query.Encode()
	return u.String()
}

func (d *Pan123) login() error {
	var body base.Json
	if utils.IsEmailFormat(d.Username) {
		body = base.Json{
			"mail":     d.Username,
			"password": d.Password,
			"type":     2,
		}
	} else {
		body = base.Json{
			"passport": d.Username,
			"password": d.Password,
			"type":     1,
		}
	}

	req := base.RestyClient.R()

	req.SetHeaders(map[string]string{
		/*			"origin":      "https://www.123pan.com",
					"referer":     "https://www.123pan.com/",*/
		"user-agent":  d.params.UserAgent,
		"platform":    d.params.Platform,
		"app-version": d.params.AppVersion,
		"osversion":   d.params.OsVersion,
		"devicetype":  d.params.DeviceType,
		"devicename":  d.params.DeviceName,
		"loginuuid":   d.params.LoginUuid,
	})

	if d.params.XChannel != "" && d.params.XAppVersion != "" {
		req.SetHeaders(map[string]string{
			"x-channel":     d.params.XChannel,
			"x-app-version": d.params.XAppVersion,
		})
	}

	req.SetQueryParam("auth-key", generateAuthKey())

	res, err := req.SetBody(body).Post(SignIn)
	//res, err := base.RestyClient.R().
	//	SetHeaders(map[string]string{
	//		/*			"origin":      "https://www.123pan.com",
	//					"referer":     "https://www.123pan.com/",*/
	//		"user-agent":  d.params.UserAgent,
	//		"platform":    d.params.Platform,
	//		"app-version": d.params.AppVersion,
	//		"osversion":   d.params.OsVersion,
	//		"devicetype":  d.params.DeviceType,
	//		"devicename":  d.params.DeviceName,
	//		//"user-agent":  base.UserAgent,
	//	}).
	//	SetBody(body).Post(SignIn)
	if err != nil {
		return err
	}
	if utils.Json.Get(res.Body(), "code").ToInt() != 200 {
		err = fmt.Errorf(utils.Json.Get(res.Body(), "message").ToString())
	} else {
		d.AccessToken = utils.Json.Get(res.Body(), "data", "token").ToString()
	}
	return err
}

//func authKey(reqUrl string) (*string, error) {
//	reqURL, err := url.Parse(reqUrl)
//	if err != nil {
//		return nil, err
//	}
//
//	nowUnix := time.Now().Unix()
//	random := rand.Intn(0x989680)
//
//	p4 := fmt.Sprintf("%d|%d|%s|%s|%s|%s", nowUnix, random, reqURL.Path, "web", "3", AuthKeySalt)
//	authKey := fmt.Sprintf("%d-%d-%x", nowUnix, random, md5.Sum([]byte(p4)))
//	return &authKey, nil
//}

func (d *Pan123) loginByQrCode() error {
	if d.Addition.UniID == "" {
		uniID, err := d.generateQrCode()
		if uniID == "" && err != nil {
			return err
		} else {
			// 保存 uniID 用于 二维码登录
			d.Addition.UniID = uniID
			op.MustSaveDriverStorage(d)
			return err
		}
	} else {
		token, err := d.getTokenByUniID()
		if token == "" && err != nil {
			return err
		} else {
			d.Addition.AccessToken = token
			op.MustSaveDriverStorage(d)
			return err
		}
	}
}

func (d *Pan123) generateQrCode() (string, error) {
	var resp QrCodeGenerateResp
	_, err := d.request(QrcodeGenerate, http.MethodGet, nil, &resp)
	if err != nil {
		return "", err
	}
	// 拼接二维码链接
	qrUrl := fmt.Sprintf(resp.Data.Url+"?uniID=%s", resp.Data.UniID+"&source=123pan&type=login")
	// 生成二维码
	qrBytes, _ := qrcode.Encode(qrUrl, qrcode.Medium, 256)
	base64Bytes := base64.StdEncoding.EncodeToString(qrBytes)
	// 展示二维码
	qrTemplate := `
	<body>
        <img src="data:image/jpeg;base64,%s"/>
		<a target="_blank" href="%s">Or Click Here</a>
    </body>`
	qrPage := fmt.Sprintf(qrTemplate, base64Bytes, qrUrl)
	return resp.Data.UniID, fmt.Errorf("need verify: \n%s", qrPage)
}

func (d *Pan123) getTokenByUniID() (string, error) {
	var resp QrCodeResultResp
	_, err := d.request(QrcodeResult, http.MethodGet, func(req *resty.Request) {
		req.SetQueryParam("uniID", d.Addition.UniID)
	}, &resp)
	if err != nil {
		return "", err
	}

	if resp.Data.LoginStatus == 4 {
		return "", errors.New("uniID expired")
	} else if resp.Data.Token == "" && resp.Data.LoginStatus == 0 {
		return "", errors.New("wait for scan qrcode")
	}

	return resp.Data.Token, nil

}

func (d *Pan123) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeaders(map[string]string{
		/*		"origin":        "https://www.123pan.com",
				"referer":       "https://www.123pan.com/",*/
		"authorization": "Bearer " + d.AccessToken,
		"user-agent":    d.params.UserAgent,
		"platform":      d.params.Platform,
		"app-version":   d.params.AppVersion,
		"osversion":     d.params.OsVersion,
		"devicetype":    d.params.DeviceType,
		"devicename":    d.params.DeviceName,
		"loginuuid":     d.params.LoginUuid,
	})

	if d.params.XChannel != "" && d.params.XAppVersion != "" {
		req.SetHeaders(map[string]string{
			"x-channel":     d.params.XChannel,
			"x-app-version": d.params.XAppVersion,
		})
	}

	req.SetQueryParam("auth-key", generateAuthKey())

	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	//authKey, err := authKey(url)
	//if err != nil {
	//	return nil, err
	//}
	//req.SetQueryParam("auth-key", *authKey)
	//res, err := req.Execute(method, GetApi(url))
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	body := res.Body()
	code := utils.Json.Get(body, "code").ToInt()
	if code != 0 && code != 200 {
		if code == 401 && d.Addition.UseQrCodeLogin == false {
			err := d.login()
			if err != nil {
				return nil, err
			}
			return d.request(url, method, callback, resp)
		} else if code == 401 && d.Addition.UseQrCodeLogin == true {
			err := d.loginByQrCode()
			if err != nil {
				return nil, err
			}
			return d.request(url, method, callback, resp)
		}
		return nil, errors.New(jsoniter.Get(body, "message").ToString())
	}
	return body, nil
}

func (d *Pan123) getFiles(ctx context.Context, parentId string, name string) ([]File, error) {
	page := 1
	total := 0
	res := make([]File, 0)
	// 2024-02-06 fix concurrency by 123pan
	for {
		if err := d.APIRateLimit(ctx, FileList); err != nil {
			return nil, err
		}
		var resp Files
		query := map[string]string{
			"driveId":              "0",
			"limit":                "100",
			"next":                 "0",
			"orderBy":              "file_id",
			"orderDirection":       "desc",
			"parentFileId":         parentId,
			"trashed":              "false",
			"SearchData":           "",
			"Page":                 strconv.Itoa(page),
			"OnlyLookAbnormalFile": "0",
			"event":                "homeListFile",
			"operateType":          "4",
			"inDirectSpace":        "false",
		}
		_res, err := d.request(FileList, http.MethodGet, func(req *resty.Request) {
			req.SetQueryParams(query)
		}, &resp)
		if err != nil {
			return nil, err
		}
		log.Debug(string(_res))
		page++
		res = append(res, resp.Data.InfoList...)
		total = resp.Data.Total
		if len(resp.Data.InfoList) == 0 || resp.Data.Next == "-1" {
			break
		}
	}
	if len(res) != total {
		log.Warnf("incorrect file count from remote at %s: expected %d, got %d", name, total, len(res))
	}
	return res, nil
}

func generateAuthKey() string {
	timestamp := time.Now().Unix()
	randomInt := rand.Intn(1e9)                                     // 生成9位的随机整数
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")     // 去掉 UUID 中的所有 -
	return fmt.Sprintf("%d-%09d-%s", timestamp, randomInt, uuidStr) // 确保随机整数是9位
}
