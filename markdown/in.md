# WSL2(Arch Linux)使用systemd

## 概述

WSL本身是由Windows负责运行的，因此使用`tree`可以看到根进程不是systemd，而这将导致无法启动服务的守护进程(deamon)，如nginx、docker、mysql等

按道理上包括ip转发、原生systemd慢慢都会实现，但是现阶段只能采用第三方方案实现。

一个比较好的解决方案是使用[arkane-systems/genie](https://github.com/arkane-systems/genie)项目  
大概原理是在一个单独的命名空间跑systemd，从而实现能够正常执行systemctl指令并启动服务  

## 安装步骤

### 安装运行 .NET

首先需要安装`.NET Core runtime`，这个可以直接`yay dotnet`，找到带runtime的安装即可
安装完成后检查下是否其安装目录，如果不在默认位置需要手动写一下环境变量，如`export DOTNET_ROOT=/opt/dotnet`

### 安装相应的编译环境

需要python、python-markdown、python-six、python-packaging、python-setuptool、python-attrs等依赖包
而且按照测试来看似乎无法自动安装，因此最好先使用yay安装上他们（这些模块使用pip安装貌似也不行）

另外还需要安装 `make`

### 安装genie

使用 `yay -S genie-systemd` 下载 `genie` 源码并安装，如果中间出现缺少某个环境的话，退回到上一步安装即可

## 使用

`genie` 有三个指令：

`genie -i` 启动systemd进程
`genie -s` 启动systemd进程，并进入该环境终端
`genie -c <command>` 启动systemd进程，并执行相应的指令

比如要运行docker，使用`genie -s`进入到环境后，执行`sudo systemctl start docker`即可
而且在这里执行`pstree`，可以看到根进程已经变成了`systemd`
执行完成后退出这个命名空间后，`systemd`不会关闭