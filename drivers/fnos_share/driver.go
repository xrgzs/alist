package fnos_share

import (
	"context"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/model"
)

type FnOSShare struct {
	model.Storage
	Addition
}

func (d *FnOSShare) Config() driver.Config {
	return config
}

func (d *FnOSShare) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *FnOSShare) Init(ctx context.Context) error {
	// TODO login / refresh token
	//op.MustSaveDriverStorage(d)
	return nil
}

func (d *FnOSShare) Drop(ctx context.Context) error {
	return nil
}

func (d *FnOSShare) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	// TODO return the files list, required
	token, err := d.getToken()
	if err != nil {
		return nil, err
	}
	_, err = d.request(d.ShareId, "POST", token, nil, &r)
	if err != nil {
		return nil, err
	}

	return nil, errs.NotImplement
}

func (d *FnOSShare) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	// TODO return link of file, required
	return nil, errs.NotImplement
}

func (d *FnOSShare) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) (model.Obj, error) {
	// TODO create folder, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) Move(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO move obj, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) Rename(ctx context.Context, srcObj model.Obj, newName string) (model.Obj, error) {
	// TODO rename obj, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) Copy(ctx context.Context, srcObj, dstDir model.Obj) (model.Obj, error) {
	// TODO copy obj, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) Remove(ctx context.Context, obj model.Obj) error {
	// TODO remove obj, optional
	return errs.NotImplement
}

func (d *FnOSShare) Put(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up driver.UpdateProgress) (model.Obj, error) {
	// TODO upload file, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) GetArchiveMeta(ctx context.Context, obj model.Obj, args model.ArchiveArgs) (model.ArchiveMeta, error) {
	// TODO get archive file meta-info, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) ListArchive(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) ([]model.Obj, error) {
	// TODO list args.InnerPath in the archive obj, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) Extract(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) (*model.Link, error) {
	// TODO return link of file args.InnerPath in the archive obj, return errs.NotImplement to use an internal archive tool, optional
	return nil, errs.NotImplement
}

func (d *FnOSShare) ArchiveDecompress(ctx context.Context, srcObj, dstDir model.Obj, args model.ArchiveDecompressArgs) ([]model.Obj, error) {
	// TODO extract args.InnerPath path in the archive srcObj to the dstDir location, optional
	// a folder with the same name as the archive file needs to be created to store the extracted results if args.PutIntoNewDir
	// return errs.NotImplement to use an internal archive tool
	return nil, errs.NotImplement
}

//func (d *FnOSShare) Other(ctx context.Context, args model.OtherArgs) (interface{}, error) {
//	return nil, errs.NotSupport
//}

var _ driver.Driver = (*FnOSShare)(nil)
