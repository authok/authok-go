![Authok 的 Go SDK](https://cdn.authok.com/website/sdks/banners/authok-go-banner.png)

<div align="center">

[![GoDoc](https://pkg.go.dev/badge/github.com/authok/authok-go.svg)](https://pkg.go.dev/github.com/authok/authok-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/authok/authok-go?style=flat-square)](https://goreportcard.com/report/github.com/authok/authok-go)
[![Release](https://img.shields.io/github/v/release/authok/authok-go?include_prereleases&style=flat-square)](https://github.com/authok/authok-go/releases)
[![License](https://img.shields.io/github/license/authok/authok-go.svg?style=flat-square)](https://github.com/authok/authok-go/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/authok/authok-go/main.yml?branch=main&style=flat-square)](https://github.com/authok/authok-go/actions?query=branch%3Amain)
[![Codecov](https://img.shields.io/codecov/c/github/authok/authok-go?style=flat-square)](https://codecov.io/gh/authok/authok-go)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fauthok%2Fauthok-go.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fauthok%2Fauthok-go?ref=badge_shield)

📚 [文档](#documentation) • 🚀 [入门指南](#getting-started) • 💬 [反馈](#feedback)

</div>



-------------------------------------

## 文档

- [Godoc](https://pkg.go.dev/github.com/authok/authok-go) - 请查看 Go SDK 文档以了解详情。
- [管理API 文档](https://authok.com/docs/api/management/v1) - 请浏览 Authok 管理 API，并了解与此 SDK 交互的方式。
- [文档站点](https://www.authok.com/docs) — 请浏览我们的文档站点并了解更多关于 Authok 的信息。
- [示例](./EXAMPLES.md) - 关于 SDK 使用的更多示例。

## 入门指南

### 要求

该库遵循与Go相同的[支持政策](https://go.dev/doc/devel/release#policy)。当前正在积极支持最近的两个重要版发布的Go版本，如果有兼容性问题将会进行修复。尽管旧版本的Go也许仍能运行，但我们将不会积极测试和修复这些旧版本与库的兼容性问题。

- Go 1.19+

### 安装步骤

```shell
go get github.com/authok/authok-go
```

### 使用说明

```go
package main

import (
	"log"

	"github.com/authok/authok-go"
	"github.com/authok/authok-go/management"
)

func main() {
	// 从您的 Authok 应用 Dashboard 获取这些信息。
  // 应用程序需要授权为Machine To Machine才能请求Authok Management API的访问令牌，
  // 具有所需的权限（作用域）。
	domain := "example.authok.com"
	clientID := "EXAMPLE_16L9d34h0qe4NVE6SaHxZEid"
	clientSecret := "EXAMPLE_XSQGmnt8JdXs23407hrK6XXXXXXX"

	// 使用域名、客户端ID和客户端秘钥初始化一个新的客户端。
  // 或者，您可以指定一个访问令牌：
  // `management.WithStaticToken("token")`
	authokAPI, err := management.New(
		domain,
		management.WithClientCredentials(clientID, clientSecret),
	)
	if err != nil {
		log.Fatalf("failed to initialize the authok management API client: %+v", err)
	}

// 现在我们可以与Authok Management API进行交互了。
// 例如：创建一个新的 client。
	client := &management.Client{
		Name:        authok.String("My Client"),
		Description: authok.String("Client created through the Go SDK"),
	}

  // 传入的 client 将会从响应中获取实例数据。
  // 这意味着，在此请求之后，我们可同一个 client 对象中访问 clientID。
	err = authokAPI.Client.Create(client)
	if err != nil {
		log.Fatalf("failed to create a new client: %+v", err)
	}

  // 使用getter函数来安全地访问字段，避免由空指针引起的宕机问题。
	log.Printf(
		"Created an authok client successfully. The ID is: %q",
		client.GetClientID(),
	)
}
```

### 频率限制

Authok 管理API对所有API客户端都实施速率限制。当达到限制时，SDK将在后台处理它，
当限制被解除时，SDK将自动重新尝试API请求。

> **注意事项**
> SDK 不能防止 `http.StatusTooManyRequests` 错误，它将根据 `X-Rate-Limit-Reset` 头部信息的剩余秒数进行限速控制。

## 反馈

### 贡献信息

我们欢迎您对该仓库进行反馈及贡献！在开始之前，请先阅读以下内容：

- [贡献指南](./CONTRIBUTING.md)
- [Authok 通用贡献准则](https://github.com/authok/open-source-template/blob/master/GENERAL-CONTRIBUTING.md)
- [Authok 行为准则](https://github.com/authok/open-source-template/blob/master/CODE-OF-CONDUCT.md)

### 发起问题

如需提供反馈或报告错误，请在[我们的问题跟踪器](https://github.com/authok/authok-go/issues)上发起一个问题。

### 漏洞报告

请不要在公共的 Github 问题跟踪器上报告安全漏洞。[责任披露计划](https://authok.com/responsible-disclosure-policy)详细说明了披露安全问题的过程。

---

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="https://cdn.authok.cn/website/sdks/logos/authok_light_mode.png" width="150">
    <source media="(prefers-color-scheme: dark)" srcset="https://cdn.authok.com/website/sdks/logos/authok_dark_mode.png" width="150">
    <img alt="Authok Logo" src="https://cdn.authok.com/website/sdks/logos/authok_light_mode.png" width="150">
  </picture>
</p>

<p align="center">Authok 是一款易于实现、适应性强的身份验证和授权平台。<br />了解更多请访问 <a href="https://authok.com/why-authok">为什么选择 Authok？</a></p>

<p align="center">该项目基于 MIT 许可证获得授权。详细信息请参阅 <a href="./LICENSE.md"> LICENSE</a> 文件。</p>
