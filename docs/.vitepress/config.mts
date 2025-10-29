import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Go Template Project",
  description: "一个现代化的 Go 项目模板，展示 Go 语言开发的最佳实践",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '首页', link: '/' },
      { text: '开发指南', link: '/development' }
    ],

    sidebar: [
      {
        text: '文档',
        items: [
          { text: '开发指南', link: '/development' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/HJH0924/go-template-project' }
    ],

    footer: {
      message: '基于 MIT 许可发布',
      copyright: 'Copyright © 2025'
    }
  }
})
