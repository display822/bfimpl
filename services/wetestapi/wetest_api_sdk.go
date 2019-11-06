package wetestapi

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/url"

	"github.com/levigross/grequests"
)

type (
	WeTestResponse struct {
		Ret int16  `json:"ret"`
		Msg string `json:"msg"`
	}

	WeTestApiClient struct {
		Appid      string
		Secret     string
		Host       string
		SESSION_ID string
		DataCache  DataCache
		Trail      Trail
	}

	DataCache struct {
		UserData interface{}
		UserCorp interface{}
		UserTeam interface{}
		Pivilege interface{}
		Product  interface{}
	}

	Trail struct {
		method string
		Errno  WeTestResponse
		Error  interface{}
	}
)

const (
	CurlErr      = 100
	paramErr     = 101
	notNull      = 102
	receiverErr  = 103
	contentErr   = 104
	mustParamErr = 105
	VERSION      = "v3"
)

func hosts(c string) string {
	var host string
	switch c {
	case "pub":
		host = "http://api.wetest.qq.com"
	case "pre":
		host = "http://api.wepub.qq.com"
	case "dev":
		host = "http://api.user.openqa.qq.com"
	}
	return host
}

func sessions(c string) string {
	var session string
	switch c {
	case "pub":
		session = "wetest_sessionid"
	case "pre":
		session = "wepub_sessionid"
	case "dev":
		session = "openqa_sessionid"
	}
	return session
}

/**
 * 初始化
 * @param  appid
 * @param  secret
 * @param  string  env
 * @return bool|null|v3sdk
 */
func NewClient(appid string, secret string, env string) *WeTestApiClient {
	if appid == "" || secret == "" {
		return nil
	}
	cons := &WeTestApiClient{Appid: appid, Secret: secret + "&", Host: hosts(env), SESSION_ID: sessions(env), Trail: struct {
		method string
		Errno  WeTestResponse
		Error  interface{}
	}{method: "GET"}}
	return cons
}

/**
 * 设置method
 * @param string method_param
 */
func (wta *WeTestApiClient) setMethod(method_param string) {
	if method_param == "POST" {
		wta.Trail.method = "POST"
	}
}

/**
 * 底层链接
 * @param  path
 * @param  params
 * @return bool|mixed
 */
func (wta *WeTestApiClient) apiCurl(path string, params map[string]string) map[string]interface{} {
	url_path := fmt.Sprintf("%s/%s/%s/", wta.Host, VERSION, path)
	sign_path := fmt.Sprintf("/%s/%s/", VERSION, path)
	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
	}
	params["appid"] = wta.Appid
	params["t"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["sign"] = wta.makeSig(sign_path, params)
	if wta.Trail.method == "GET" {
		str := ""
		for k, v := range params {
			str += "&" + fmt.Sprintf("%s=%s", k, v)
		}
		str = strings.TrimLeft(str, "&")
		if str != "" {
			has := strings.Contains(url_path, "?")
			if has == false {
				url_path += "?" + str
			} else {
				url_path += "&" + str
			}
		}
		resp, err := grequests.Get(url_path, nil)
		if err != nil {
			wta.Trail.Errno.Ret = CurlErr
			wta.Trail.Errno.Msg = err.Error()
			return nil
		} else {
			m := make(map[string]interface{})
			json.Unmarshal([]byte(resp.String()), &m)
			return m
		}
	} else {
		resp, err := grequests.Post(url_path,
			&grequests.RequestOptions{Data: params})
		if err != nil {
			wta.Trail.Errno.Ret = CurlErr
			wta.Trail.Errno.Msg = "json fatal"
			return nil
		} else {
			m := make(map[string]interface{})
			json.Unmarshal([]byte(resp.String()), &m)
			return m
		}
	}
}

/**
 * 创建sign
 * @param  url_path
 * @param  params
 * @return string
 */
func (wta *WeTestApiClient) makeSig(url_path string, params map[string]string) string {
	mk := wta.makeSource(url_path, params)
	key := []byte(wta.Secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(mk))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return sign
}

/**
 * 创建有序url
 * @param  url_path
 * @param  params
 * @return string
 */
