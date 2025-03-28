package fnos_share

type Request struct {
	Method string
	URL    string
	Params map[string]string
	Data   interface{}
}

type BaseResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data interface{}
}

type ShareData struct {
	BaseResp
	Data struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	} `json:"data"`
}

type ListResp struct {
	BaseResp
	Data struct {
		Files []struct {
			FileID  int    `json:"fileId"`
			Path    string `json:"path"`
			IsDir   int    `json:"isDir"`
			File    string `json:"file"`
			Size    int    `json:"size"`
			ModTime int    `json:"modTime"`
		} `json:"files"`
	} `json:"data"`
}
