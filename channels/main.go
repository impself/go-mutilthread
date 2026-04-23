package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	windRegex     = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z]*(\d{5}KT|VRB\d{2}KT).*=`)
	trfValidation = regexp.MustCompile(`.*TAF.*`)
	comment       = regexp.MustCompile(`\w*#.*`)
	metarClose    = regexp.MustCompile(`.*=`)
	variableWind  = regexp.MustCompile(`.*VRB\d{2}KT`)
	validWind     = regexp.MustCompile(`\d{5}KT`)
	windDirOnly   = regexp.MustCompile(`(\d{3})\d{2}KT`)
	windDist      [8]int
)

// parseToArray 将原始METAR文本按行解析为独立的报文字符串切片
// 处理逻辑：逐行读取，跳过注释行(#开头)，拼接有效行内容，
// 遇到 "=" 闭合标记时将拼接结果作为一条完整报文加入切片，
// 遇到TAF行时停止解析（TAF是预报，不是观测报文，不需要处理）
func parseToArray(textChannel chan string, metarChannel chan []string) {
	for text := range textChannel {
		lines := strings.Split(text, "\n")
		// 预分配容量，避免频繁扩容
		metarSlice := make([]string, 0, len(lines))
		// 用于拼接多行为一条完整报文（METAR报文可能跨行）
		metarStr := ""
		for _, line := range lines {
			// 遇到TAF预报段，说明METAR部分已结束，停止解析
			if trfValidation.MatchString(line) {
				break
			}
			// 跳过注释行（以#开头的行），只拼接有效内容
			if !comment.MatchString(line) {
				metarStr += strings.Trim(line, " ")
			}
			// METAR报文以 "=" 结尾，遇到闭合标记说明一条报文拼接完毕
			if metarClose.MatchString(line) {
				metarSlice = append(metarSlice, metarStr)
				metarStr = "" // 重置，准备拼接下一条报文
			}
		}
		metarChannel <- metarSlice
	}
	close(metarChannel)
}

// extractWindDirection 从每条METAR报文中提取风向风速信息
// 使用 windRegex 匹配格式如 "25010KT"（前3位风向+2位风速）或 "VRB05KT"（风向不定）
// 返回第一个捕获组的内容，即风向风速子串（如 "25010KT" 或 "VRB05KT"）
func extractWindDirection(metarChannel chan []string, windsChannel chan []string) {
	for metars := range metarChannel {
		winds := make([]string, 0, len(metars))
		for _, metar := range metars {
			if windRegex.MatchString(metar) {
				// FindAllStringSubmatch 返回 [][]string，[0][1] 取第一条匹配的第一个捕获组
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}
		windsChannel <- winds
	}
	close(windsChannel)
}

// mineWindDistribution 将风向数据归类到8个方位桶中，统计风向频率分布
// windDist 数组的8个索引对应8个方位：
//
//	0=N(北), 1=NE(东北), 2=E(东), 3=SE(东南),
//	4=S(南), 5=SW(西南), 6=W(西), 7=NW(西北)
//
// 两种风向处理方式：
//   - VRB（风向不定）：8个方向全部+1，表示无法确定具体方向
//   - 固定风向（如250°）：角度/45取整后映射到0-7索引，对应方位桶+1
func mineWindDistribution(windsChannel chan []string, resultsChannel chan [8]int) {
	for winds := range windsChannel {
		for _, wind := range winds {
			if variableWind.MatchString(wind) {
				// 风向不定(VRB)，所有方位都计入
				for i := 0; i < 8; i++ {
					windDist[i]++
				}
			} else if validWind.MatchString(wind) {
				// 从 "25010KT" 格式中提取前3位风向角度，如 "250"
				windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]
				if d, err := strconv.ParseFloat(windStr, 64); err == nil {
					// 角度转方位索引：360° / 8方向 = 45°每方向
					// 例如 0°→0(N), 45°→1(NE), 90°→2(E) ... 250°→6(W偏SW)
					dirIndex := int(math.Round(d/45.0)) % 8
					windDist[dirIndex]++
				}
			}
		}
	}
	resultsChannel <- windDist
	close(resultsChannel)
}

func main() {

	textChannel := make(chan string)
	metarChannel := make(chan []string)
	windsChannel := make(chan []string)
	resultsChannel := make(chan [8]int)
	go parseToArray(textChannel, metarChannel)
	go extractWindDirection(metarChannel, windsChannel)
	go mineWindDistribution(windsChannel, resultsChannel)
	abspath, err := filepath.Abs("../metarfiles")
	if err != nil {
		log.Fatalf("wrong path: %v", err)
	}
	fmt.Println(abspath)
	files, _ := os.ReadDir(abspath)
	start := time.Now()
	for _, file := range files {
		dat, err := os.ReadFile(filepath.Join(abspath, file.Name()))
		if err != nil {
			panic(err)
		}
		text := string(dat)
		textChannel <- text
	}
	close(textChannel)
	result := <-resultsChannel
	elapsed := time.Since(start)
	fmt.Printf("%v\n", result)
	fmt.Printf("Processing took %s\n", elapsed) //没有channle花费2.8s ,使用channel花费1s
}
