---
layout: post
title:  "打造基于WSL的开发环境"
date:   2018-10-29 20:13:00 +0800
categories: DevOps
---
现有工作电脑是只有8G内存的Dell笔记本，操作系统是window10家庭版，由于Linux上钉钉、Office等的支持不太好，又没有Macbook可用，笔者一直以来都是在虚拟机上用Linux来工作，但是8G内存对于后端开发实在不太够用，写代码的时候总是缺乏一种爽快感。

最近发现WSL（windows subsystem for Linux）可以在win10家庭版中启用，经过一番折腾，顺利得把开发调试环境从虚拟机迁移到了WSL中，内存降了不少。这里把笔者的折腾经历记录下来，希望能令有相同需求的小伙伴少走点弯路。

笔者的WSL目标有：

- [x] 支持python、golang、javascript前后端开发
- [x] 良好的terminal支持，包括powerline字体、tmux等
- [x] neovim的良好支持，包括YouCompleteMe、cpsm、vim-go等插件的支持
- [x] 能使用oh-my-zsh
- [x] 支持全屏
- [x] 可以使用docker（这里坑比较多）

{:toc}

## 1. 在Windows上启用WSL

1. 按windows键，输入控制面板，打开控制面板，点击程序 -> 启用或关闭Windows功能。![控制面板]({{ site.url }}/assets/wsl/control_panel.png)
1. 在Windows功能窗口，打开“适用于Linux的Windows子系统”，重启系统生效。![windows功能]({{ site.url }}/assets/wsl/windows_feature.png)

## 2. 在应用商店安装Ubuntu 18.04

1. 打开windows应用商店，搜索“Ubuntu 18.04”，安装。![应用商店]({{ site.url }}/assets/wsl/windows_store.png)
1. 安装完毕后打开Ubuntu，确保能正确运行，如果出现运行问题，检查WSL是否开启，以及是否忘了重启
1. 修改软件源为163源
```sh
sudo echo 'deb http://mirrors.163.com/ubuntu/ bionic main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ bionic-security main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ bionic-updates main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb http://mirrors.163.com/ubuntu/ bionic-backports main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ bionic main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ bionic-security main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ bionic-updates main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src http://mirrors.163.com/ubuntu/ bionic-backports main restricted universe multiverse
' > /etc/apt/sources.list
sudo apt update
```

## 3. 使用wsl-terminal代替默认的终端

