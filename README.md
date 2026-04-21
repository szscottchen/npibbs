
##Introduction

This project is a discussion forum developed based on bbs go (an open-source community forum system developed in Go language), which can be used for internal discussions and sharing of topics within organizations such as research and development. For detailed information about bbs go, please refer to:
- Github: [ https://github.com/mlogclub/bbs-go ]( https://github.com/mlogclub/bbs-go )
- Gitee: [ https://gitee.com/mlogclub/bbs-go ]( https://gitee.com/mlogclub/bbs-go )

The source code storage address for this project is:
https://github.com/szscottchen/npibbs
This project follows the GPL3.0 protocol of the original bbs go and is open sourced.

This project is based on BBS GO and focuses on enhancing discussions and exchanges within the R&D department of an enterprise. Compared to BBS Go, the added content mainly includes:

-Attachments such as images and articles are deployed on the machine where the application server is located, which better takes into account the information security considerations of the internal R&D department.
-Add the function of establishing users through batch import in the backend management, and design the import template based on the organizational form of general enterprises.
-Added a new type of posting called 'Call for Support', which means calling for help. In this way, the person who posts a topic can call for more people to discuss a key topic or idea, and the topic owner can add points to the comments they think are helpful and valuable, increasing the points of relevant reply users.
-Added the function of accessing and following up on topics on the enterprise WeChat platform, so that users can use the forum on the enterprise WeChat platform.
-AI summarization function has been added to each topic, which can summarize the topic, extract the highlights of the discussion, and help identify viewpoints and ideas.


! [BBS GO Function Introduction] (docs/images/features. jpg)

##System architecture and technology stack

This project is based on * * bbs go * * and adopts a typical front-end and back-end separation architecture.

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

##Function Introduction

-Topic Management - Posting, Responding, Node Classification
-Article Management - Article Publishing
-* * Comment System * * - Comment Function
-* * AI Summary * * - Topic Intelligence Summary
-* * Enterprise WeChat * * - Enterprise WeChat Integration
-* * Search Function * * - Full Text Search
-* * User Center * * - Personal Settings, User List, Batch Import
-Topic Management - Topic Review, Node Management
-Article Management - Article Review and Tag Management
-* * System Settings * * - Menu, Role, API Management

##Installation and deployment

##Contact Information

## Contributors