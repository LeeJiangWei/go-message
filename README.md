# go-message

使用微信测试号和企业微信应用搭建自己的消息推送服务。只要向自己的服务器发送一个简单的 GET 请求，即可将消息推送至你的手机微信上。

## 配置

在推送之前，需要先在系统上配置好微信相关的设置。登录系统后台，注册一个新用户，然后参考以下信息配置。

微信测试号和企业微信应用可以只配置其中一个，但必须填写所有字段，否则无法推送。

### 用户信息

+ 用户名：随意，同时用于在推送消息时 URL 中辨识用户。例如：`GET http://你的域名/template/用户名`
+ 密码：随意，仅用于登录后台系统。
+ 消息推送 Token：随意，在推送消息时携带，用于验证推送者身份。例如：`GET http://你的域名/template/用户名?token=hello`

### 微信测试号

1. 前往 [微信公众平台](https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index)
2. 在 **测试号信息** 一栏中获取 appID 和 appsecret。
3. 在 **测试号二维码** 一栏中扫描二维码关注自己的测试号，并在右侧列表中获取自己的微信号，作为接收者 ID。
4. 在 **模板消息接口** 一栏中新增模板，模板标题随意填写，模板内容填入：`From: {{from.DATA}} {{description.DATA}} {{remark.DATA}}`，提交后获取模板 ID。
5. 在 **接口配置信息** 一栏中，URL 填入自己服务器的域名 / IP 加上后缀，例如：`http://你的域名/verify/用户名`；填入任意 Token 作为接口配置 Token。填完先不要点击验证，先在回到推送系统中填好接口配置 Token 并成功修改后，再点击验证。

### 企业微信应用

1. 在 [企业微信](https://work.weixin.qq.com/) 注册一个企业微信号（不需要企业资质）。
2. 在 [企业信息](https://work.weixin.qq.com/wework_admin/frame#profile) 页面最下方获取企业 ID。
3. 在 [微信插件](https://work.weixin.qq.com/wework_admin/frame#profile/wxPlugin) 页面用微信扫码关注，这样消息才会直接推送到微信上。
4. 在 [应用管理](https://work.weixin.qq.com/wework_admin/frame#apps) 页面创建一个新应用，并获取应用的 AgentId 和 Secret。消息将通过这个应用推送。
5. 在 [通讯录](https://work.weixin.qq.com/wework_admin/frame#contacts) 页面右侧列表中，点击自己进入详情，获取自己的账号最为接收者 ID（一般是自己名字的拼音）。
6. 卡片消息默认 URL 可以随意填写，用于卡片消息点进去的页面，但是必须要有而且要保证是合法 URL，否则微信服务器会报错。可以使用系统的主页作为 URL。

## 推送

发送简单 GET 或 POST 请求到自己服务器上即可完成推送。

注意：**所有推送**必须携带 `token` 参数来验证身份，值为用户信息配置中的消息推送 Token。可以在 URL 中用 query params 的方式，也可以用 form data 的方式。

### 推送至微信测试号

微信测试号仅支持模板消息。

#### 模板消息

推送 URL（GET / POST）：`http://你的域名/template/用户名`

可选参数字段（用于显示在模板的不同位置）：

+ from
+ desc
+ remark

注：如果用 GET 请求，参数需要全部用 query params 的方式。如果使用 POST 请求，需要全部用 form data 的方式。

示例：`http://你的域名/template/用户名?token=token&from=webpage&desc=Hello%20Wrold&remark=这是一条测试消息`

效果：

<a href="https://sm.ms/image/GHiuts3IolUwjpn" target="_blank"><img src="https://s2.loli.net/2022/03/09/GHiuts3IolUwjpn.jpg" width="350"></a>

### 推送至企业微信应用

企业微信应用支持纯文本消息、文字卡片消息。

#### 纯文本消息

推送 URL（GET / POST）：`http://你的域名/plaintext/用户名`

必需参数字段：

+ content

注：如果用 GET 请求，参数需要全部用 query params 的方式。如果使用 POST 请求，需要全部用 form data 的方式。

示例：`http://你的域名/plaintext/用户名?token=token&content=Hello,%20traveller`

#### 文字卡片消息

推送 URL（GET / POST）：`http://你的域名/textcard/用户名`

必需参数字段：

+ title
+ desc
+ url（如果不填，则会使用配置时的默认 URL）

注：如果用 GET 请求，参数需要全部用 query params 的方式。如果使用 POST 请求，需要全部用 form data 的方式。

示例：`http://你的域名/textcard/用户名?token=token&title=Hello&desc=World`

效果：

<a href="https://sm.ms/image/rnMPmTevh3B2AxH" target="_blank"><img src="https://s2.loli.net/2022/03/09/rnMPmTevh3B2AxH.jpg" width="350"></a>

## 部署

直接运行打包好的二进制文件即可。初次运行会产生一个配置文件，里面可以配置 JWT 签发相关设置，是否启用 Redis 作为缓存（默认否），服务运行的端口（默认80）。重启服务器后生效。

## 编译

### Windows

```sh
GOPROXY=https://goproxy.cn,direct GOOS=windows GOARCH=amd64 go build -o go-message go-message-pusher
```

### Linux

```sh
GOPROXY=https://goproxy.cn,direct CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o go-message go-message-pusher
```