默认的终端效果很差，我们希望能够有一个类似Linux或Mac上终端效果的替代品，这里我们考虑[wsl-terminal](https://github.com/goreliu/wsl-terminal)。

1. 下载wsl-terminal的最新release，我下载的是：https://github.com/goreliu/wsl-terminal/releases/download/v0.8.11/wsl-terminal-0.8.11.7z
1. 解压缩wsl-terminal到一个ntfs分区的目录下
1. 点击open-wsl.exe，即可打开wsl terminal。![wsl terminal]({{ site.url }}/assets/wsl/wsl_terminal.png)
1. 其它设置见wsl-terminal的github主页

## 4. 使用zsh和oh-my-zsh

1. zsh相较于bash在定制化等方面都有不小的提升，替换之，`sudo apt install zsh`
1. 安装[oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh)
```sh
bash -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
```
1. 修改wsl-terminal的启动方式，在wsl-terminal目录中找到快捷方式，右键 -> 属性，将目标中的bash修改为zsh
1. 将wsl-terminal的快捷方式放到任务栏，方便后续打开

## 5. 修改主题、颜色和使用powerline字体

1. wsl-terminal支持主题设置，在窗口标题栏右键 -> 选项，即可设置，我设置的是“base16-seti-ui.minttyrc”
1. wsl-terminal默认的配色目录的背景色一片绿，实在太伤眼睛了，修改之：
```sh
dircolors  # 查看当前的颜色配置
# 修改其中的 tw和ow，放到~/.zshrc里面，以下为我的设置
LS_COLORS='rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=01;34:ow=01;34:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.Z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:';
export LS_COLORS
```
1. 为了在终端支持一些图标字符，我们通常需要安装[powerline字体](https://github.com/powerline/fonts)，这里我安装的是DejaVuSansMono，把对应文件夹下面的ttf文件下载下来安装，然后在wsl-terminal中修改字体就可以了。 ![windows功能]({{ site.url }}/assets/wsl/powerline_font.png)

## 6. 安装、启动docker

网上搜索wsl上安装docker，大部分都是让你用windows版本的docker，但笔者的win10家庭版不支持hyper-v啊，一翻折腾后终于找到一个可行的办法，步骤如下：

1. 在/etc/apt/sources.list里面加入xenial的docker源，注意：这里不能用bionic的源，因为我们要安装旧版本的docker-ce。
```sh
sudo echo 'deb [arch=amd64] https://download.docker.com/linux/ubuntu xenial stable' >> /etc/apt/sources.list
```
1. 安装17.09.0版本的docker-ce，小版本也不能错，不然安装好的docker跑不起来
```sh
sudo apt install docker-ce=17.09.0~ce-0~ubuntu  # 这里的版本不能超过17.09.0
```
1. 用administrator启动wsl-terminal，执行以下命令启动docker服务
```sh
sudo cgroupfs-mount  # 在每次启动docker服务前必须执行cgroupfs-mount
sudo service docker start
sudo service docker status  # 一次不成功就多试几次
```
1. 运行`docker run --rm hello-world`检查docker是否安装成功，如果存在问题，检查上述步骤是否有错误
1. 最后的效果如下：
![docker]({{ site.url }}/assets/wsl/docker.png)

## 7. 关闭对Linux子系统的病毒实时扫描

在linux子系统中安装软件或者git clone经常会发现系统CPU比较高，查看任务管理器，发现病毒扫描服务占用了大量的CPU，关闭它可以顺畅很多。

1. 在用户目录的AppData\Local\Packages下找到对应的Ubuntu的目录，我的是`C:\Users\feiyu\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu18.04onWindows_79rhkp1fndgsc`
1. 打开设置 -> 更新和安全 -> Windows安全 -> 病毒和威胁防护，点击病毒和威胁防护设置，点击添加或删除排除项，将上面的目录加入到排除项中。
![antirus]({{ site.url }}/assets/wsl/antirus.png)

## 8. 安装tmux、vim等开发工具

1. 安装tmux和neovim：`sudo apt install tmux neovim`
1. 相关配置可以参考[笔者的配置](https://github.com/feiyuw/vim.d/)
1. vim如果需要支持YouCompleteMe和CPSM插件，安装相关library：`sudo apt install libboost-all-dev python3-dev`

## 9. 总结

Good:

1. 内存消耗比虚拟机降低了很多很多，之前我的8G内存根本不够用，现在基本在4.8G左右
1. 可以很方便地操作windows上的文件
1. 可以使用windows上的服务，如代理服务等
1. 启动的服务可以直接用windows上的浏览器访问和调试

Bad：

1. 磁盘IO性能一般
1. 部分命令得不到预期结果，如netstat

总的来说，wsl现在已经比较完善了，一般的开发已经可以胜任。

## 10. 参考

1. [https://huyinjie.xyz/2018/02/17/%E4%BD%BF%E7%94%A8wsl-terminal%E7%BE%8E%E5%8C%96WSL/](https://huyinjie.xyz/2018/02/17/%E4%BD%BF%E7%94%A8wsl-terminal%E7%BE%8E%E5%8C%96WSL/)
1. [https://www.reddit.com/r/bashonubuntuonwindows/comments/8cvr27/docker_is_running_natively_on_wsl/](https://www.reddit.com/r/bashonubuntuonwindows/comments/8cvr27/docker_is_running_natively_on_wsl/)
1. [https://blog.royink.li/post/ls_color_in_wsl](https://blog.royink.li/post/ls_color_in_wsl)
1. [https://medium.com/@leandrw/speeding-up-wsl-i-o-up-than-5x-fast-saving-a-lot-of-battery-life-cpu-usage-c3537dd03c74](https://medium.com/@leandrw/speeding-up-wsl-i-o-up-than-5x-fast-saving-a-lot-of-battery-life-cpu-usage-c3537dd03c74)
