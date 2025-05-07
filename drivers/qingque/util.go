package qingque

import (
	"encoding/json"
	"errors"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/google/uuid"
)

// do others that not defined in Driver interface

func (d *Qingque) request(method string, path string, callback base.ReqCallback, out any) error {
	u := "https://docs.qingque.cn/merlot/api" + path
	req := d.client.R()
	req.SetQueryParam("um", "false")
	req.SetHeaders(map[string]string{
		"Accept": "application/json",
		"rid":    uuid.NewString(),
	})
	if d.IdentityId != "" {
		req.SetHeader("identityId", d.IdentityId)
	}

	var r BaseResp
	req.SetResult(&r)
	if callback != nil {
		callback(req)
	}

	resp, err := req.Execute(method, u)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.New(resp.String())
	}

	if r.Code != 0 {
		return errors.New(r.ResultCode)
	}

	if out != nil && r.Result != nil {
		var marshal []byte
		marshal, err = json.Marshal(r.Result)
		if err != nil {
			return err
		}
		err = json.Unmarshal(marshal, &out)
		if err != nil {
			return err
		}
	}
	return nil
}