func (wta *WeTestApiClient) makeSource(url_path string, params map[string]string) string {
	strs := wta.Trail.method + "&" + url.QueryEscape(url_path) + "&"
	var sslice []string
	for key, _ := range params {
		sslice = append(sslice, key)
	}
	sort.Strings(sslice)
	query_string := ""
	for _, v := range sslice {
		query_string = query_string + "&" + fmt.Sprintf("%s=%s", v, params[v])
	}
	query_string = strings.TrimLeft(query_string, "&")
	query_string = strings.Replace(url.QueryEscape(query_string), "~", "%7E", -1)
	return strs + query_string
}

/**
 * 获取用户
 * @param  bool showCorp
 * @param  bool showProj
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getUser(showCorp bool, showProj bool) bool {
	path := "get_user"
	params := map[string]string{
		"sessionid": wta.SESSION_ID,
	}
	if showCorp == true {
		params["corporation"] = "1"
	}
	if showProj == true {
		params["team"] = "1"
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		if _, ok := result["corporation"]; ok {
			wta.DataCache.UserCorp = result["corporation"]
			delete(result, "corporation")
		}
		if _, ok := result["team"]; ok {
			wta.DataCache.UserTeam = result["team"]
			delete(result, "team")
		}
	} else {
		wta.Trail.Error = result
		return false
	}
	delete(result, "ret")
	wta.DataCache.UserData = result
	return true
}

/**
 * 获取用户数据
 * @return bool|null
 */
func (wta *WeTestApiClient) getUserData() interface{} {
	if wta.DataCache.UserData == nil {
		if wta.getUser(false, false) == false {
			if wta.Trail.Errno != (WeTestResponse{}) {
				return wta.Trail.Errno
			} else {
				return wta.Trail.Error
			}
		}
	}
	return wta.DataCache.UserData
}

/**
 * 获取用户企业
 * @return null
 */
func (wta *WeTestApiClient) getUserCorp() interface{} {
	if wta.DataCache.UserCorp == nil {
		if wta.getUser(true, false) == false {
			if wta.Trail.Errno != (WeTestResponse{}) {
				return wta.Trail.Errno
			} else {
				return wta.Trail.Error
			}
		}
	}
	return wta.DataCache.UserCorp
}

/**
 * 获取用户团队
 * @return null
 */
func (wta *WeTestApiClient) getUserTeam() interface{} {
	if wta.DataCache.UserTeam == nil {
		if wta.getUser(false, true) == false {
			if wta.Trail.Errno != (WeTestResponse{}) {
				return wta.Trail.Errno
			} else {
				return wta.Trail.Error
			}
		}
	}
	return wta.DataCache.UserTeam
}

/**
 * 获取团队详情
 * @param array  list
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getTeamDetail(list []interface{}) interface{} {
	path := "get_team"
	if list == nil {
		return WeTestResponse{notNull, "参数不能为空"}
	}
	lists := strings.Replace(strings.Trim(fmt.Sprint(list), "[]"), " ", ",", -1)
	params := map[string]string{
		"team": lists,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["team"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 获取企业详情
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getCorpDetail() interface{} {
	path := "get_corporation"
	params := map[string]string{
		"sessionid": wta.SESSION_ID,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["corporation"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 下单接口
 * @param  key
 * @param  content
 * @param  num
 * @return bool|mixed
 */
