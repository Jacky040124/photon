[English](README.md) | [中文](README.zh.md)

# Photon

<p align="center">
  <img src="public/social-preview.png" alt="Photon 把终端中的一个问题整理为 AI 摘要与关键要点" width="100%">
</p>

<p align="center"><strong>不离开终端，把一个问题整理成简洁的 AI 摘要与关键要点。</strong></p>

<p align="center">
  <a href="https://github.com/Jacky040124/photon/releases/latest"><img src="https://img.shields.io/github/v/release/Jacky040124/photon?style=flat-square&color=22d3ee" alt="最新版本"></a>
  <a href="https://github.com/Jacky040124/photon/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/Jacky040124/photon/release.yml?style=flat-square&label=release" alt="发布流程"></a>
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.24">
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux-111827?style=flat-square" alt="macOS 与 Linux">
</p>

## 真实交付证明

| 已发布 | 可安装 | 有完整产品体验 |
| --- | --- | --- |
| 版本化 Release 同时发布 macOS 与 Linux 的 Intel 和 ARM 原生二进制文件。 | 独立维护的 Homebrew Formula 会把可执行文件安装为 `ptn`。 | CLI 包含模型目录、交互式选择、持久化配置、加载状态与超时降级。 |

[查看 Releases](https://github.com/Jacky040124/photon/releases) · [查看发布工作流](.github/workflows/release.yml) · [查看 Homebrew Formula](https://github.com/Jacky040124/homebrew-photon/blob/main/photon.rb)

## 我实现了什么

* 基于 Cobra 的清晰命令入口：用 `ptn "..."` 提问，用 `ptn model` 查看或切换当前模型。
* OpenRouter 客户端与结果解析器，把模型输出稳定整理为 `Summary` 和 `Key Points`。
* 基于 Bubble Tea 与 Lip Gloss 的终端体验，包含交互式模型选择器与结构化结果样式。
* 自动构建四个平台二进制文件的发布流程，以及 Homebrew 分发路径。

## 真实 CLI 演示

<p align="center">
  <img src="public/demo.gif" alt="真实 Photon 终端会话：回答问题并输出摘要与关键要点" width="800">
</p>

上方录屏来自真实 CLI。回复内容与等待时间取决于当前选择的 OpenRouter 模型。

## 安装与运行

### Homebrew

```bash
brew tap Jacky040124/photon
brew install photon
```

创建 [OpenRouter API Key](https://openrouter.ai/keys)，然后提供给 Photon：

```bash
export PHOTON_OPEN_ROUTER_KEY="your_openrouter_api_key"
ptn "Transformer 架构解决了什么问题？"
```

Photon 会把回答渲染为一段简短摘要和一组编号要点。

### 下载二进制文件

[最新 Release](https://github.com/Jacky040124/photon/releases/latest) 提供以下构建：

* Intel 与 Apple 芯片的 macOS
* amd64 与 arm64 的 Linux

下载对应文件后，将其重命名为 `ptn`，添加可执行权限，并放入系统 `PATH`。

### 从源码构建

当前 [`go.mod`](go.mod) 声明 Go 1.24.3。

```bash
git clone https://github.com/Jacky040124/photon.git
cd photon
go build -o ptn ./cmd
./ptn --help
```

## 模型控制

```bash
ptn model list
ptn model current
ptn model set kimi
ptn model info kimi
ptn model reset
```

当前模型保存在 `~/.photon/config.json`。推荐通过上方展示的 `PHOTON_OPEN_ROUTER_KEY` 环境变量提供凭证。

## 工作原理

```text
终端问题
    ↓
Cobra 命令与本地模型配置
    ↓
OpenRouter Chat Completion 请求
    ↓
Summary 与 Key Points 解析
    ↓
Bubble Tea 与 Lip Gloss 终端输出
```

实现刻意保持紧凑。命令处理位于 [`cmd/`](cmd)，API 客户端、模型目录、解析与终端渲染位于 [`pkg/`](pkg)。

## 当前边界

Photon 是轻量的 AI 回答客户端，不是网页搜索引擎。当前版本：

* 通过 OpenRouter 把一个文本问题发送给所选模型
* 把模型回复整理为摘要与关键要点
* 不会浏览网页、验证事实或附加来源引用
* 受 OpenRouter 的模型可用性、限制、延迟与定价影响

这个边界是有意保留的。Photon 优化的是快速的终端工作流，不声称模型回答可以替代基于来源的研究。

## 发布流程

匹配 `v*` 的 Tag 会触发 [GitHub Actions](.github/workflows/release.yml)，构建并把以下文件附加到 GitHub Release：

```text
ptn-darwin-amd64
ptn-darwin-arm64
ptn-linux-amd64
ptn-linux-arm64
```

---

由 [Jacky Zhong](https://github.com/Jacky040124) 构建。
