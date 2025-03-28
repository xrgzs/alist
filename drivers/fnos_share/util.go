package fnos_share

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/pkg/utils/random"
	"github.com/dlclark/regexp2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// do others that not defined in Driver interface
func (d *FnOSShare) request(path string, method string, token string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	url := d.Host + path
	req := base.RestyClient.R()

	if token != "" {
		req.SetHeader("Auth", token)
		request := &Request{
			Method: method,
			URL:    url,
		}
		signature := d.getSign(request, token)
		fmt.Println("Generated Signature:", signature)
		req.SetHeader("Authx", signature)
	}

	if callback != nil {
		callback(req)
	}
	var r BaseResp
	req.SetResult(&r)
	res, err := req.Execute(method, url)
	log.Debugln(res.String())
	if err != nil {
		return nil, err
	}

	// 业务状态码检查（优先于HTTP状态码）
	if r.Code != 0 {
		return res.Body(), errors.New(r.Msg)
	}
	if resp != nil {
		err = utils.Json.Unmarshal(res.Body(), resp)
		if err != nil {
			return nil, err
		}
	}
	return res.Body(), nil
}

func (d *FnOSShare) stringifyParams(params map[string]string) string {
	v := url.Values{}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v.Add(k, params[k])
	}
	return strings.ReplaceAll(v.Encode(), "+", "%20")
}

func (d *FnOSShare) parseUrl(rawURL string) (string, map[string]string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, make(map[string]string)
	}

	queryParams := make(map[string]string)
	for key, values := range u.Query() {
		if len(values) > 0 && values[0] != "undefined" && values[0] != "null" {
			queryParams[key] = values[0]
		}
	}
	return u.Path, queryParams
}

func (d *FnOSShare) hashSignatureData(s string) string {
	// 处理无效的%符号
	re := regexp2.MustCompile(`%(?![0-9A-Fa-f]{2})`, 0)
	processed, err := re.Replace(s, "%25", -1, -1)
	if err != nil {
		processed = s
	}

	// 尝试URL解码
	decoded, err := url.QueryUnescape(processed)
	if err != nil {
		decoded = s
	}

	// 返回参数的MD5校验和
	return fmt.Sprint(md5.Sum([]byte(decoded)))
}

func (d *FnOSShare) getSign(req *Request, secret string) string {
	if req == nil {
		return ""
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GenerateSignature error:", r)
		}
	}()

	// 处理HTTP方法
	method := strings.ToUpper(req.Method)

	// 解析URL
	path, queryParams := d.parseUrl(req.URL)

	// 构建签名基础数据
	var b string
	if method == "GET" {
		// 合并参数
		mergedParams := make(map[string]string)
		maps.Copy(mergedParams, req.Params)
		maps.Copy(mergedParams, queryParams)
		b = d.stringifyParams(mergedParams)
	} else if req.Data != nil {
		// 序列化请求体
		data, err := utils.Json.Marshal(req.Data)
		if err != nil {
			log.Errorf("failed to marshal json: %+v", err)
			return ""
		}
		b = string(data)
	}

	// 计算数据哈希
	hashedB := d.hashSignatureData(b)

	// 生成随机数和时间戳
	nonce := strconv.FormatInt(random.RangeInt64(100000, 999999), 10)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// 构建最终签名字符串
	tt := []string{
		"NDzZTVxnRKP8Z0jXg1VAMonaG8akvh",
		path,
		nonce,
		timestamp,
		hashedB,
		secret,
	}
	signStr := strings.Join(tt, "_")

	// 计算最终签名
	sign := fmt.Sprint(md5.Sum([]byte(signStr)))

	// 构建返回参数
	params := url.Values{}
	params.Add("nonce", nonce)
	params.Add("timestamp", timestamp)
	params.Add("sign", sign)

	return params.Encode()
}

func (d *FnOSShare) getToken(id string) (string, error) {

	res, err := d.request(d.ShareId, "GET", "", nil, nil)
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return "", err
	}

	// 解析 HTML DOM，找到 <script id="share-data">
	var shareDataJSON string
	var findShareData func(n *html.Node)
	findShareData = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "share-data" {
					if n.FirstChild != nil {
						shareDataJSON = n.FirstChild.Data
						return
					}
				}
			}
		}
		// 递归遍历子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findShareData(c)
		}
	}
	findShareData(doc)
	if shareDataJSON == "" {
		return "", fmt.Errorf("Failed to find share-data in the HTML: %s", doc)
	}

	var data ShareData
	err = utils.Json.Unmarshal([]byte(shareDataJSON), data)
	if err != nil {
		return "", err
	}
	return data.Data.Token, nil
}
