package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// 字符去重
func removeDuplicate(txt string) string {
	var distinctStr strings.Builder
	tmpMap := make(map[rune]interface{})
	for _, val := range txt {
		if _, ok := tmpMap[val]; !ok {
			distinctStr.WriteRune(val)
			tmpMap[val] = nil
		}
	}
	return distinctStr.String()
}

// 填充替代?占位符，返回填充好的字符串
func genTxt(txt string, dic string) string {
	re := regexp.MustCompile(`\?+`)
	strIndex := re.FindAllStringIndex(txt, -1)
	for _, v := range strIndex {
		tmpLen := v[1] - v[0]
		txt = txt[:v[0]] + dic[:tmpLen] + txt[v[1]:]
		dic = dic[tmpLen:]
	}
	return txt
}

// 笛卡尔乘积，同python中itertools.product
func nextPassword(n int, c string) func() string {
	r := []rune(c)
	p := make([]rune, n)
	x := make([]int, len(p))
	return func() string {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = r[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(r) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return string(p)
	}
}

func Get32MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Get16MD5Encode(data string) string {
	return Get32MD5Encode(data)[8:24]
}

func routine(txt string, pwd string, lenMD5 int, verbose bool, startwith string, endwith string, instr string, startTime int64) {
	var md5Str string
	var isMatch bool
	dstTxt := genTxt(txt, pwd)
	if lenMD5 == 32 {
		md5Str = Get32MD5Encode(dstTxt)
	} else {
		md5Str = Get16MD5Encode(dstTxt)
	}
	if verbose {
		fmt.Println("Trying : " + dstTxt + "  " + md5Str)
	}

	if len(startwith) > 0 {
		isMatch = strings.HasPrefix(md5Str, startwith)
	}
	if len(endwith) > 0 {
		isMatch = strings.HasSuffix(md5Str, endwith)
	}
	if len(instr) > 0 {
		isMatch = strings.Contains(md5Str, instr)
	}
	if isMatch {
		fmt.Printf("Bingo!! Here is what you want : %s  %s\n", dstTxt, md5Str)
		fmt.Printf("Time escaped : %d ms", (time.Now().UnixNano()-startTime)/1000000)
		os.Exit(3)
	}
}

func main() {
	startTime := time.Now().UnixNano()
	var lowercase string = "abcdefghijklmnopqrstuvwxyz"
	var uppercase string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var digits string = "1234567890"
	var hexdigits string = "1234567890abcdefABCDEF"
	var punctuation string = "!\"#$%&()*+,-./:;<=>?@[\\]^_`{|}~ "
	var printable string = digits + lowercase + uppercase + punctuation
	var verbose bool
	var lenMD5 int
	var txt, dstTxt, startwith, endwith, instr, dic, diyDic, finalDic string
	var bFinalDic strings.Builder
	var np func() string

	flag.StringVar(&txt, "a", "", "设置明文格式，支持?占位符，如flag{?????}")
	flag.StringVar(&dic, "b", "", "按顺序组合爆破字符集(字符集先后顺序会严重影响爆破速度，请尽量精确)\n数字d | 小写字母l | 大写字母u | 16进制字符集h | 特殊字符p | 所有可见字符r\n例如：指定爆破字符集为数字、字母 -b=dlu")
	flag.StringVar(&diyDic, "bb", "", "自定义爆破字符集")
	flag.StringVar(&startwith, "s", "", "设置目标MD5值起始字符串")
	flag.StringVar(&endwith, "e", "", "设置目标MD5值结束字符串")
	flag.StringVar(&instr, "c", "", "设置目标MD5值包含字符串")
	flag.IntVar(&lenMD5, "i", 32, "设置目标MD5位数16位或32位")
	flag.BoolVar(&verbose, "v", false, "显示爆破进度(影响爆破速度)")
	// 必须在所有flag都注册好而未访问其值时执行
	flag.Parse()

	if len(txt)*(len(startwith)+len(endwith)+len(instr))*(len(dic)+len(diyDic)) == 0 {
		fmt.Println(`
  未设置必要参数，查看帮助 bruteMD5 -h
  示例：
    用自定义字符集穷举"code??{q????w}"明文，32位MD5结尾为"930bac91"
      > bruteMD5 -a=code??{q????w} -bb=ABCcopqrstuvwxyz_ -e=930bac91
    用自定义字符集穷举"c???new???"明文，32位MD5包含字符串"3b605234ed"
      > bruteMD5 -a=c???new??? -bb=abcdefnutuvw_ -c=3b605234ed
    用数字、大写字母穷举明文"flag{?????}"(?代表未知5位)，16位MD5开头为"b6dff925"
      > bruteMD5 -a=flag{?????} -b=du -s=b6dff925 -i=16
		`)
		os.Exit(3)
	}

	if len(dic) > 0 {
		for _, v := range strings.ToLower(dic) {
			switch {
			case v == 'l':
				bFinalDic.WriteString(lowercase)
			case v == 'u':
				bFinalDic.WriteString(uppercase)
			case v == 'd':
				bFinalDic.WriteString(digits)
			case v == 'h':
				bFinalDic.WriteString(hexdigits)
			case v == 'p':
				bFinalDic.WriteString(punctuation)
			case v == 'r':
				bFinalDic.WriteString(printable)
			}
		}
		finalDic = removeDuplicate(bFinalDic.String())
	} else if len(diyDic) > 0 {
		finalDic = removeDuplicate(diyDic)
	}

	fmt.Println("Brute-force range : " + finalDic)

	i := strings.Count(txt, "?")
	if i > 0 {
		np = nextPassword(i, finalDic)
	} else {
		fmt.Printf("Your plaintext and MD5 is : %s  %s", txt, Get32MD5Encode(dstTxt))
		os.Exit(3)
	}

	for {
		pwd := np()
		if len(pwd) == 0 {
			break
		}
		go routine(txt, pwd, lenMD5, verbose, startwith, endwith, instr, startTime)
	}
}
