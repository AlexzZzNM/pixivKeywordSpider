package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"pixivKwywordSpider/config"
	"pixivKwywordSpider/download"
	"pixivKwywordSpider/useType"
	"pixivKwywordSpider/utils"
	"strconv"
	"sync"
)

func main()  {
	fmt.Println("请输入爬取生成的文件路径和文件名（例：./xxx.json）：")
	dir := ""
	fmt.Scanf("%v", &dir)
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("打开文件失败=>%v", err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件失败 %v", err)
	}

	var ArrDataSclice []useType.Combination = []useType.Combination{}
	err = json.Unmarshal(bytes, &ArrDataSclice)
	if err != nil {
		log.Fatalf("解析文本数据失败 %v", err)
	}

	// 和 main 的一样了， 懒得写了
	fmt.Printf("本文件共%d条数据， 请输入您需要下载的数量：\n", len(ArrDataSclice))
	imgNumber := 0
	imgNumberTxt := ""
	for {
		fmt.Scanf("%s", &imgNumberTxt)
		fmt.Println(imgNumberTxt, "imgNumberTxt")
		imgNumber, err = strconv.Atoi(imgNumberTxt)
		if err != nil {
			fmt.Println(err, "err")
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
	saveDir := ""
	fmt.Println("请输入保存图片文件夹的路径路径（例：./save）：")
	fmt.Scanf("%v", &saveDir)
	_, err = utils.NoPathOnCreate(saveDir)
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
			err = download.DownloadAndSave(url, saveDir, imgName)
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

	fmt.Println("下载完成")
}