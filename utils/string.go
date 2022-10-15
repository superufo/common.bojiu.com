package utils

import (
	"math/rand"
	"strings"
	"time"
)

var (
	charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// GenRandString 生成指定长度的随机字符串
func GenRandString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}

// IsEmptyString 为空判断
func IsEmptyString(text string) bool {
	s := strings.TrimSpace(text)
	if len(s) > 0 {
		return false
	}

	return true
}

// IsSpaceOrEmpty 判断是否空字符串
func IsSpaceOrEmpty(text string) bool {
	count := len(text)
	if count == 0 {
		return true
	}

	for i := 0; i < count; i++ {
		if text[i] != ' ' {
			return false
		}
	}
	return true
}

// DistinctStringList slice 去重
func DistinctStringList(strList []string) []string {
	distinctMap := make(map[string]bool, 0)
	var distinctList []string

	for _, str := range strList {
		if distinctMap[str] == true {
			continue
		} else {
			distinctMap[str] = true
			distinctList = append(distinctList, str)
		}
	}
	return distinctList
}

func ExistInString(strFind string, str string) bool {
	return strings.Contains(str, strFind)
}

// ExistInStringList in str array
func ExistInStringList(strFind string, strList []string) bool {
	flag := false
	for _, str := range strList {
		if str == strFind {
			flag = true
			break
		}
	}
	return flag
}

// RemoveStringList 移除str元素
func RemoveStringList(key string, strList []string) []string {
	removeList := make([]string, len(strList))
	removeList = strList[:]
	for i, str := range strList {
		if str == key {
			removeList = append(strList[:i], strList[i+1:]...)
			break
		}
	}
	return removeList
}

// RemoveDuplicate 移除重复元素
func RemoveDuplicate(duplicateList []string) (pureList []string) {
	strMap := make(map[string]bool, 0)
	pureList = make([]string, 0)
	for _, str := range duplicateList {
		if strMap[str] == true {
			continue
		} else {
			strMap[str] = true
			pureList = append(pureList, str)
		}
	}
	return pureList
}
