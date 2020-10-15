package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/htruong/go-md2"
	"golang.org/x/crypto/md4"
)

var startTime int64
var mHashMask map[int]string
var szLowercase, szUppercase, szDigits, szHexdigits, szPunctuation, szPrintable string
var txt, dic, diyDic, finalDic string
var bIsRandTxt bool
var iLenMd5, iCryptoMode, iTotal, iShown int
var bFinalDic strings.Builder
var np func() string

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

// 将输入的"??6377????666"解析为{2:"6377";10:"666"}格式(起始值0)字典方便匹配调用
func parseHashMask(word string) map[int]string {
	hashMask := make(map[int]string)
	reg := regexp.MustCompile(`[^?]+`)
	allStr := reg.FindAllString(word, -1)
	allStrIndex := reg.FindAllStringIndex(word, -1)
	for i, v := range allStrIndex {
		hashMask[v[0]] = allStr[i]
	}
	return hashMask
}

// 将输入的"3:6377|11:666"(其实值1)解析为{2:"6377";10:"666"}格式字典方便匹配调用
func parseHashMaskSep(word string) map[int]string {
	hashMask := make(map[int]string)
	for _, v := range strings.Split(word, "|") {
		flag := strings.Split(v, ":")
		pos, _ := strconv.Atoi(flag[0])
		if pos > 0 {
			hashMask[pos-1] = flag[1]
		}
	}
	return hashMask
}

// 将输入的"*6377*"解析为{-1:"6377"}格式(故意设置key为-1，以区分上面另外两种情况)
func parseHashMaskStar(word string) map[int]string {
	hashMask := make(map[int]string)
	hashMask[-1] = strings.Replace(word, "*", "", -1)
	return hashMask
}

// 笛卡尔乘积，同python中itertools.product
func nextPassword(n int, c string) func() string {
	r := []rune(c)
	p := make([]rune, n)
	x := make([]int, len(p))
	return func() string {
		//p := p[:len(x)]
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

var asciiBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	asciiIdxBits = 6
	asciiIdxMask = 1<<asciiIdxBits - 1
	asciiIdxMax  = 63 / asciiIdxBits
)

// 获取长度为n的随机字符串
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), asciiIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), asciiIdxMax
		}
		if idx := int(cache & asciiIdxMask); idx < len(asciiBytes) {
			b[i] = asciiBytes[idx]
			i--
		}
		cache >>= asciiIdxBits
		remain--
	}
	return string(b)
}

