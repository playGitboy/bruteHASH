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
	"strings"
	"time"

	"golang.org/x/crypto/md4"
)

var startTime int64
var szLowercase, szUppercase, szDigits, szHexdigits, szPunctuation, szPrintable string
var txt, startwith, endwith, instr, dic, diyDic, finalDic string
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
		szhash = GetMD4(dstTxt)
	} else if iCryptoMode == 1 {
		if iLenMd5 == 32 {
			szhash = Get32MD5(dstTxt)
		} else {
			szhash = Get16MD5(dstTxt)
		}
	} else {
		szhash = GetSha1(dstTxt)
	}

	if len(startwith) > 0 {
		isMatch = strings.HasPrefix(szhash, startwith)
	}
	if len(endwith) > 0 {
		isMatch = strings.HasSuffix(szhash, endwith)
	}
	if len(instr) > 0 {
		isMatch = strings.Contains(szhash, instr)
	}
	if isMatch {
		fmt.Printf("Bingo!! Here is what you want : %s  %s\n", dstTxt, szhash)
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
	var pwd, szhash string
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
	flag.StringVar(&dic, "b", "", "按顺序组合穷举字符集(字符集顺序会严重影响穷举速度，请尽量精确)\nd 数字 | l 小写字母 | u 大写字母 | h 十六进制字符集 | p 特殊字符 | r 可见字符\n例如：指定穷举字符集为数字、字母 -b=dlu")
	flag.StringVar(&diyDic, "bb", "", "自定义穷举字符集")
	flag.IntVar(&iCryptoMode, "m", 1, "设置HASH算法\n0 MD4 | 1 MD5 | 2 SHA1 | 3 SHA224 | 4 SHA256 | 5 SHA384 | 6 SHA512")
	flag.StringVar(&startwith, "s", "", "设置目标HASH值起始字符串")
	flag.StringVar(&endwith, "e", "", "设置目标HASH值结束字符串")
	flag.StringVar(&instr, "c", "", "设置目标HASH值包含字符串")
	flag.IntVar(&iLenMd5, "i", 32, "设置目标MD5位数16位或32位")
	flag.IntVar(&iTotal, "t", 3, "使用-aa选项随机穷举HASH时，设置最少输出条数")
	// 必须在所有flag都注册好而未访问其值时执行
	flag.Parse()

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
		if iCryptoMode == 0 {
			szhash = GetMD4(txt)
		} else if iCryptoMode == 1 {
			if iLenMd5 == 32 {
				szhash = Get32MD5(txt)
			} else {
				szhash = Get16MD5(txt)
			}
		} else if iCryptoMode == 2 {
			szhash = GetSha1(txt)
		} else if iCryptoMode == 3 {
			szhash = GetSha224(txt)
		} else if iCryptoMode == 4 {
			szhash = GetSha256(txt)
		} else if iCryptoMode == 5 {
			szhash = GetSha384(txt)
		} else {
			szhash = GetSha512(txt)
		}
		fmt.Printf("Your plaintext and hash is : %s  %s", txt, szhash)
		os.Exit(3)
	}

	if len(txt)*(len(startwith)+len(endwith)+len(instr))*(len(dic)+len(diyDic)) == 0 {
		if !((len(startwith)+len(endwith)+len(instr) != 0) && bIsRandTxt) {
			fmt.Println(`
  未设置必要参数，查看帮助 bruteHASH -h
  示例：
    随机字符穷举，HASH中包含"6377666"的SHA1
      > bruteHASH -aa -c=6377666 -m=2
    随机字符穷举，"0e"开头的MD4
      > bruteHASH -aa -s=0e -m=0
    用自定义字符集穷举"c???new???"明文，32位MD5包含字符串"3b605234ed"
      > bruteHASH -a="c???new???" -bb=abcdefnutuvw_ -c=3b605234ed
    用数字、大写字母穷举明文"flag{?????}"(?代表未知5位)，16位MD5开头为"b6dff925"
      > bruteHASH -a="flag{?????}" -b=du -s=b6dff925 -i=16
			`)
			os.Exit(3)
		}
	}

	startwith = strings.ToLower(startwith)
	endwith = strings.ToLower(endwith)
	instr = strings.ToLower(instr)

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
