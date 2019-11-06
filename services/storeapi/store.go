package storeapi

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type StoreGo struct {
	buid          string
	key           string
	httpstorehost string
	timeout       time.Duration
	append        int
	token         int
	savedays      int
	unzip         int
}

var errFileExists = errors.New("-15_file exist")

type RetMsg struct {
	Retcode int    `json:"retcode"`
	Msg     string `json:"msg"`
}

// store_dir_list
type Dir struct {
	Dir string `json:"dir"`
}
type File struct {
	File string `json:file`
}

type DirListRetMsg struct {
	RetMsg
	DirList  []Dir  `json:"dir_list"`
	FileList []File `json:"file_list"`
}

type FileExistRetMsg struct {
	RetMsg
	FileList map[string]int `json:"file_list"`
}

const (
	API_STORE_FILE_UPLOAD    = "/store_file_upload?"
	API_STORE_FILE_DIRCREATE = "/store_dir_create?"
	API_STORE_FILE_EXIST     = "/store_file_exist?"
	API_STORE_DIR_LIST       = "/store_dir_list?"
	API_STORE_FILE_DELETE    = "/store_file_delete?"
	API_STORE_DIR_DELETE     = "/store_dir_delete?"
	METHOD_GET               = "GET"
	METHOD_POST              = "POST"
)

func NewStore(storehost, buid, key string) *StoreGo {
	s := &StoreGo{buid, key, storehost, time.Second * 20, 0, 1, 365, 0}
	return s
}

// optoins
func (s *StoreGo) SetTimeOut(newtimeout time.Duration) {
	s.timeout = newtimeout
}

func (s *StoreGo) SetAppend(app int) {
	s.append = app
}

func (s *StoreGo) SetNeedToken(token int) {
	s.token = token
}

func (s *StoreGo) SetSaveDays(savedays int) {
	s.savedays = savedays
}

func (s *StoreGo) SetUnzip(unzip int) {
	s.unzip = unzip
}

func (s *StoreGo) geneateUrlValuesToken() *url.Values {
	var values url.Values = url.Values{}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	randnum := strconv.Itoa(rand.Intn(9999999999))
	var keys []string = []string{s.buid, s.key, now, randnum}
	sort.Strings(keys)
	tmpstr := strings.Join(keys, "")
	sign := fmt.Sprintf("%x", sha1.Sum([]byte(tmpstr)))
	values.Set("buid", s.buid)
	values.Set("t", now)
	values.Set("r", randnum)
	values.Set("sign", sign)
	return &values
}

func (s *StoreGo) generateUrlValuesFileToken() *url.Values {
	values := s.geneateUrlValuesToken()
	values.Set("append", strconv.Itoa(s.append))
	values.Set("token", strconv.Itoa(s.token))
	values.Set("save_days", strconv.Itoa(s.savedays))
	return values
}

