# bruteHASH  

## 别问，问就是为了CTF  

### 功能  
随机或穷举指定格式HASH值，输出符合条件的"明文 HASH"  

明文/HASH都支持使用占位符"?"  
支持随机字符/自定义字符集穷举  
支持常见HASH类型(MD4/MD5/SHA1/SHA224/SHA256/SHA384/SHA512)  

### 帮助  
```
Usage of bruteHASH v1.3.2:
  -a string
        设置明文格式，支持?占位符，如flag{?????}(Linux下字符串请使用引号包裹)
  -aa
        不限制明文，随机穷举指定格式HASH
  -b string
        按顺序组合穷举字符集(顺序会严重影响穷举速度，请尽量精确)
        d 数字 | l 小写字母 | u 大写字母 | h 十六进制字符集 | p 特殊字符 | r 可见字符
        例如：穷举字符集为数字、字母 -b=dlu
  -bb string
        自定义穷举字符集
  -i int
        设置目标MD5位数16位或32位 (default 32)
  -m int
        设置HASH算法
        0 MD4 | 1 MD5 | 2 SHA1 | 3 SHA224 | 4 SHA256 | 5 SHA384 | 6 SHA512 (default 1)
  -s string
        设置HASH值字符串格式，支持3种模式
        ? 占位符模式，如HASH第3位开始是6377，直接写'??6377'即可
        | 分隔符模式，如HASH第3位开始是6377第11位开始是66，直接写'3:6377|11:66'即可
        * 通配符模式，如fuzz包含7366的hash值，直接写'*7366*'即可
  -t int
        使用-aa选项随机穷举HASH时，设置最少输出条数 (default 3)
  -v    显示当前版本号
```  

![bruteHASH帮助](https://github.com/playGitboy/bruteHASH/blob/master/img/bruteHASH_help.jpg)  

### 示例  
```
直接输出"HelloWorld"字符串的多种HASH值
  > bruteHASH -a=HelloWorld
随机字符穷举，输出至少6条hash开头是"6377"的SHA1
  > bruteHASH -aa -s=6377 -m=2 -t=6
限制数字穷举，hash第7位是"6377"的SHA256
  > bruteHASH -aa -b=d -s="??????6377" -m=4
  > bruteHASH -aa -b=d -s="7:6377" -m=4
随机字符穷举，hash第3位是"63"第11位是"77"的SHA224
  > bruteHASH -aa -s="??63??????77" -m=3
  > bruteHASH -aa -s="3:63|11:77" -m=3
随机字符穷举，hash包含"6377"的md4
  > bruteHASH -aa -s="*6377*" -m=0
自定义字符集穷举"c???new???"明文，以"95ce2a"结尾的16位MD5
  > bruteHASH -a="c???new???" -bb=abcdefnutvw_ -s="??????????95ce2a" -i=16
  > bruteHASH -a="c???new???" -bb=abcdefnutvw_ -s="11:95ce2a" -i=16
```  

![bruteHASH测试](https://github.com/playGitboy/bruteHASH/blob/master/img/bruteHASH_test.jpg)  

### Fuzz特殊HASH  
使用该工具Fuzz出一些CTF常见特殊HASH，有备无患(￣▽￣)"  
* 明文和md5都以0e开头   
0e215962017  0eBkcqQpv  0eKfoob  0edpGW  0embO4G  0eqb  
* 明文和md4都以0e开头  
0e30  0e189  0e311  0e77961763272  0e001233333333333334557778889  
* 明文和sha1都以0e开头  
0ecJFe  0e6NM  0eYAu0dPt  
* md5包含"276f7227"  
ffifdyop  d0Fqvwtr2PitRUJyqT  hwqc5H27HdV6WhcBbKDVX  
> 用于web构造mysql注入md5($password,true)  

### 声明  
CTF偶尔要用"特殊"HASH，如MISC已知个别明文字符和部分HASH，穷举flag明文；WEB中构造MYSQL注入，要用指定字符集构造一个以"xxxxxxxx"开头的MD5等等。但找了半天，满天飞的都是"爆破"HASH的工具，一个好用的穷举生成HASH的工具都没有  

**“先从无到有，再从有到精”**  

**代码不精简可能有BUG且效率未达最佳，如欲吐槽请fork后show your code...**  
