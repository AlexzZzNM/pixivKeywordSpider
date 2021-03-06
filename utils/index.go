package utils

import (
	"errors"
	"os"
	"pixivKwywordSpider/useType"
	"strings"
)

// url 获取文件名
func UrlGetFileName(url string) (string, error) {
	lastIndex := strings.LastIndex(url, "/")
	lastQueryIndex := strings.LastIndex(url, "?")
	if lastIndex == -1 {
		return "", errors.New("无法找到'/'字符")
	}

	if lastQueryIndex == -1 {
		return string([]rune(url)[lastIndex+1:]), nil
	} else {
		return string([]rune(url)[lastIndex:lastQueryIndex]), nil
	}
}

func UrlGetPath(url string) (string) {
	lastIndex := strings.LastIndex(url, "/")
	if lastIndex == -1 {
		return ""
	}
	return string([]rune(url)[0:lastIndex])
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 如果路径不存在就创建
func NoPathOnCreate(dir string) (bool, error) {
	exist, err := PathExists(dir)
	if err != nil {
		return false, err
	}

	if exist {
		// pass
	} else {
		// 创建文件夹
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// 获取map key 集合
func GetMaPkey(newMap map[string]useType.DetailsInfoID) []string {
	arr := []string{}
	for key := range newMap{
		arr = append(arr, key)
	}
	return arr
}