func (s *StoreGo) doRequest(request *http.Request) ([]byte, error) {
	client := http.Client{
		Timeout: s.timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 500 {
		return nil, errors.New("server 500")
	}
	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respbytes, nil
}

func (s *StoreGo) doSimpleRequest(method string, cgiurl string, values *url.Values, body io.Reader, contenttype string) error {
	url := s.httpstorehost + cgiurl + values.Encode()
	req, err := http.NewRequest(method, url, body)
	if contenttype != "" {
		req.Header.Add("Content-Type", contenttype)
	}
	if err != nil {
		return err
	}
	respbytes, err := s.doRequest(req)
	if err != nil {
		return err
	}
	var retmsg RetMsg
	err = json.Unmarshal(respbytes, &retmsg)
	if err != nil {
		return errors.New(string(respbytes))
	}
	if retmsg.Retcode != 0 {
		return errors.New(fmt.Sprintf("%d_%s", retmsg.Retcode, retmsg.Msg))
	}
	return nil
}

func (s *StoreGo) doReturnRequest(method string, cgiurl string, values *url.Values, body io.Reader) ([]byte, error) {
	url := s.httpstorehost + cgiurl + values.Encode()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return s.doRequest(req)
}

func (s *StoreGo) CreateDir(dirpath, name string) error {
	params := s.geneateUrlValuesToken()
	params.Set("dir_path", dirpath)
	params.Set("name", name)
	return s.doSimpleRequest(METHOD_GET, API_STORE_FILE_DIRCREATE, params, nil, "")
}

func (s *StoreGo) FileExists(dirpath string, names []string) (bool, error) {
	params := s.geneateUrlValuesToken()
	params.Set("dir_path", dirpath)
	params.Set("name", strings.Join(names, ","))
	respbytes, err := s.doReturnRequest(METHOD_GET, API_STORE_FILE_EXIST, params, nil)
	if err != nil {
		return false, err
	}
	var filelist FileExistRetMsg
	err = json.Unmarshal(respbytes, &filelist)
	if err != nil {
		return false, errors.New(string(respbytes))
	}
	if filelist.Retcode != 0 {
		return false, errors.New(filelist.Msg)
	}
	for _, n := range names {
		if v, ok := filelist.FileList[n]; ok == false || v == 0 {
			return false, nil
		}
	}
	return true, nil
}

func (s *StoreGo) UploadFile(filepath, dirpath, name string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	contenttype := writer.FormDataContentType()
	writer.Close()
	params := s.generateUrlValuesFileToken()
	params.Set("dir_path", dirpath)
	params.Set("name", name)
	return s.doSimpleRequest(METHOD_POST, API_STORE_FILE_UPLOAD, params, body, contenttype)
}

func (s *StoreGo) DirList(dirpath string) (*DirListRetMsg, error) {
	params := s.generateUrlValuesFileToken()
	params.Set("dir_path", dirpath)
	respbytes, err := s.doReturnRequest(METHOD_GET, API_STORE_DIR_LIST, params, nil)
	if err != nil {
		return nil, err
	}
	var dirlist DirListRetMsg
	err = json.Unmarshal(respbytes, &dirlist)
	if err != nil {
		return nil, errors.New(string(respbytes))
	}
	if dirlist.Retcode != 0 {
		return nil, errors.New(dirlist.Msg)
	}
	return &dirlist, nil
}

func (s *StoreGo) DeleteFile(dirpath string, name string) (*RetMsg, error) {
	params := s.generateUrlValuesFileToken()
	params.Set("dir_path", dirpath)
	params.Set("name", name)
	respbytes, err := s.doReturnRequest(METHOD_GET, API_STORE_FILE_DELETE, params, nil)
	if err != nil {
		return nil, err
	}
	var ret RetMsg
	err = json.Unmarshal(respbytes, &ret)
	if err != nil {
		return nil, errors.New(string(respbytes))
	}
	if ret.Retcode != 0 {
		return nil, errors.New(ret.Msg)
	}
	return &ret, nil
}

func (s *StoreGo) DeleteDir(dirpath string, name string) (*RetMsg, error) {
	params := s.generateUrlValuesFileToken()
	params.Set("dir_path", dirpath)
	params.Set("name", name)
	respbytes, err := s.doReturnRequest(METHOD_GET, API_STORE_DIR_DELETE, params, nil)
	if err != nil {
		return nil, err
	}
	var ret RetMsg
	err = json.Unmarshal(respbytes, &ret)
	if err != nil {
		return nil, errors.New(string(respbytes))
	}
	if ret.Retcode != 0 {
		return nil, errors.New(ret.Msg)
	}
	return &ret, nil
}

func (s *StoreGo) GetDownloadPath(dirpath, name string) string {
	params := url.Values{}
	randnum := strconv.Itoa(rand.Intn(9999999999))
	content := fmt.Sprintf("%s|%s|%s|%s|%s", s.buid, dirpath, s.key, name, randnum)
	token := fmt.Sprintf("%x", sha1.Sum([]byte(content)))
	params.Set("token", token)
	params.Set("r", randnum)
	return fmt.Sprintf("%s/gqop/%s%s%s?%s", s.httpstorehost, s.buid, dirpath, name, params.Encode())
}

func (s *StoreGo) GetDownloadPathStoreFileDownload(dirpath, name string, randn uint64) string {
	var randnum string
	if randn == 0 {
		randnum = strconv.Itoa(rand.Intn(9999999999))
	} else {
		randnum = strconv.FormatUint(randn, 10)
	}
	content := fmt.Sprintf("%s|%s|%s|%s|%s", s.buid, dirpath, s.key, name, randnum)
	token := fmt.Sprintf("%x", sha1.Sum([]byte(content)))
	return fmt.Sprintf("%s/store_file_download?buid=%v&dir_path=%v&r=%v&token=%v&name=%v",
		s.httpstorehost, s.buid, dirpath, randnum, token, name)
}

func (s *StoreGo) GetDownloadPathStoreFilePartDownload(dirpath, name string, suffix string, randn uint64, fname string) string {
	var randnum string
	if randn == 0 {
		randnum = strconv.Itoa(rand.Intn(9999999999))
	} else {
		randnum = strconv.FormatUint(randn, 10)
	}
	content := fmt.Sprintf("%s|%s|%s|%s|%s", s.buid, dirpath, s.key, name, randnum)
	token := fmt.Sprintf("%x", sha1.Sum([]byte(content)))
	return fmt.Sprintf("%s/store_file_part_download?buid=%v&dir_path=%v&r=%v&token=%v&ext=%s&name=%v&fname=%v",
		s.httpstorehost, s.buid, dirpath, randnum, token, suffix, name, fname)
}

func (s *StoreGo) ComputeMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(result)), nil
}
