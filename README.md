fileserver
===

简易文件服务器, 将文件存储到Telegram Bot上, 方便对外进行分享?

目前支持最大单文件4G, 同时也支持Range方式下载(断点续传能力)

## 使用方式

### 服务端配置

这里是一份完整的配置模板

```json
{
    "log_info": {
        "level": "debug",  
        "console": false,
        "file": "/tmp/fileserver.log",
        "file_size": 20971520,
        "file_count": 5,
        "keep_days": 30
    },
    "file_db_info": {
        "host": "db",
        "port": 3306,
        "user": "abc",
        "pwd": "123456",
        "db": "filedb"
    },
    "server_info": {
        "address": ":9901"
    },
    "upload_fs": "tgbot",
    "fs_info": {
        "tgbot": [
            {
                "chatid": 111111,
                "token": "666666:111111-BiYhYznWWdaC_pprQ-of4z"
            }
        ]
    },
    "auth_info": {
        "test": "761f6035-5f16-474c-9a1c-82567x5c4a7b"
    },
    "io_info": {
        "max_upload_thread": 50,
        "max_download_thread": 100
    },
    "idgen_info": {
        "worker_id": 2
    },
    "fake_s3_info":{
        "enable": true,
        "bucket_list":["hackmd"]
    }
}
```

| 参数         | 备注                                                                                                       |
| ------------ | ---------------------------------------------------------------------------------------------------------- |
| log_info     | 日志项配置                                                                                                 |
| file_db_info | db配置                                                                                                     |
| server_info  | 服务器监听地址配置                                                                                         |
| upload_fs    | 上传使用的文件系统, 固定为tgbot                                                                            |
| fs_info      | 文件系统配置信息, 子项只能选用tgbot, 列表中可以填写多个机器人配置, 上传的时候会以轮询的方式进行上传        |
| auth_info    | 鉴权信息, kv结构, key为user(**ak**), value为password(**sk**), 支持多用户                                   |
| io_info      | 限制上传下载的连接数用的                                                                                   |
| idgen_info   | 目前只用于单机部署, 所以这里的值随便填一个10以内的值即可                                                   |
| fake_s3_info | 支持以s3的方式进行上传下载(只支持GetObject, PutObject 2种), 这个能力主要给hackmd用的, 正常情况下应该用不到 |

### 服务器部署

目前只支持docker部署, 拉取镜像`xxxsen/file_server:latest` 运行即可。

下面是一份测试配置, 配置文件的编写可以参考前面的内容.

```yml
version: "3.0"
services:
  fileserver:
    image: xxxsen/file_server:v0.0.1
    restart: always
    volumes:
      - "./config:/config"
    ports:
      - 127.0.0.1:9901:9901
    command: -config=/config/config.json
    depends_on:
      - db
  db:
    image: mariadb:10.4
    command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW
    restart: always
    volumes:
      - ./data/db/data:/var/lib/mysql
      - ./data/db/init:/docker-entrypoint-initdb.d
    environment:
      - MYSQL_ROOT_PASSWORD=${your password}
      - MYSQL_PASSWORD=${your password}
      - MYSQL_DATABASE=filedb
      - MYSQL_USER=${your username}
```

**程序运行需要依赖DB, 按这上面的配置, 需要将建表sql放到`./data/db/init` 中(将代码目录中的sql下面`.sql`结尾的文件放到init目录即可), 程序首次启动会创建对应的DB表**

### 本地命令行

#### shell 脚本

下面的配置是在自己的机器上运行的, 不是服务器!!!

linux下测试是ok的, windows应该是不行的(WSL是ok的), 如需要在非linux下使用, 可以使用下面的`golang客户端`

进入到`scripts/fsrz` 目录, 执行`sudo ./fsrz install` 进行安装即可, 安装脚本会创建`/etc/fsrz/`目录, 并生成模板文件`config.tplt`, 这里需要复制这个模板文件, 重命名为`config`(最终配置路径: `/etc/fsrz/config`)并修改其中的配置, 下面是参考模板。

```shell
AK=abc     # 用户名
SK=123456   # 密码
SCHEMA=https  # 上传/下载链接的schema
DOMAIN=file.mydomain.cc  # 默认的上传下载域名, 如果DOWNLOAD_DOMAIN/UPLOAD_DOMAIN 不填写, 那么则使用这个域名
DOWNLOAD_DOMAIN=fs.mydomain.cc # 指定下载域名
UPLOAD_DOMAIN=      # 指定上传域名
```

按上面配置完成后, 就可以使用`fsrz ${file_location}`的方式进行文件上传了, 上传完成后, 会拿到下面的结果, 其中的链接就是可下载链接。

```text
abc@machine:/fileserver/scripts/fsrz$ fsrz /etc/fstab 
read ak:abc
read sk:123456
read schema:https
read upload domain:file.mydomain.cc
read download domain:fs.mydomain.cc
read downkey from server:000134a51f3ba085
========
https://fs.mydomain.cc/file?down_key=000134a51f3ba085
```

#### golang客户端

通过镜像`xxxsen/fsrz` 获取, 使用方式参考`docker run -it --rm xxxsen/fsrz --help`