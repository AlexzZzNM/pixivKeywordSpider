package download

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"pixivKwywordSpider/config"
	"pixivKwywordSpider/utils"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// 生成代理 Client
func NewProxyClient(proxyStrSlice ...string) (*http.Client) {
	var useProxyUrl string
	if len(proxyStrSlice) == 1 {
		useProxyUrl = proxyStrSlice[0]
	} else {
		useProxyUrl = proxyStrSlice[rand.Intn(len(proxyStrSlice))]
	}
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(useProxyUrl)
	}

	transport := &http.Transport{Proxy: proxy}

	client := &http.Client{Transport: transport}
	return client
}

// 生成新的下载请求
func NewDownload(url string) (*http.Response, error) {
	var client *http.Client = &http.Client{}
	if len(config.ProxyList) != 0 {
		// 随机 生成 一个连接代理的 client
		client = NewProxyClient(config.ProxyList...)
	}
	//生成要访问的url
	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//增加header选项
	reqest.Header.Add("Referer", "https://www.pixiv.net")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36")

	if err != nil {
		return nil, err
	}
	//处理返回结果
	response, err := client.Do(reqest)
	return response, err
}

// 下载文件并且保存
func DownloadAndSave(url, dir, fileName string) (error) {
	// 对实参进行处理 s
	if dir == "" {
		dir = "./"
	}

	if fileName == "" {
		name, err := utils.UrlGetFileName(url)
		if err != nil {
			return err
		}
		fileName = name
	}
	// 对实参进行处理 e

	newDate := time.Now()
	res, err := NewDownload(url)
	if err != nil {
		return err
	}
	imgPath := path.Join(dir, fileName)
	fileInfo, err := os.Stat(imgPath)
	if err != nil || fileInfo.Size() == 0 {
		// 不存在 或者大小为0
		// 需要下载
	} else {
		// 如果存在且大小不为0， 不下载
		return nil
	}
	f, err := os.OpenFile(imgPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	bys, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(bys)
	if err != nil {
		return err
	}
	fmt.Println( "下载文件", fileName, "耗时：==》" ,time.Now().Sub(newDate))
	return nil
}
