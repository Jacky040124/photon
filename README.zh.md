[English](README.md) | [中文](README.zh.md)

---



# Photon

<div align="center">
  <img src="public/logo.png" alt="Photon Logo" width="800" height="200">
</div>

专为追求效率的技术人员打造的终端研究神器。随手一问，秒出答案，让知识触手可及。

## 演示

<div align="center">
  <img src="public/demo.gif" alt="Photon Demo" width="800">
</div>

## 核心优势

**闪电响应** — 毫秒级 AI 回答  
**极简轻量** — 仅 3MB，开箱即用  
**完全免费** — 无限查询，永不收费  
**全球可达** — 无需翻墙

## 专为命令行极客打造

```bash
ptn "quantum computing basics"
# → 瞬间获得结构化答案
# → 不打断工作流
# → 专注高效
```

**极简工具，极致体验。**

## 安装

### Homebrew（推荐）
```bash
brew tap Jacky040124/photon
brew install photon
```

### 直接下载
从 [GitHub Releases](https://github.com/Jacky040124/photon/releases) 下载对应系统的可执行文件，放入系统 PATH 即可。

### 源码编译
```bash
git clone https://github.com/Jacky040124/photon.git
cd photon
go build -o ptn ./cmd
mv ptn /usr/local/bin/
```

## 配置

配置 OpenRouter API 密钥：
```bash
export PHOTON_OPEN_ROUTER_KEY="your_openrouter_api_key_here"
```

> **提示：** 在 [openrouter.ai](https://openrouter.ai) 免费注册获取 API 密钥，无需绑卡即可使用！

## 使用指南

一行命令搞定：

```
./ptn "你想了解的任何问题"
```

使用示例：
```
ptn "机器学习入门"
ptn "Go 和 Rust 性能对比"
ptn "Docker 最佳实践"
ptn "区块链技术原理"
```

## 输出格式

**Photon** 为你精心整理信息：
- **核心要点**：一目了然的关键信息
- **深度解析**：3-5 个精华知识点，条理清晰

---

**Photon** — 极简工具，极致体验。