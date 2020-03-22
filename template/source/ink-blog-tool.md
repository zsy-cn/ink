title: "静态博客构建工具 - Ink"
date: 2020-03-22 00:00:00 +0800
update: 2020-03-22 00:00:00 +0800
author: xavier
cover: "-/images/example.png"
tags:
    - 设计
    - 写作
preview: Ink 是一个 Golang 语言编写的开源静态博客构建工具，可以快速搭建博客网站。它无依赖跨平台，配置简单构建快速，注重简洁易用与更优雅的排版，设计主要参考 [ InkPaper ](https://github.com/InkProject/ink)。
---

## Ink 简介

Ink 是一个 Golang 语言编写的开源静态博客构建工具，可以快速搭建博客网站。它无依赖跨平台，配置简单构建快速，注重简洁易用与更优雅的排版，设计主要参考 [ InkPaper ](https://github.com/InkProject/ink) 

![纸小墨 - 简洁的静态博客构建工具](-/images/example.png)

### 开始上手

- 下载并解压 [Ink](http://github.com/zsy-cn/ink/)，运行命令 `ink preview`

  > 注意：Linux/macOS下，使用 `./ink preview`

- 使用浏览器访问 `http://localhost:9292` 预览。

### 特性介绍
- YAML格式的配置
- Markdown格式的文章
- 无依赖跨平台
- 超快的构建速度
- 不断改善的主题与排版
- 多文章作者支持
- 归档与标签自动生成
- 保存时实时预览页面
- 离线的全文关键字搜索

### 配置网站
编辑`config.yml`，使用如下格式：

``` yaml
site:
    title: 网站标题
    subtitle: 网站子标题
    limit: 每页可显示的文章数目
    theme: 网站主题目录
    comment: 评论插件变量(默认为Disqus账户名)
    root: 网站根路径 #可选
    lang: 网站语言 #支持I18n，可在theme/config.yml配置
    url: 网站链接 #用于RSS生成
    link: 文章链接形式 #默认为{title}.html，支持{year},{month},{day},{title}变量

authors:
    作者ID:
        name: 作者名称
        intro: 作者简介
        avatar: 作者头像路径

build:
    output: 构建输出目录 #可选, 默认为 "public"
    port: 预览端口
    copy:
        - 构建时将会复制的目录/文件
    publish: |
        ink publish 命令将会执行的脚本
```

### 创建文章
在`source`目录中建立任意`.md`文件（可置于子文件夹），使用如下格式：

``` yaml
title: 文章标题
date: 年-月-日 时:分:秒 #创建时间，可加时区如" +0800"
update: 年-月-日 时:分:秒 #更新时间，可选，可加时区如" +0800"
author: 作者ID
cover: 题图链接 #可选
draft: false #草稿，可选
top: false #置顶文章，可选
preview: 文章预览，也可在正文中使用<!--more-->分割 #可选
tags: #可选
    - 标签1
    - 标签2
type: post #指定类型为文章(post)或页面(page)，可选
hide: false #隐藏文章，只可通过链接访问，可选

---

Markdown格式的正文
```

### 发布博客
- 在博客目录下运行`ink publish`命令自动构建博客并发布。
- 或运行`ink build`命令将生成的`public`目录下的内容手动部署。

> Tips: 在使用`ink preview`命令时，编辑保存文件后，博客会自动重新构建并刷新浏览器页面。

## 定制支持

### 修改主题

默认主题在`theme`目录下，修改源代码后在该目录下运行`npm install`与`npm run build`重新构建。

页面包含`page.html`（文章列表）及`article.html`（文章）等，所有页面均支持[GO 语言 HTML 模板](http://golang.org/pkg/html/template/)语法，可引用变量。

### 添加页面

在`source`目录下创建的任意`.html`文件将被复制，这些文件中可引用`config.yml`中site字段下的所有变量。

### 源码编译

本地运行

1. 配置[ GO ](http://golang.org/doc/install)语言环境。
2. 运行命令`go get github.com/zsy-cn/ink`，编译并获取ink。
3. 运行命令`ink preview $GOPATH/src/github.com/zsy-cn/ink/template`，预览博客。

## 主题

- Theme: [https://github.com/zsy-cn/ink/template/theme](https://github.com/zsy-cn/ink/template/theme)

## 反馈贡献

如有问题可报告至 [https://github.com/zsy-cn/ink/issues](https://github.com/zsy-cn/ink/issues)。

## 更新历史

- [2020-03-22] 更改主题，删除部分第三方依赖。

## 更新计划

- 更换 Markdown 转 HTML 为 lute
- 命令行构建工具改为官方 flag
- 每篇文章可以自定义 bannder

## 正在使用

- [https://zsy-cn.github.io](https://zsy-cn.github.io)