func (wta *WeTestApiClient) generateOrder(key string, content string, num interface{}) interface{} {
	path := "generate_order"
	params := map[string]string{
		"sessionid": wta.SESSION_ID,
		"key":       key,
		"content":   content,
	}
	switch num.(type) {
	case int:
		params["num"] = num.(string)
	case []interface{}:
		params["type"] = "choices"
		choice, _ := json.Marshal(num)
		params["choices"] = string(choice)
	default:
		return WeTestResponse{paramErr, "num参数错误"}
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["order_id"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 验证订单
 * @param  orders
 * @return bool|mixed
 */
func (wta *WeTestApiClient) checkOrder(orders interface{}) interface{} {
	if orders == nil {
		return WeTestResponse{notNull, "参数不能为空"}
	}
	path := "check_order"
	var order string
	switch orders.(type) {
	case []interface{}:
		a, _ := orders.([]interface{})
		order = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		order = orders.(string)
	}
	params := map[string]string{
		"order": order,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["order"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 退还次数
 * @param  order_id
 * @param  num
 * @return bool|mixed
 */
func (wta *WeTestApiClient) refundOrder(order_id string, num int64) interface{} {
	path := "refund_order"
	params := map[string]string{
		"order": order_id,
		"num":   strconv.FormatInt(num, 10),
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 创建团队
 * @param  name
 * @return bool|mixed
 */
func (wta *WeTestApiClient) createTeam(name string) interface{} {
	path := "create_team"
	params := map[string]string{
		"name":      name,
		"sessionid": wta.SESSION_ID,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["team"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 创建用户
 * @param  user array('name'=>'xx','qq'=>111,'rtx'=>'xx','email'=>'xx','phone'=>'xxx')
 * @return bool|mixed
 */
func (wta *WeTestApiClient) createUser(user map[string]string) interface{} {
	path := "create_user"
	result := wta.apiCurl(path, user)
	if result != nil && result["ret"] == 0 {
		delete(result, "ret")
		return result
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 获取当前工具下所有权限
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getPrivilege() interface{} {
	if wta.DataCache.Pivilege == nil {
		path := "get_privilege"
		params := map[string]string{
			"sessionid": wta.SESSION_ID,
		}
		result := wta.apiCurl(path, params)
		if result != nil && result["ret"] == 0 {
			wta.DataCache.Pivilege = result["privilege"]
		} else {
			if wta.Trail.Errno != (WeTestResponse{}) {
				wta.DataCache.Pivilege = wta.Trail.Errno
			} else {
				wta.DataCache.Pivilege = result
			}
		}
	}
	return wta.DataCache.Pivilege
}

/**
 * 获取用户所有产品
 * @return bool|mixed|null
 */
func (wta *WeTestApiClient) getProduct() interface{} {
	if wta.DataCache.Product == nil {
		path := "get_product"
		params := map[string]string{
			"sessionid": wta.SESSION_ID,
		}
		result := wta.apiCurl(path, params)
		if result != nil && result["ret"] == 0 {
			wta.DataCache.Product = result["product"]
		} else {
			if wta.Trail.Errno != (WeTestResponse{}) {
				wta.DataCache.Product = wta.Trail.Errno
			} else {
				wta.DataCache.Product = result
			}
		}
	}
	return wta.DataCache.Product
}

/**
 * 获取产品隶属的用户
 * @param  product_id
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getProductOwner(product_id string) interface{} {
	path := "get_product_owner"
	result := wta.apiCurl(path, map[string]string{"product": product_id})
	if result != nil && result["ret"] == 0 {
		return result
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 创建产品
 * @param  params
 * name:产品名称
 * team_id:产品归属团队id，没有则不传
 * apk:产品的apk包名，没有则不传
 * icon:产品图标url
 * url:产品介绍url
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) createProduct(params map[string]string) interface{} {
	path := "create_product"
	if _, ok := params["name"]; !ok {
		return WeTestResponse{paramErr, "参数错误"}
	}
	params["sessionid"] = wta.SESSION_ID
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["product"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 获取用户列表数据
 * @param  users
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getUsersInfo(users interface{}) interface{} {
	if users == nil {
		return WeTestResponse{notNull, "参数不能为空"}
	}
	path := "get_user_information"
	var user string
	switch users.(type) {
	case []interface{}:
		a, _ := users.([]interface{})
		user = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		user = users.(string)
	}
	params := map[string]string{
		"uin": user,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["user"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 获取用户列表数据
 * @param  uins
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getUsersId(uins interface{}) interface{} {
	if uins == nil {
		return WeTestResponse{notNull, "参数不能为空"}
	}
	path := "get_user_id"
	var uin string
	switch uins.(type) {
	case []interface{}:
		a, _ := uins.([]interface{})
		uin = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		uin = uins.(string)
	}
	param := map[string]string{
		"uin": uin,
	}
	result := wta.apiCurl(path, param)
	if result != nil && result["ret"] == 0 {
		return result["user"]
	} else {
		return result
	}
}

/**
 * 获取用户idcard信息
 * @param  users
 * @return bool|mixed
 */
func (wta *WeTestApiClient) getUsersIdCard(users interface{}) interface{} {
	if users == nil {
		return WeTestResponse{notNull, "参数不能为空"}
	}
	path := "get_user_idcard_information"
	var user string
	switch users.(type) {
	case []interface{}:
		a, _ := users.([]interface{})
		user = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		user = users.(string)
	}
	params := map[string]string{
		"uin": user,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return result["user"]
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 上报测试数据
 * @param  report
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) reportTest(report map[string]string) interface{} {
	path := "report_test"
	if _, ok := report["uin"]; !ok {
		return WeTestResponse{mustParamErr, "缺少必填参数"}
	}
	if _, ok := report["id"]; !ok {
		return WeTestResponse{mustParamErr, "缺少必填参数"}
	}
	if _, ok := report["type"]; !ok {
		return WeTestResponse{mustParamErr, "缺少必填参数"}
	}
	if _, ok := report["status"]; !ok {
		return WeTestResponse{mustParamErr, "缺少必填参数"}
	}
	if _, ok := report["name"]; !ok {
		return WeTestResponse{mustParamErr, "缺少必填参数"}
	}
	result := wta.apiCurl(path, report)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * tof发消息
 * @param  params
 * @return bool|mixed
 */
func (wta *WeTestApiClient) _tofMessage(params map[string]string) interface{} {
	path := "tof_message"
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 发送短信
 * @param  receivers
 * @param  content
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendTofSMS(receivers interface{}, content string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" {
		return WeTestResponse{contentErr, "内容不能为空"}
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"type":     "sms",
		"receiver": receiver,
		"content":  content,
	}
	return wta._tofMessage(params)
}

/**
 * 发送邮件
 * @param  receivers
 * @param  title
 * @param  content
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendTofEmail(receivers interface{}, title string, content string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" || title == "" {
		return WeTestResponse{contentErr, "内容/标题不能为空"}
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"type":     "email",
		"receiver": receiver,
		"content":  content,
		"title":    title,
	}
	return wta._tofMessage(params)
}

/**
 * 发送rtx
 * @param  receivers
 * @param  title
 * @param  content
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendTofRtx(receivers interface{}, title string, content string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" || title == "" {
		return WeTestResponse{contentErr, "内容/标题不能为空"}
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"type":     "rtx",
		"receiver": receiver,
		"content":  content,
		"title":    title,
	}
	return wta._tofMessage(params)
}

/**
 * 发微信消息
 * @param  receivers
 * @param  content
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendTofWeChat(receivers interface{}, content string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" {
		return WeTestResponse{contentErr, "内容不能为空"}
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"type":     "wechat",
		"receiver": receiver,
		"content":  content,
	}
	return wta._tofMessage(params)
}

/**
 * 发送即时消息
 * @param  uins
 * @param  content string 内容
 * @param  msg string 链接内文字
 * @param  url string 链接
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendMessage(uins interface{}, content string, msg string, url string) interface{} {
	path := "tool_message"
	if uins == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" || msg == "" || url == "" {
		return WeTestResponse{contentErr, "内容不能为空"}
	}
	var uin string
	switch uins.(type) {
	case []interface{}:
		a, _ := uins.([]interface{})
		uin = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		uin = uins.(string)
	}
	params := map[string]string{
		"uin":     uin,
		"content": content,
		"msg":     msg,
		"url":     url,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 发送用户邮件
 * @param  receivers
 * @param  title
 * @param  content
 * @param int  code
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendEmail(receivers interface{}, title string, content string, code string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" || title == "" {
		return WeTestResponse{contentErr, "内容/标题不能为空"}
	}
	if code == "" {
		code = "10000"
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"to":      receiver,
		"title":   title,
		"content": content,
		"code":    code,
	}
	result := wta.apiCurl("user_email", params)
	if result != nil {
		return result
	} else {
		return wta.Trail.Errno
	}
}

/**
 * 发送
 * @param  receivers
 * @param  title
 * @param  content
 * @return array|int
 */
func (wta *WeTestApiClient) sendRawEmail(receivers interface{}, title string, content string) interface{} {
	if receivers == nil {
		return WeTestResponse{receiverErr, "接收人错误"}
	}
	if content == "" || title == "" {
		return WeTestResponse{contentErr, "内容/标题不能为空"}
	}
	var receiver string
	switch receivers.(type) {
	case []interface{}:
		a, _ := receivers.([]interface{})
		receiver = strings.Replace(strings.Trim(fmt.Sprint(a), "[]"), " ", ",", -1)
	default:
		receiver = receivers.(string)
	}
	params := map[string]string{
		"to":      receiver,
		"title":   title,
		"content": content,
	}
	result := wta.apiCurl("user_raw_email", params)
	if result != nil {
		return result
	} else {
		return wta.Trail.Errno
	}
}

/**
 * 发送用户短信
 * @param  phone
 * @param  msg
 * @param  code
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendSMS(phone string, msg string, code string) interface{} {
	if phone == "" || msg == "" {
		return WeTestResponse{contentErr, "内容/标题不能为空"}
	}
	params := map[string]string{
		"phone":   phone,
		"content": msg,
	}
	if code != "" {
		params["code"] = code
	}
	result := wta.apiCurl("send_sms", params)
	if result != nil {
		return result
	} else {
		return wta.Trail.Errno
	}
}

/**
 * 创建工具申请
 * @param  tool_application_id  工具维护的申请id
 * @param  kwargs             工具维护弹窗参数
 * @param  content            申请内容
 * @param  reason            申请理由
 * @return bool|mixed
 */
func (wta *WeTestApiClient) createToolApplication(tool_application_id string, content string, reason string, kwargs interface{}) interface{} {
	path := "tool_application"
	if tool_application_id == "" || content == "" {
		return WeTestResponse{contentErr, "申请id不能为空"}
	}
	kwarg, _ := json.Marshal(kwargs)
	params := map[string]string{
		"sessionid":      wta.SESSION_ID,
		"kwargs":         string(kwarg),
		"op":             "0",
		"content":        content,
		"reason":         reason,
		"application_id": tool_application_id,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 工具处理申请
 * @param  tool_application_id
 * @param  permit
 * @param  reason
 * @return bool|mixed
 */
func (wta *WeTestApiClient) dealToolApplication(tool_application_id string, permit string, reason string) interface{} {
	path := "tool_application"
	if tool_application_id == "" {
		return WeTestResponse{contentErr, "申请id不能为空"}
	}
	if reason == "" {
		return WeTestResponse{contentErr, "拒绝申请理由不能为空"}
	}
	var op string
	if permit != "" {
		op = "1"
	} else {
		op = "2"
	}
	params := map[string]string{
		"sessionid":      wta.SESSION_ID,
		"op":             op,
		"reason":         reason,
		"application_id": tool_application_id,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 撤销申请
 * @param  tool_application_id
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) cancelToolApplication(tool_application_id string) interface{} {
	path := "tool_application"
	if tool_application_id == "" {
		return WeTestResponse{contentErr, "申请id不能为空"}
	}
	params := map[string]string{
		"sessionid":      wta.SESSION_ID,
		"op":             "3",
		"application_id": tool_application_id,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 修改申请
 * @param  tool_application_id
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) modifyToolApplication(tool_application_id string) interface{} {
	path := "tool_application"
	if tool_application_id == "" {
		return WeTestResponse{contentErr, "申请id不能为空"}
	}
	params := map[string]string{
		"sessionid":      wta.SESSION_ID,
		"op":             "4",
		"application_id": tool_application_id,
	}
	result := wta.apiCurl(path, params)
	if result != nil && result["ret"] == 0 {
		return true
	} else {
		if wta.Trail.Errno != (WeTestResponse{}) {
			return wta.Trail.Errno
		} else {
			return result
		}
	}
}

/**
 * 发送企业微信机器人消息
 * @param  phone
 * @param  msg
 * @param  code
 * @return array|bool|mixed
 */
func (wta *WeTestApiClient) sendRobotMsg(rtx string, team string, markdown string) interface{} {
	if rtx == "" || team == "" {
		return WeTestResponse{contentErr, "收信人不能为空"}
	}
	if markdown == "" {
		return WeTestResponse{contentErr, "内容不能为空"}
	}
	params := map[string]string{
		"rtx":      rtx,
		"team":     team,
		"markdown": markdown,
	}
	wta.Trail.method = "POST"
	result := wta.apiCurl("wx_robot_msg_push", params)
	wta.Trail.method = "GET"
	if result != nil {
		return result
	} else {
		return wta.Trail.Errno
	}
}
