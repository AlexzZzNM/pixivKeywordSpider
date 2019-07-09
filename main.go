package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
	"pixivKwywordSpider/config"
	"pixivKwywordSpider/download"
	"pixivKwywordSpider/useType"
	"pixivKwywordSpider/utils"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex sync.Mutex

func main() {

	fmt.Println("请输入要搜索的关键字：")
	keyWord := ""
	fmt.Scanf("%v", &keyWord)

	if keyWord == "" || keyWord == " " {
		keyWord = config.KeyWord
	}

	// page 列表页
	var CrawlUrl string = fmt.Sprintf("https://www.pixiv.net/search.php?word=%s&order=date_d&p=", keyWord)
	// page 详情页
	var DetailsUrl string = "https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" // 需要加上 作品 id

	var AllImgMap map[string]useType.Combination = map[string]useType.Combination{} // 所有的爬虫数据

	// 登陆等操作 s
	loginC := CreateCollector()
	//loginData := LoginData{}
	postKeyData := useType.PostKeyType{}
	loginC.OnHTML("#init-config", func(e *colly.HTMLElement) {
		fmt.Println(e.Request.URL, "url")
		el := e.DOM
		value, ok := el.Attr("value")
		if ok && len(value) > 10 {
			err := json.Unmarshal([]byte(value), &postKeyData)
			if err != nil {
				fmt.Println(err, "转json出错")
				os.Exit(1)
			} else {
				// success
			}
		}
	})

	err := loginC.Visit("https://accounts.pixiv.net/login")
	if err != nil {
		log.Fatalf("%s => 请求登陆页失败，检查网络或者代理", err)
	}

	logonMap := map[string]string{
		"pixiv_id":  config.UserName,
		"password":  config.Password,
		"post_key":  postKeyData.PostKey,
		"source":    "accounts",
		"return_to": "https://www.pixiv.net/",
	}
	err = loginC.Post("https://accounts.pixiv.net/api/login?lang=zh", logonMap)
	if err != nil {
		fmt.Println(err, "---登陆出错---")
		os.Exit(1)
	}

	// 登陆等操作 e

	// 开始爬虫 开始

	// 生成 详情 爬取 Collector s-----
	detailsC := loginC.Clone()
	detailsC.Async = true
	detailsC.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: config.MaxGoroutine}) // 设置最大并发

	detailsC.OnHTML("head", func(e *colly.HTMLElement) {
		el := e.DOM
		el.Find("script").Each(func(i int, selection *goquery.Selection) {
			txt := selection.Text()
			targetData := strings.Index(txt, "globalInitData")
			if targetData > 0 {
				// 获取json数据
				index := strings.Index(txt, `"illustId"`)
				lastInde := strings.Index(txt, `"likeData"`)
				jsonStr := `{` + txt[index:lastInde-1] + `}`
				detailsInfo := useType.DetailsInfo{}
				err := json.Unmarshal([]byte(jsonStr), &detailsInfo)
				if err != nil {
					fmt.Println(err, "=> 解析详情数据失败")
				} else {
					// 入库操作
					mutex.Lock()
					b := AllImgMap[detailsInfo.IllustId]
					AllImgMap[detailsInfo.IllustId] = useType.Combination{
						BaseDataItem: useType.BaseDataItem{
							IllustId:      b.IllustId,
							IllustTitle:   b.IllustTitle,
							Url:           b.Url,
							UserId:        b.UserId,
							UserName:      b.UserName,
							UserImage:     b.UserImage,
							Width:         b.Width,
							Height:        b.Height,
							BookmarkCount: b.BookmarkCount,
						},
						Urls: detailsInfo.Urls,
					}
					mutex.Unlock()
				}
			}
		})
	})
	// 生成 详情 爬取 Collector e----

	// 生成page爬取 Collector s---
	pageC := loginC.Clone()

	// 页面数据处理
	pageC.OnHTML("#js-mount-point-search-result-list", func(e *colly.HTMLElement) {
		//fmt.Println("页面数据处理", time.Now().Format("2006-01-02 15:04:05"))
		pageDataStr, ok := e.DOM.Attr("data-items")
		if ok {
			// 解析 数据
			baseDataList := []useType.BaseDataItem{}
			err := json.Unmarshal([]byte(pageDataStr), &baseDataList)
			if err != nil {
				log.Fatalf("%s => 解析page列表数据失败")
			}
			// 遍历 判断
			for _, item := range baseDataList {
				// 如果收藏数大于 最小收藏， 请求详情页
				if item.BookmarkCount >= config.MinCollection {
					// 存入 到map
					mutex.Lock()
					AllImgMap[item.IllustId] = useType.Combination{
						BaseDataItem: item,
					}
					mutex.Unlock()
					err := detailsC.Visit(DetailsUrl + item.IllustId)
					if err != nil {
						fmt.Println("%s => 请求图片详情页面错误", err)
						fmt.Println("正在尝试重复请求...", )
						err := detailsC.Visit(DetailsUrl + item.IllustId)
						if err != nil {
							log.Fatalf("重复请求失败， 请检查网络或者代理")
						}
					}
				}
			}
		} else {
			// 爬取完毕
		}
	})

	// 下一页 处理
	pageC.OnHTML(".next", func(e *colly.HTMLElement) {
		detailsC.Wait() // 等待爬取完毕
		url, ok := e.DOM.Children().Attr("href")
		//fmt.Println(url, time.Now().Format("2006-01-02 15:04:05"))
		if ok {
			abUrl := e.Request.AbsoluteURL(url)
			fmt.Println("当前爬取路径： ", abUrl)
			err := e.Request.Visit(abUrl)
			if err != nil {
				fmt.Println("下一页请求失败", err)
				fmt.Println("正尝试再次请求...")
				err := pageC.Visit(abUrl)
				if err != nil {
					log.Fatalf("尝试再次请求失败...请检查网络或者代理 %v", err)
				}
			}
		}
	})

	fmt.Println("当前爬取路径： ", CrawlUrl)
	err = pageC.Visit(CrawlUrl)
	if err != nil {
		fmt.Printf("%s => 爬取错误 \n", err)
		fmt.Println("正在尝试重新请求")
		err = pageC.Visit(CrawlUrl)
		if err != nil {
			log.Fatalf("再次请求失败=》%v", err)
		}
	}
	detailsC.Wait() // 等待爬取完毕

	// 将map 转为 切片 方便排序
	ArrDataSclice := useType.CombinationSlice{}
	for _, v := range AllImgMap {
		ArrDataSclice = append(ArrDataSclice, v)
	}

	sort.Sort(ArrDataSclice)
	// 排序完成

	// 将数据变为json 格式
	bt, err := json.Marshal(ArrDataSclice)
	if err != nil {
		fmt.Println(err, "=> 爬取的数据转json失败")
	}
	// 将信息写入json 文件
	fileName := keyWord + time.Now().Format("2006-01-02") + fmt.Sprintf("共%d条数据", len(ArrDataSclice)) + ".json"

	f, err := os.OpenFile("./"+fileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err, "=> 打开文件失败")
	}
	defer f.Close()
	_, err = f.Write(bt)
	if err != nil {
		fmt.Println(err, "=> 爬取数据写入失败")
	}

	// 图片处理 s
	fmt.Println("是否需要下载图片（y/n）：")
	downloadImgTxt := ""
	_, err = fmt.Scanf("%v", &downloadImgTxt)
	if err != nil {
		fmt.Println("读取扫描失败=>", err)
	}
	if downloadImgTxt == "y" || downloadImgTxt == "Y" {
		fmt.Printf("本次共爬取%d条数据， 请输入您需要下载的数量（按收藏排序下载）：\n", len(ArrDataSclice))
		imgNumber := 0
		imgNumberTxt := ""
		for {
			fmt.Scanf("%v", &imgNumberTxt)
			imgNumber, err = strconv.Atoi(imgNumberTxt)
			if err != nil {
				fmt.Println("您输入的不是数字， 请重新输入")
			} else {
				break
			}
		}
		if imgNumber > len(ArrDataSclice) {
			imgNumber = len(ArrDataSclice)
		}
		dowLenSclice := ArrDataSclice[0:imgNumber]
		// 保存的文件夹
		dir := "./" + keyWord + time.Now().Format("2006-01-02")
		_, err = utils.NoPathOnCreate(dir)
		if err != nil {
			log.Fatalf("创建文件夹失败=》%v", err)
		}

		// 创建通道，用于多线程下载
		var downloadChan chan string = make(chan string, config.MaxGoroutine) // 控制最大线程
		var wg sync.WaitGroup // 保证线程执行完毕

		// 对要下载的图片下载 s
		for index, item := range dowLenSclice {
			downloadChan <- ""
			wg.Add(1)

			url := item.Urls["original"]
			imgName, err := utils.UrlGetFileName(url)
			if err != nil {
				imgName = url
			}
			imgName = fmt.Sprintf("%d_%s", index+1, imgName)
			go func() {
				fmt.Println("imgName",imgName)
				err = download.DownloadAndSave(url, dir, imgName)
				if err != nil {
					fmt.Println("下载并且保存图片错误=》%v", err)
				}
				defer func() {
					<-downloadChan
					wg.Done()
				}()
			}()
		}

		wg.Wait()
		// 对要下载的图片下载 s
	} else {
		// pass
	}

	// 图片处理 e

	fmt.Println("爬取完毕")
}

// 创建 基础Collector
func CreateCollector() *colly.Collector {
	c := colly.NewCollector()
	if len(config.ProxyList) != 0 {
		if p, err := proxy.RoundRobinProxySwitcher(config.ProxyList...); err == nil {
			c.SetProxyFunc(p)
		} else {
			log.Fatalf("代理错误 %v", err)
		}
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "zh-CN,zh-TW;q=0.9,zh;q=0.8,en;q=0.7")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36")
		r.Headers.Set("Referer", "https://www.google.com/recaptcha/api2/anchor?ar=1&k=6LfJ0Z0UAAAAANqP-8mvUln2z6mHJwuv5YGtC8xp&co=aHR0cHM6Ly9hY2NvdW50cy5waXhpdi5uZXQ6NDQz&hl=zh-CN&v=v1561357937155&size=invisible&cb=apv8u3bpvdmf")
	})
	return c
}