// 获取字符串MD2值
func GetMD2(data string) string {
	h := md2.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取字符串的MD4值
func GetMD4(data string) string {
	h := md4.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取字符串的32位MD5值
func Get32MD5(data string) string {
	sum := md5.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

// 获取字符串的16位MD5值
func Get16MD5(data string) string {
	return Get32MD5(data)[8:24]
}

// 获取字符串的SHA1值
func GetSha1(data string) string {
	sum := sha1.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

// 获取字符串的SHA224值
func GetSha224(data string) string {
	sum := sha256.Sum224([]byte(data))
	return hex.EncodeToString(sum[:])
}

// 获取字符串的SHA256值
func GetSha256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

// 获取字符串的SHA384值
func GetSha384(data string) string {
	sum := sha512.Sum384([]byte(data))
	return hex.EncodeToString(sum[:])
}

// 获取字符串的SHA512值
func GetSha512(data string) string {
	sum := sha512.Sum512([]byte(data))
	return hex.EncodeToString(sum[:])
}

func produce(pwd string, p chan<- string) {
	if bIsRandTxt {
		p <- RandStringBytesMaskImpr(rand.Intn(30) + 1)
	} else {
		p <- genTxt(txt, pwd)
	}
}

func routine(c <-chan string) {
	var szhash string
	isMatch := false

	dstTxt := <-c
	if iCryptoMode == 0 {
		szhash = GetMD2(dstTxt)
	} else if iCryptoMode == 1 {
		szhash = GetMD4(dstTxt)
	} else if iCryptoMode == 2 {
		if iLenMd5 == 32 {
			szhash = Get32MD5(dstTxt)
		} else {
			szhash = Get16MD5(dstTxt)
		}
	} else if iCryptoMode == 3 {
		szhash = GetSha1(dstTxt)
	} else if iCryptoMode == 4 {
		szhash = GetSha224(dstTxt)
	} else if iCryptoMode == 5 {
		szhash = GetSha256(dstTxt)
	} else if iCryptoMode == 6 {
		szhash = GetSha384(dstTxt)
	} else {
		szhash = GetSha512(dstTxt)
	}

	for k, v := range mHashMask {
		if k == -1 {
			isMatch = strings.Contains(szhash, v)
		} else {
			isMatch = strings.HasPrefix(szhash[k:], v)
		}
		if !isMatch {
			break
		}
	}

	if isMatch {
		fmt.Printf("Bingo!! It's your goal : %s  %s\n", dstTxt, szhash)
		if bIsRandTxt {
			if iShown < iTotal {
				iShown++
			} else {
				os.Exit(3)
			}
		} else {
			fmt.Printf("Time escaped : %d ms\n", (time.Now().UnixNano()-startTime)/1000000)
			os.Exit(3)
		}
	}
	return
}

func main() {
	var pwd, hashMask string
	var bShowVersion bool
	iShown = 1
	startTime = time.Now().UnixNano()
	szLowercase = "abcdefghijklmnopqrstuvwxyz"
	szUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	szDigits = "1234567890"
	szHexdigits = "1234567890abcdefABCDEF"
	szPunctuation = "@$_&-!\"#%()*+,./:;<=>?[\\]^`{|}~ "
	szPrintable = szDigits + szLowercase + szUppercase + szPunctuation
	rand.Seed(time.Now().Unix())

	flag.StringVar(&txt, "a", "", "设置明文格式，支持?占位符，如flag{?????}(Linux下字符串请使用引号包裹)")
	flag.BoolVar(&bIsRandTxt, "aa", false, "不限制明文，随机穷举指定格式HASH")
	flag.StringVar(&dic, "b", "", "按顺序组合穷举字符集(顺序会严重影响穷举速度，请尽量精确)\nd 数字 | l 小写字母 | u 大写字母 | h 十六进制字符集 | p 特殊字符 | r 可见字符\n例如：穷举字符集为数字、字母 -b=dlu")
	flag.StringVar(&diyDic, "bb", "", "自定义穷举字符集")
	flag.IntVar(&iCryptoMode, "m", 2, "设置HASH算法\n0 MD2 | 1 MD4 | 2 MD5 | 3 SHA1 | 4 SHA224 | 5 SHA256 | 6 SHA384 | 7 SHA512")
	flag.StringVar(&hashMask, "s", "", "设置HASH值字符串格式，支持3种模式\n? 占位符模式，如HASH第3位开始是6377，直接写'??6377'即可\n| 分隔符模式，如HASH第3位开始是6377第11位开始是66，直接写'3:6377|11:66'即可\n* 通配符模式，如fuzz包含7366的hash值，直接写'*7366*'即可")
	flag.IntVar(&iLenMd5, "i", 32, "设置目标MD5位数16位或32位")
	flag.IntVar(&iTotal, "t", 3, "使用-aa选项随机穷举HASH时，设置最少输出条数")
	flag.BoolVar(&bShowVersion, "v", false, "显示当前版本号")
	// 必须在所有flag都注册好而未访问其值时执行
	flag.Parse()

	if bShowVersion {
		fmt.Println("Version : 1.3.3")
		os.Exit(3)
	}

	if len(dic) > 0 {
		for _, v := range strings.ToLower(dic) {
			switch {
			case v == 'l':
				bFinalDic.WriteString(szLowercase)
			case v == 'u':
				bFinalDic.WriteString(szUppercase)
			case v == 'd':
				bFinalDic.WriteString(szDigits)
			case v == 'h':
				bFinalDic.WriteString(szHexdigits)
			case v == 'p':
				bFinalDic.WriteString(szPunctuation)
			case v == 'r':
				bFinalDic.WriteString(szPrintable)
			}
		}
		finalDic = removeDuplicate(bFinalDic.String())
	} else if len(diyDic) > 0 {
		finalDic = removeDuplicate(diyDic)
	}

	i := strings.Count(txt, "?")
	if i > 0 {
		np = nextPassword(i, finalDic)
	} else if len(txt) > 0 {
		fmt.Printf("MD2     : %s\n", GetMD2(txt))
		fmt.Printf("MD4     : %s\n", GetMD4(txt))
		fmt.Printf("MD5(16) : %s\n", Get16MD5(txt))
		fmt.Printf("MD5(32) : %s\n", Get32MD5(txt))
		fmt.Printf("SHA1    : %s\n", GetSha1(txt))
		fmt.Printf("SHA224  : %s\n", GetSha224(txt))
		fmt.Printf("SHA256  : %s\n", GetSha256(txt))
		fmt.Printf("SHA384  : %s\n", GetSha384(txt))
		fmt.Printf("SHA512  : %s\n", GetSha512(txt))
		os.Exit(3)
	}

	if len(txt)*len(hashMask)*(len(dic)+len(diyDic)) == 0 {
		if !(len(hashMask) != 0 && bIsRandTxt) {
			fmt.Println(`
  未设置必要参数，查看帮助 bruteHASH -h
  示例：
    直接输出"HelloWorld"字符串的多种HASH值
      > bruteHASH -a=HelloWorld
    随机字符穷举，输出至少6条hash开头是"6377"的SHA1
      > bruteHASH -aa -s=6377 -m=3 -t=6
    限制数字穷举，hash第7位是"6377"的SHA256
      > bruteHASH -aa -b=d -s="??????6377" -m=5
      > bruteHASH -aa -b=d -s="7:6377" -m=5
    随机字符穷举，hash第3位是"63"第11位是"77"的SHA224
      > bruteHASH -aa -s="??63??????77" -m=4
      > bruteHASH -aa -s="3:63|11:77" -m=4
    随机字符穷举，hash包含"6377"的md4
      > bruteHASH -aa -s="*6377*" -m=1
    自定义字符集穷举"c???new???"明文，以"95ce2a"结尾的16位MD5
      > bruteHASH -a="c???new???" -bb=abcdefnutvw_ -s="??????????95ce2a" -i=16
      > bruteHASH -a="c???new???" -bb=abcdefnutvw_ -s="11:95ce2a" -i=16
			`)
			os.Exit(3)
		}
	}

	if strings.Index(hashMask, ":") > 0 {
		mHashMask = parseHashMaskSep(strings.ToLower(hashMask))
	} else if strings.Index(hashMask, "*") >= 0 {
		mHashMask = parseHashMaskStar(strings.ToLower(hashMask))
	} else {
		mHashMask = parseHashMask(strings.ToLower(hashMask))
	}

	if len(finalDic) > 0 {
		if bIsRandTxt {
			asciiBytes = finalDic
		}
		fmt.Println("Brute-force range : " + finalDic)
	} else {
		fmt.Println("Brute-force range : " + asciiBytes)
	}

	iChanNum := 0
	for {
		if bIsRandTxt {
			iChanNum = 100
		} else {
			pwd = np()
			if len(pwd) == 0 {
				break
			}
		}
		ch := make(chan string, iChanNum)
		go produce(pwd, ch)
		go routine(ch)
	}
}
