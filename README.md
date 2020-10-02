# bruteHASH

## 别问，问就是为了CTF

## 主要功能
随机或穷举指定格式HASH值，输出符合条件的"明文 HASH"  

支持指定明文格式  
不限定明文格式随机字符穷举  
自定义穷举字符集  
CTF常见HASH(MD4/MD5/SHA1)  
设置HASH开头、结尾或包含字符串  

### 示例
```
Usage of bruteHASH.exe:
  -a string
        设置明文格式，支持?占位符，如flag{?????}(Linux下字符串请使用引号包裹)
  -aa
        不限制明文，随机穷举指定格式HASH
  -b string
        按顺序组合穷举字符集(字符集顺序会严重影响爆破速度，请尽量精确)
        d 数字 | l 小写字母 | u 大写字母 | h 十六进制字符集 | p 特殊字符 | r 可见字符
        例如：指定爆破字符集为数字、字母 -b=dlu
  -bb string
        自定义穷举字符集
  -c string
        设置目标HASH值包含字符串
  -e string
        设置目标HASH值结束字符串
  -i int
        设置目标MD5位数16位或32位 (default 32)
  -m int
        设置HASH算法
        0 MD4 | 1 MD5 | 2 SHA1 (default 1)
  -s string
        设置目标HASH值起始字符串
  ```  

#### 具体用法 
随机字符穷举，HASH中包含"6377666"的SHA1  
> bruteHASH -aa -c=6377666 -m=2  

随机字符穷举，"0e"开头的MD4  
> bruteHASH -aa -s=0e -m=0  

用自定义字符集穷举"c???new???"明文，32位MD5包含字符串"3b605234ed"  
> bruteHASH -a="c???new???" -bb=abcdefnutuvw_ -c=3b605234ed  

用数字、大写字母穷举明文"flag{?????}"(?代表未知5位)，16位MD5开头为"b6dff925"  
> bruteHASH -a="flag{?????}" -b=du -s=b6dff925 -i=16  

![help](https://github.com/playGitboy/bruteMD5/tree/master/img/bruteMD5_help.png)  

![test](https://github.com/playGitboy/bruteMD5/tree/master/img/bruteMD5_test.png)  

## 声明  
CTF偶尔需要用到"特殊"HASH，比如MISC中已知个别明文字符和部分HASH，要穷举flag明文；WEB中构造MYSQL注入，要用指定字符集构造一个以"xxxxxxxx"开头的MD5等等。但找了半天，满天飞的都是"爆破"HASH的工具，一个好用的穷举生成HASH的工具都没有  

虽然"人生苦短，该用python"，但为了兼顾性能和开发效率，做了一个艰难的决定——  
学用golang试一试？  

首次用golang，本着能跑就行的初心聚合"云智慧"完成——  
**代码不精简有BUG且效率未达最佳，如需吐槽请fork后show me your code...**  
