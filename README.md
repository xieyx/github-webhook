# gitlab-webhook
[![](https://img.shields.io/github/release/xieyx/gitlab-webhook?style=flat-square)](https://github.com/xieyx/gitlab-webhook/releases)
[![](https://img.shields.io/github/license/xieyx/gitlab-webhook?style=flat-square)](https://github.com/xieyx/gitlab-webhook/blob/master/LICENSE)
![](https://img.shields.io/github/repo-size/xieyx/gitlab-webhook?style=flat-square)


## 📡 Overview
- The Gitlab-webhook is a webhook tool on gitlab,
- That can trigger bash scripts after monitoring git's push behavior
- The a line command handles the automatic build
- Built-in queue for tasks, quick response to Gitlab Webhook, 100% response 200 guaranteed

1. gitlab-webhook 是gitlab webhook自动构建工具.能监听git push行为,自动触发脚本.
1. 一条命令搞定webhook自动构建,无需复杂的配置.
1. 内置队列执行任务，迅速响应 gitlab webhook, 保证100% response 200

## 📜 Usage
### 1. Download && Install
- [releases](https://github.com/xieyx/gitlab-webhook/releases)
```shell script
cd ~
wget https://github.com/xieyx/gitlab-webhook/releases/download/v1.6.0/gitlab-webhook1.6.0.linux-amd64.tar.gz
tar -zxvf gitlab-webhook1.6.0.linux-amd64.tar.gz
cp ~/gitlab-webhook /usr/local/sbin
chmod u+x /usr/local/sbin/gitlab-webhook
```

run script
```
/usr/local/sbin/gitlab-webhook --bash /home/sh/test.sh
```

## 3. Command
- Daemonize run:  `nohup gitlab-webhook --bash /home/my.sh --secret mysecret -q &`
- Monitor run: `gitlab-webhook --bash /home/my.sh --secret mysecret`
- Quiet mode run: `gitlab-webhook --bash /home/my.sh --secret mysecret --quiet`
- Custom port mode run: `gitlab-webhook --bash /home/my.sh --secret mysecret --port 6100 --quiet`
- Hidden secret mode run: `gitlab-webhook --bash /home/my.sh  --quiet`

add systemd service
> /home/sh/hugo2www.sh is your script bash file
```shell script
cat > /lib/systemd/system/webhook << EOF
[Unit]
Description=gitlab-webhook
Documentation=https://github.com/xieyx/gitlab-webhook
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/gitlab-webhook --bash /home/sh/hugo2www.sh --secret qweqwe
Restart=on-failure
RestartSec=42s

[Install]
WantedBy=multi-user.target
EOF
```
```shell script
systemctl daemon-reload
systemctl start webhook
systemctl status webhook
```


## 4. WebHook
- Default port: 2020
- Http path: /web-hook
- Test URL: `http://ip:2020/ping`
- WebHook URL: `http://ip:2020/web-hook`


## 💌 Features
- Just run the binaries file
- Custom your bash script
- Custom your enter secret
- Custom your port. 0 ~ 65535
- Quiet operation

中文
- 直接运行二进制文件
- 自定义脚本路径
- 自定义密码
- 自定义端口. 0 ~ 65535
- 安静模式

```text
GLOBAL OPTIONS:
   --bash value, -b value    Execute the script path. eg: /home/hook.sh
   --port value, -p value    http port (default: 2020)
   --secret value, -s value  gitlab hook secret
   --quiet, -q               quiet operation (default: false)
   --verbose, --vv           print verbose (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```
中文
```text
GLOBAL OPTIONS:
   --bash value, -b value    Execute the script path. eg: /home/hook.sh 自定义脚本
   --port value, -p value    http port (default: 2020) 自定义端口,默认6666
   --secret value, -s value  gitlab hook secret 自定义密码, 不允许为空
   --verbose, --vv           print verbose (default: false) 打印更多详细信息
   --quiet, -q               quiet operation (default: false) 安静模式,默认关闭. -q 开启,不输出任何信息
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)

```
# How it works

![gitlab-webhook](https://raw.githubusercontent.com/xieyx/images/main/2022/03/upgit_20220322_1647929918.png)


- step 1:: Run your gitlab-webhook server

  - notice: port default 2020, http-path: /web-hook
  - 注意: 端口默认为 2020, 可以更改, http的路由: /web-hook
  - 查看自己的外网Ip: `curp ip.sb`

- step 2: Add webhook
  - 添加 webhook 参数

    ![配置第一步](https://raw.githubusercontent.com/xieyx/images/main/2022/03/upgit_20220322_1647915877.png)

    ![配置第二步](https://raw.githubusercontent.com/xieyx/images/main/2022/03/upgit_20220322_1647916039.png)

- step 3: run shell script
 - notice: Make sure that the last line write: exit 0
 - shell脚本的最后一行一定要写上 `exit 0` 代码
```
#!/bin/bash
echo "hello webhook"
exit 0
```
## 👋 Thanks
- See [GitbookIO](https://github.com/GitbookIO/go-github-webhook)
