# go-qrcode-test

### 运行

默认端口 8085
```
go run .
```

更改端口
```
go run . -port=8889
```


### 访问

【 GET / POST 】http://localhost:8085


|  参数   |  类型  |  必填  |  描述  |
|  ----  |  ----  |  ----  |  ----  |
| content |  string |  是   |   内容  |
| fg     |  string |   否   |   前景颜色，Hex，默认 #000000， 传递时不需要#号|
| bg     |  string |   否   |   背景颜色，Hex，默认 #ffffff， 传递时不需要#号|
| block  |  integer |   否   |   二维码块的大小，默认 10 |
| border |  integer |   否   |   边框的大小，默认 20 |
| logo   |   file   |   否  |  logo文件，multipart/form-data |