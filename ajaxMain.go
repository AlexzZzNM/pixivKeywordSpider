package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"os"
	"pixivKwywordSpider/config"
	"pixivKwywordSpider/download"
	"pixivKwywordSpider/useType"
	"sort"
	"sync"
	"time"
)

// p站改版后 为 ajax 下载， 原有的爬取失效, 改为本页爬取

var nowPage = 2460

var getPageChan chan string = make(chan string, config.MaxGoroutine) // 控制最大线程
var wg sync.WaitGroup                                                // 保证线程执行完毕
var AllList []useType.DetailsInfoTwo

/*func sortKey()  {
	
}*/


func SortAllList(list []useType.DetailsInfoTwo) {
	sort.Slice(list, func(i, j int) bool {
		k1 := getMaPkey(list[i].Illust)[0]
		k2 := getMaPkey(list[j].Illust)[0]
		return list[i].Illust[k1].BookmarkCount > list[j].Illust[k2].BookmarkCount
	})
}

func getMaPkey(newMap map[string]useType.DetailsInfoID) []string {
	arr := []string{}
	for key := range newMap{
		arr = append(arr, key)
	}
	return arr
}

func main() {

	star()
	wg.Wait()

	// 总数据拉取完毕后进行排序写入
	SortAllList(AllList)
	writeFile()

	/*str := utils.UrlGetPath("https://i.pximg.net/img-original/img/2020/08/12/17/05/16/83625783_p0.jpg")
	fmt.Println(str)

	fileName, _ := utils.UrlGetFileName("https://i.pximg.net/img-original/img/2020/08/12/17/05/16/83625783_p0.jpg")

	newFileName := strings.Replace(fileName, "p0", fmt.Sprintf("p%d", 5), 1)
	fmt.Println(newFileName)*/
}

func writeFile()  {
	jsonBytes, err := json.Marshal(AllList)
	if err != nil {
		fmt.Println("爬取数据转 json 失败")
		return
	}

	// 将信息写入json 文件
	fileName := config.KeyWord + time.Now().Format("2006-01-02") + ".json"

	f, err := os.OpenFile("./"+fileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err, "=> 打开文件失败")
	}
	defer f.Close()
	_, err = f.Write(jsonBytes)
	if err != nil {
		fmt.Println(err, "=> 爬取数据写入失败")
	}
}

// 获取关键字页面 的url
func getPageUrl(keyword string, nowPage int) string {
	var enKey = url.QueryEscape(keyword)
	return fmt.Sprintf("https://www.pixiv.net/ajax/search/artworks/%s?word=%s&order=date_d&mode=r18&p=%d&s_mode=s_tag&type=all&lang=zh", enKey, enKey, nowPage)
}

// 获取关键字页面 的列表数据
func getPageData(url string) (useType.AjaxPageData, error) {
	var pageData useType.AjaxPageData
	res, err := download.NewDownload(url)
	if err != nil {
		return pageData, err
	}
	defer res.Body.Close()
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return pageData, err
	}
	err = json.Unmarshal(resByte, &pageData)
	if err != nil {
		return pageData, err
	}
	return pageData, err
}

func star() {
	for true {
		pageUrl := getPageUrl(config.KeyWord, nowPage)
		pageData, err := getPageData(pageUrl)
		idMap := make(map[string]bool)

		getPageChan <- ""
		wg.Add(1)
		fmt.Printf("当前正在请求第 %d 页, 请求地址：%s \n", nowPage, pageUrl)

		// 如果失败，每页尝试三次
		if err != nil {
			for i := 0; i < 6; i++ {
				fmt.Printf("重复请求第%d次： 当前正在请求第 %d 页, 请求地址：%s \n",i, nowPage, pageUrl)
				pageData, err = getPageData(pageUrl)
				if err == nil {
					break
				}
			}
		}

		// 无论是否成功， 必须开始下一页
		nowPage++

		// 每次执行本方法后进行数据写入，防止因意外得不到数据
		if len(AllList) != 0 {
			writeFile()
		}



		if err != nil {
			fmt.Printf("页面 %d 请求失败， 本页数据将缺失 \n", nowPage)
		} else if len(pageData.Body.IllustManga.Data) >= 0 {
			// 拿到了数据, 先判断是否拉取完毕
			if len(pageData.Body.IllustManga.Data) == 0 {
				fmt.Printf("爬取完成")
				<-getPageChan
				wg.Done()
				break
			}
			// 开始详情页
			for _, v := range pageData.Body.IllustManga.Data {
				//fmt.Println(v)
				nowData := v
				// 判断是否以及 有数据s
				if idMap[nowData.ID] {
					continue
				} else {
					idMap[nowData.ID] = true
				}
				// 判断是否以及 有数据e

				getPageChan <- ""
				wg.Add(1)

				go func() {

					detailsData, err := getDetailsInfo(nowData.ID)

					<-getPageChan
					wg.Done()

					if err == nil {
						AllList = append(AllList, detailsData)
					} else {
						fmt.Println("----详细拉取失败---")
					}

				}()
			}
		}


		<-getPageChan
		wg.Done()
	}

}

// 根据id 获取详情信息
func getDetailsInfo(id string) (useType.DetailsInfoTwo, error) {
	var detailsInfo useType.DetailsInfoTwo
	detailsPgaeUrl := fmt.Sprintf("https://www.pixiv.net/artworks/%s", id)
	res, err := download.NewDownload(detailsPgaeUrl)

	fmt.Printf("当前正在请求第 %d 页, id为：%s  \n", nowPage - 1, id )

	// 依旧是尝试三次
	if err != nil {
		for i := 0; i < 6; i++ {
			fmt.Printf("重复请求 第 %d次： 当前正在请求第 %d 页, id为：%s  \n", i, nowPage - 1, id )
			res, err = download.NewDownload(detailsPgaeUrl)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		fmt.Printf("详情页 %s 请求失败， 数据将缺失 \n", id)
	} else {
		defer res.Body.Close()
		/*if err != nil {
			fmt.Printf("详情页 %s 转换失败， 数据将缺失 \n", id)
		} else {

		}*/
		dom,err :=goquery.NewDocumentFromReader(res.Body)
		if err != nil {
		} else {
			dom.Find("#meta-preload-data").Each(func(i int, selection *goquery.Selection) {
				dataStr, _ := selection.Attr("content")
				err = json.Unmarshal([]byte(dataStr), &detailsInfo)
			})
		}
	}
	return detailsInfo, err

}

