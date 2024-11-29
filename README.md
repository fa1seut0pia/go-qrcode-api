### 仓库地址

https://github.com/fa1seut0pia/go-qrcode-api

一个生成简单二维码的api, 使用`go`语言实现, 支持嵌套`logo`和水印, 支持`docker`部署

这里选用等宽字体 `SarasaFixedSC-Regular.ttf`, 方便计算字体宽度

### 构建和运行

```shell
# 直接运行
docker run -p 7688:7688 -d --name go-qrcode-api fa1seut0pia/go-qrcode-api

# 或者构建后运行
docker build -t go-qrcode-api .
docker run -p 7688:7688 -d --name go-qrcode-api go-qrcode-api

```

### 使用例子

- 常规 http://127.0.0.1:7688/qrcode?t=helloWorld
- 带logo http://127.0.0.1:7688/qrcode?t=https://github.com&l=https://avatars.githubusercontent.com/u/9919

![常规](https://raw.githubusercontent.com/fa1seut0pia/go-qrcode-api/refs/heads/master/data/normal.png) ![带logo](https://raw.githubusercontent.com/fa1seut0pia/go-qrcode-api/refs/heads/master/data/logo.png)

- 带水印 http://127.0.0.1:7688/qrcode?t=https://github.com&l=https://avatars.githubusercontent.com/u/9919&w=Github.com
- 水印 + logo + base64内容(可在text内容比较复杂时使用)  
  http://127.0.0.1:7688/qrcode?t=aHR0cHM6Ly9naXRodWIuY29tL2ZhbHNlVXQwcGlhL2dvLXFyY29kZS1hcGk=&l=https://avatars.githubusercontent.com/u/9919&w=github.com/fa1seUt0pia&e=1

![带水印](https://raw.githubusercontent.com/fa1seut0pia/go-qrcode-api/refs/heads/master/data/logo-watermark.png) ![带logo](https://raw.githubusercontent.com/fa1seut0pia/go-qrcode-api/refs/heads/master/data/base64-logo-watermark.png)

### 参数说明

| 参数 | 说明                             |
|----|--------------------------------|
| t  | text 内容, 支持base64 (需要带上标识参数e)  |
| l  | logo URL 地址, 可选                |
| w  | watermark 底部水印内容, 可选           |
| e  | 任意值, 表示 t 是 base64 编码之后的内容, 可选 |
