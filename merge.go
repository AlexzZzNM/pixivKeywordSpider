package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"pixivKwywordSpider/useType"
	"sort"
)

func main()  {
	dir := ""
	arr := []string{}
	outFile := ""

	targetArr := []useType.DetailsInfoTwo{}


	for true {
		fmt.Println("请输入爬取生成的文件路径和文件名, 使用空格隔开（例：./xxx.json）, null 结束输入：")
		fmt.Scanln(&dir)
		if dir == "null" {
			break
		} else {
			arr = append(arr, dir)
		}
	}

	fmt.Println("请输入输出文件夹")
	fmt.Scanln(&outFile)

	for _, v := range arr{
		if v != " " {
			file, err := os.Open(v)
			if err != nil {
				log.Fatalf("打开文件失败=>%v", err)
			}

			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalf("读取文件失败 %v", err)
			}

			var ArrDataSclice  = []useType.DetailsInfoTwo{}
			err = json.Unmarshal(bytes, &ArrDataSclice)
			if err != nil {
				log.Fatalf("解析文本数据失败 %v", err)
			}

			targetArr = append(targetArr, ArrDataSclice...)
		}
	}

	targetArr = duplicateRemoval(targetArr)

	// 排序
	SortAllListTwo(targetArr)

	f, err := os.OpenFile("./"+outFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err, "=> 打开文件失败")
	}
	defer f.Close()
	targetBts, err := json.Marshal(targetArr)
	if err != nil{
		log.Fatalf("将数据转为字节失败 %v", err)
	}

	_, err = f.Write(targetBts)
	if err != nil {
		fmt.Println(err, "=> 总数据写入失败")
	}

}


func SortAllListTwo(list []useType.DetailsInfoTwo) {
	sort.Slice(list, func(i, j int) bool {
		keys1 := getMaPkeyTow(list[i].Illust)
		keys2 := getMaPkeyTow(list[j].Illust)
		if len(keys1) == 0 {
			return false
		}
		if len(keys2) == 0 {
			return true
		}

		k1 := keys1[0]
		k2 := keys2[0]
		return list[i].Illust[k1].BookmarkCount > list[j].Illust[k2].BookmarkCount
	})
}

func getMaPkeyTow(newMap map[string]useType.DetailsInfoID) []string {
	arr := []string{}
	for key := range newMap{
		arr = append(arr, key)
	}
	return arr
}

// 去重重复数据, 以及空白的数据
func duplicateRemoval(list []useType.DetailsInfoTwo) []useType.DetailsInfoTwo {
	idMap := make(map[string]bool)
	targetArr := []useType.DetailsInfoTwo{}

	for i, item := range list{
		keys := getMaPkeyTow(item.Illust)
		if len(keys) == 0 {
			fmt.Println(i)
		} else {
			if idMap[keys[0]] {
				fmt.Printf("index: %d  ids:%v \n", i, keys)
			} else {
				idMap[keys[0]] = true
				targetArr = append(targetArr, item)
			}

		}
	}

	return targetArr
}