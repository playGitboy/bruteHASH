# bruteMD5

## 别问，问就是为了CTF

## 声明
众所周知，CTF偶尔需要用到"特殊"MD5，比如MISC中已知个别字符和部分MD5，要穷举出flag明文；WEB中构造MYSQL注入时，要用指定字符集构造一个以"xxxxxxxx"开头的MD5等等。但找了半天，满天飞的都是爆破MD5的工具，一个好用的"穷举MD5"的工具都没有(如果hashcat能穷举指定格式的值多好)  

虽然"人生苦短，该用python"，但为了兼顾穷举性能和开发效率，于是做了一个艰难的决定——  
golang试一试？  

本人首次用golang，本着能跑就行的初心聚合"云智慧"完成——  
**代码不精简有BUG且效率未达最佳，如需吐槽请fork后show me your code...**  

### 示例
```
Usage of bruteMD5.exe:
  -a string
        设置明文格式，支持?占位符，如flag{?????}
  -b string
        按顺序组合爆破字符集(字符集先后顺序会严重影响爆破速度，请尽量精确)
        数字d | 小写字母l | 大写字母u | 16进制字符集h | 特殊字符p | 所有可见字符r
        例如：指定爆破字符集为数字、字母 -b=dlu
  -bb string
        自定义爆破字符集
  -c string
        设置目标MD5值包含字符串
  -e string
        设置目标MD5值结束字符串
  -i int
        设置目标MD5位数16位或32位 (default 32)
  -s string
        设置目标MD5值起始字符串
  -v    显示爆破进度(影响爆破速度)
  ```  

#### 具体用法 
用自定义字符集穷举"code??{q????w}"明文，32位MD5结尾为"930bac91"  
> bruteMD5 -a=code??{q????w} -bb=ABCcopqrstuvwxyz_ -e=930bac91  

用自定义字符集穷举"c???new???"明文，32位MD5包含字符串"3b605234ed"  
> bruteMD5 -a=c???new??? -bb=abcdefnutuvw_ -c=3b605234ed  

用数字、大写字母穷举明文"flag{?????}"(?代表未知5位)，16位MD5开头为"b6dff925"  
> bruteMD5 -a=flag{?????} -b=du -s=b6dff925 -i=16  

![help](https://github.com/playGitboy/bruteMD5/tree/master/img/bruteMD5_help.png)  

![test](https://github.com/playGitboy/bruteMD5/tree/master/img/bruteMD5_test.png)  

### 留坑
增加参数，允许不设置"-a"参数，即允许穷举任意字符来构造指定格式md5值
