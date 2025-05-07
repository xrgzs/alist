package qingque

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/cookie"
	"github.com/go-resty/resty/v2"
)

type Qingque struct {
	model.Storage
	Addition
	client     *resty.Client
	IdentityId string
}

func (d *Qingque) Config() driver.Config {
	return config
}

func (d *Qingque) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Qingque) Init(ctx context.Context) error {
	// TODO login / refresh token
	//op.MustSaveDriverStorage(d)
	d.client = base.NewRestyClient()
	d.client.SetCookieJar(nil)
	c := cookie.Parse(d.Cookie)
	d.client.SetCookies(c)
	if cookie.GetCookie(c, "Recent-Identity-Id") != nil {
		d.IdentityId = cookie.GetCookie(c, "Recent-Identity-Id").Value
	}
	return nil
}

func (d *Qingque) Drop(ctx context.Context) error {
	d.Cookie = cookie.ToString(d.client.Cookies)
	return nil
}

func (d *Qingque) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	var r FileResp
	var f []model.Obj
	var pageNum int64 = 1
	for {
		err := d.request(http.MethodGet, "/docs/subfolder/{docID}", func(req *resty.Request) {
			req.SetPathParam("docID", dir.GetID())
			req.SetQueryParams(map[string]string{
				"docTypeEn":    "all",
				"orderTypeEn":  "asc",
				"ownerTypeEn":  "ownerAll",
				"pageNum":      strconv.FormatInt(pageNum, 10),
				"pageSize":     "30",
				"spaceCosmoId": "mine",
				"timeTypeEn":   "title",
			})
		}, &r)
		if err != nil {
			return nil, err
		}
		for _, l := range r.CosmoExtVoPage.List {
			// filter online document
			// TODO: Implement online document
			if l.DocTypeEn == "folder" || l.DocTypeEn == "yFile" {
				f = append(f, &model.Object{
					ID:       l.DocID,
					Name:     l.DocName,
					Size:     l.FileSize,
					Modified: time.UnixMilli(l.LastModifiedTime),
					Ctime:    time.UnixMilli(l.CreateTime),
					IsFolder: l.DocTypeEn == "folder",
				})
			}
		}

		if !r.CosmoExtVoPage.HasNext {
			break
		}
	}
	return f, nil
}

func (d *Qingque) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	var r DownloadResp
	err := d.request(http.MethodGet, "/docs/yfile/download-url/{docID}", func(req *resty.Request) {
		req.SetPathParam("docID", file.GetID())
		req.SetQueryParams(map[string]string{
			"anonToken": "true", // true: support download without cookie
		})
	}, &r)
	if err != nil {
		return nil, err
	}
	return &model.Link{URL: r.FileURL}, nil
}

func (d *Qingque) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	var r FolderNewResp
	err := d.request(http.MethodPost, "/docs", func(req *resty.Request) {
		req.SetBody(base.Json{
			"docTypeEn":   "folder",
			"shareTypeEn": "normal",
			"parentId":    parentDir.GetID(),
			"docName":     dirName,
			"userPhotoId": "",
			"sendMessage": "true",
		})
	}, &r)
	if err != nil {
		return err
	}
	return nil
}

func (d *Qingque) Move(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO move obj, optional
	return nil, errs.NotImplement
}

func (d *Qingque) Rename(ctx context.Context, srcObj model.Obj, newName string) (model.Obj, error) {
	// TODO rename obj, optional
	return nil, errs.NotImplement
}

func (d *Qingque) Copy(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO copy obj, optional
	return nil, errs.NotImplement
}

func (d *Qingque) Remove(ctx context.Context, obj model.Obj) error {
	// TODO remove obj, optional
	return errs.NotImplement
}

func (d *Qingque) Put(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up driver.UpdateProgress) (model.Obj, error) {
	// TODO upload file, optional
	return nil, errs.NotImplement
}

func (d *Qingque) GetArchiveMeta(ctx context.Context, obj model.Obj, args model.ArchiveArgs) (model.ArchiveMeta, error) {
	// TODO get archive file meta-info, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *Qingque) ListArchive(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) ([]model.Obj, error) {
	// TODO list args.InnerPath in the archive obj, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *Qingque) Extract(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) (*model.Link, error) {
	// TODO return link of file args.InnerPath in the archive obj, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *Qingque) ArchiveDecompress(ctx context.Context, srcObj, dstDir model.Obj, args model.ArchiveDecompressArgs) ([]model.Obj, error) {
	// TODO extract args.InnerPath path in the archive srcObj to the dstDir location, optional
	// a folder with the same name as the archive file needs to be created to store the extracted results if args.PutIntoNewDir
	// return errs.NotImplement to use an internal archive tool
	return nil, errs.NotImplement
}

//func (d *Template) Other(ctx context.Context, args model.OtherArgs) (interface{}, error) {
//	return nil, errs.NotSupport
//}

var _ driver.Driver = (*Qingque)(nil)
