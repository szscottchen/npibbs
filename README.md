
## Introduction

This project is a discussion forum developed based on bbs go (an open-source community forum system developed in Go language), which can be used for internal discussions and sharing of topics within organizations such as research and development. For detailed information about bbs go, please refer to:
- Github: [ https://github.com/mlogclub/bbs-go ]( https://github.com/mlogclub/bbs-go )
- Gitee: [ https://gitee.com/mlogclub/bbs-go ]( https://gitee.com/mlogclub/bbs-go )

The source code storage address for this project is:
https://github.com/szscottchen/npibbs
This project follows the GPL3.0 protocol of the original bbs go and is open sourced.

This project is based on BBS GO and focuses on enhancing discussions and exchanges within the R&D department of an enterprise. Compared to BBS Go, the added content mainly includes:

- Attachments such as images and articles are deployed on the machine where the application server is located, which better takes into account the information security considerations of the internal R&D department.
- Add the function of establishing users through batch import in the backend management, and design the import template based on the organizational form of general enterprises.
- Added a new type of posting called 'Call for Support', which means calling for help. In this way, the person who posts a topic can call for more people to discuss a key topic or idea, and the topic owner can add points to the comments they think are helpful and valuable, increasing the points of relevant reply users.
- Added the function of accessing and following up on topics on the enterprise WeChat platform, so that users can use the forum on the enterprise WeChat platform.
- AI summarization function has been added to each topic, which can summarize the topic, extract the highlights of the discussion, and help identify viewpoints and ideas.



## System architecture and technology stack

This project is based on bbs go and adopts a typical front-end and back-end separation architecture.

bbs-go/
∝ - Server/# Backend Services (Go)
∝ - Site/# Front end Site (Nuxt.exe+Vue 3)
∝ - admin/# Management Console (Vue 3+Arco Design)

### server

[![bbs-go-server]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-server.yml/badge.svg )]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-server.yml )

>Built on Golang, providing interface data support.

tech stack

- iris ([ https://github.com/kataras/iris ]( https://github.com/kataras/iris ））MVC framework for Go language
- gorm ([ http://gorm.io ]( http://gorm.io ））The best Go language database ORM framework to use
- resty ([ https://github.com/go-resty/resty ]( https://github.com/go-resty/resty ））A user-friendly HTTP client for Go language
- cron ([ https://github.com/robfig/cron ]( https://github.com/robfig/cron ））Timed Task Framework
- goquery ([ https://github.com/PuerkitoBio/goquery ]( https://github.com/PuerkitoBio/goquery ））HTML DOM element parsing

### site

[![bbs-go-site]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-site.yml/badge.svg )]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-site.yml )

>Front end page rendering service, built on nuxt. js.

tech stack

- vue.js ([ https://vuejs.org ]( https://vuejs.org ））Progressive JavaScript Framework
- nuxt.js ([ https://nuxtjs.org ]( https://nuxtjs.org ））Vue based server-side rendering framework with explosive efficiency

### admin

[![bbs-go-admin]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-admin.yml/badge.svg )]( https://github.com/mlogclub/bbs-go/actions/workflows/bbs-go-admin.yml )

>Management backend system, built based on ` vue.js+element ui `.

tech stack

- vue.js ([ https://vuejs.org ]( https://vuejs.org ））Progressive JavaScript Framework
- element-ui ([ https://element.eleme.cn ]( https://element.eleme.cn ））Ele.me's open-source frontend library based on Vue.js

## Function Introduction

- Topic Management - Posting, Responding, Node Classification
- Article Management - Article Publishing
- Comment System - Comment Function
- AI Summary - Topic Intelligence Summary
- Enterprise WeChat - Enterprise WeChat Integration
- Search Function - Full Text Search
- User Center - Personal Settings, User List, Batch Import
- Topic Management - Topic Review, Node Management
- Article Management - Article Review and Tag Management
- System Settings - Menu, Role, API Management

## Installation and deployment
### prepare
- Download the source code to <local source code directory>.
- Download Go, with a version of Go 1.24.6 or above. install
- Install pnpm
- Install MySQL database
- Create a <project operation root directory> on the server

The above installation process refers to the reference documentation of the relevant system.

### Installation:
- Configure several key files in <local source code directory>: server\bbs-go.yaml, site\env.production, and admin\.env.production. Refer to the explanations in the original files for configuration content.
- Compile the source code using the command `go build -o <正式使用的程序名> main.go`.
Execute `pnpm build` in the `site` directory of `<local source code directory>`
Execute `pnpm build` in the `admin` directory of `<local source code directory>`
- Create a new project database in Mydql: create database <project database>; grant all privileges on <project database>.* to '<database administrator>'@'localhost'
- Configure the application callback domain, trusted domain, and enterprise trusted IP in the WeChat Work backend. For specific details, refer to the WeChat Work documentation.
- Copy <the officially used program name> and bbs-go.yaml to <the project's root directory for operation
- Copy the two directories, locales and migrations, located in the <local source code directory>\server directory to the <project runtime root directory
- Copy all contents in the \dist directory under the <local source code directory>\site directory to <project running root directory>\site
- Copy all contents from the \dist directory under the \admin directory in the <local source code directory> to the \admin directory in the <project runtime root directory
Create two empty directories named \temp and \uploads under the <project root directory for operation

If nginx and Docker configurations are involved, please refer to the relevant system documentation.

### Start the service
For the first login, access the browser interface at http://<server address and port number>/install to complete the initialization of the database and administrator.

## Contact Information

### Mailbox

<cayuanyuan@aliyun.com>

### QQ
![开发者QQ](docs/images/qq.png)

### WeChat

![开发者微信](docs/images/wechat.png)

## Contributors