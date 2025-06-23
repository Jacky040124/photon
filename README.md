[English](README.md) | [中文](README.zh.md)

---



# Photon

<div align="center">
  <img src="public/logo.png" alt="Photon Logo" width="1000" height="400">
</div>

**Research. Instantly.**

The terminal research tool for people who value speed and simplicity. Ask anything, get answers in seconds, work from anywhere.

## Demo

<div align="center">
  <img src="public/demo.gif" alt="Photon Demo" width="800">
</div>

## What makes it different

**Fast** — Sub-second AI responses  
**Light** — 3MB binary, zero setup  
**Free** — Unlimited queries, always  
**Global** — Works everywhere, including China  

## For developers who think in terminal

```bash
ptn "quantum computing basics"
# → Instant, structured insights
# → No context switching
# → No distractions
```

**Simple tools. Powerful results.**

## Installation

### Homebrew (Recommended)
```bash
brew tap Jacky040124/photon
brew install photon
```

### Manual Download
Download pre-built binaries from [GitHub Releases](https://github.com/Jacky040124/photon/releases) and add to your PATH.

### Build from Source
```bash
git clone https://github.com/Jacky040124/photon.git
cd photon
go build -o ptn ./cmd
mv ptn /usr/local/bin/
```

## Setup

Set your OpenRouter API key:
```bash
export PHOTON_OPEN_ROUTER_KEY="your_openrouter_api_key_here"
```

> **Note:** Get a free OpenRouter API key from [openrouter.ai](https://openrouter.ai) - no credit card required for the free tier!

## Usage

Just run:

```
./ptn "your research topic here"
```

Examples:
```
ptn "machine learning basics"
ptn "rust vs go performance"
ptn "docker best practices"
```

## Output Format

**Photon** provides clean, structured output:
- **Summary**: Concise 2-3 sentence overview
- **Key Points**: 3-5 essential insights, clearly numbered

---

**Photon** — Simple tools. Powerful results. 
