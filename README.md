# SecretServer
服务端加密的文件服务器


## 使用方法

使用`curl`上传文件
```bash
curl -T file 'http://your.addr/filename?api_key=ru8w9qr&file_key=KKEEYY'
```

下载文件
```bash
curl 'http://your.addr/filename?file_key=KKEEYY'
```

生成 API_KEY：https://emn178.github.io/online-tools/sha3_256.html
