# Hooks 远程命令执行

# 使用方法
根据平台下载对应版本的程序(支持windows,linux)
赋予程序执行权限
程序会自动监听8888端口并执行Post过来的命令(数组方式执行组合命令)

## Post请求

# 请求地址
```
http://服务器IP:8888/hook
```

# 数据格式
```
{
	"command":["cd /root/test", "git pull"]
}
```
# 正确返回结果
```
{"error":0,"message":"Already up-to-date.\n"}
```

## GET请求
# 请求地址
```
http://服务器IP:8888/hook?经过base64后的命令1&经过base64后的命令2...
```
# 实例
如果需要执行命令
```
cd root
ls
```
则需要发送的格式为
```
http://服务器IP:8888/hook?Y2QgL3Jvb3Q=&bHM=
```
