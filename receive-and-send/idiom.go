package main

import "fmt"

func idiomExists(phrase string, idiom map[string]string) bool {
	_, ok := idiom[phrase]
	return ok == true && len([]rune(phrase)) == 4
}

func idiomIsRight(phrase string, pre string) bool {
	if phrase != "" && pre != "" {
		flag := string([]rune(phrase)[:1]) == string([]rune(pre)[3:])
		fmt.Println("flag = ", flag)
		return flag
	}
	return false
}

func getMeaning(phrase string, idiom map[string]string) string {
	return idiom[phrase]
}

func idiomSelect(phrase string, idiom map[string]string) string {
	last_word := string([]rune(phrase)[3:])
	for key, _ := range idiom {
		if string([]rune(key)[:1]) == last_word {
			return key
		}
	}
	return ""
}

func getRandKey(idiom map[string]string) string {
	index := 1
	for num := range idiom {
		if index == 2 {
			return num
		}
		index += 1
	}
	return ""
}
