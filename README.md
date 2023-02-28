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


### 部署

#### 生成可执行文件
```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build . 
```

#### 创建目录并上载可执行文件，假设放在 `/var/www/go-qrcode`

#### nohup 执行
```
nohup ./go-qrcode &
```

（可选）更改端口
```
nohup ./go-qrcode -port 8087 &
```

#### Apache 反代

在 Apache 的配置文件中启用 mod_proxy 模块
```
LoadModule proxy_module modules/mod_proxy.so
LoadModule proxy_http_module modules/mod_proxy_http.so
```

在 /etc/apache2/sites-available 中创建 go-qrcode.conf
```
<VirtualHost *:80>
    ServerAdmin webmaster@localhost
    DocumentRoot /var/www/go-qrcode

    ServerName qr.example.com
    ProxyRequests Off
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:8085
    ProxyPassReverse / http://127.0.0.1:8085

    <Directory /var/www/go-qrcode>
        Options Indexes FollowSymLinks MultiViews
        AllowOverride FileInfo Options
        Order allow,deny
        allow from all
        Require all granted
    </Directory>

</VirtualHost>
```

启用新站点
```
sudo a2ensite go-qrcode.conf
```

重启 Apache
```
sudo systemctl reload apache2
```

完毕！