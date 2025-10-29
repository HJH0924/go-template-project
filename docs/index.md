---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Go Template Project"
  text: "现代化的 Go 项目模板"
  tagline: 展示 Go 语言开发的最佳实践
  actions:
    - theme: brand
      text: 快速开始
      link: /development
    - theme: alt
      text: GitHub
      link: https://github.com/HJH0924/go-template-project

features:
  - icon: 🚀
    title: gRPC/Connect
    details: 使用 Connect RPC 框架，支持 HTTP/1.1、HTTP/2 和 gRPC，提供现代化的 API 服务
  - icon: 📝
    title: 结构化日志
    details: 使用 Go 1.21+ 的 log/slog 包，提供高性能的结构化日志记录
  - icon: ⚙️
    title: 配置管理
    details: 基于 Viper 的配置管理，支持 YAML 配置文件和环境变量
  - icon: 🎯
    title: 命令行工具
    details: 使用 Cobra 构建强大的 CLI 工具，支持子命令和参数解析
  - icon: 🏗️
    title: 分层架构
    details: Handler、Service 层清晰分离，易于维护、测试和扩展
  - icon: 🧪
    title: 单元测试
    details: 完整的单元测试示例和覆盖率报告，确保代码质量
  - icon: 🔍
    title: 代码检查
    details: 集成 golangci-lint 和 pre-commit hook，自动化代码质量保证
  - icon: 📦
    title: Protocol Buffers
    details: 使用 Buf 管理 proto 文件和代码生成，标准化 API 定义
---

