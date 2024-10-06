package _123Share

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
)

const (
	Api          = "https://www.123pan.com/api"
	AApi         = "https://www.123pan.com/a/api"
	BApi         = "https://www.123pan.com/b/api"
	MainApi      = Api
	FileList     = MainApi + "/share/get"
	DownloadInfo = MainApi + "/share/download/info"
	//AuthKeySalt      = "8-8D$sL8gPjom7bk#cY"
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

func (d *Pan123Share) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
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
	//res, err := req.Execute(method, GetApi(url))
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	body := res.Body()
	code := utils.Json.Get(body, "code").ToInt()
	if code != 0 {
		return nil, errors.New(jsoniter.Get(body, "message").ToString())
	}
	return body, nil
}

func (d *Pan123Share) getFiles(ctx context.Context, parentId string) ([]File, error) {
	page := 1
	res := make([]File, 0)
	for {
		if err := d.APIRateLimit(ctx, FileList); err != nil {
			return nil, err
		}
		var resp Files
		query := map[string]string{
			"limit":          "100",
			"next":           "0",
			"orderBy":        "file_id",
			"orderDirection": "desc",
			"parentFileId":   parentId,
			"Page":           strconv.Itoa(page),
			"shareKey":       d.ShareKey,
			"SharePwd":       d.SharePwd,
		}
		_, err := d.request(FileList, http.MethodGet, func(req *resty.Request) {
			req.SetQueryParams(query)
		}, &resp)
		if err != nil {
			return nil, err
		}
		page++
		res = append(res, resp.Data.InfoList...)
		if len(resp.Data.InfoList) == 0 || resp.Data.Next == "-1" {
			break
		}
	}
	return res, nil
}

// do others that not defined in Driver interface

func generateAuthKey() string {
	timestamp := time.Now().Unix()
	randomInt := rand.Intn(1e9)                                     // 生成9位的随机整数
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")     // 去掉 UUID 中的所有 -
	return fmt.Sprintf("%d-%09d-%s", timestamp, randomInt, uuidStr) // 确保随机整数是9位
}
