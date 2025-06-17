# Photon

**Photon** is a lightning-fast, developer-friendly terminal research tool that delivers packets of pure knowledge at light speed. Powered by OpenRouter's free Mistral 3.1 model, just type a query like:

```
ptn "quantum computing"
```

and get an instant, beautifully formatted report with key insights and summaries. Perfect for developers, students, and curious minds who want fast, accurate knowledge without leaving the terminal.

## Features
- ⚡ **Light Speed:** Packets of pure knowledge delivered instantly via OpenRouter's free Mistral 3.1 small 24b model.
- 🖥️ **Beautiful TUI:** Clean, modern terminal interface with elegant spinner animations.
- 🎯 **Focused Output:** Structured summaries and key points without information overload.
- 🧑‍💻 **Developer Friendly:** Easy to install, configure, and use in any workflow.
- 🚀 **Instant Results:** Get actionable knowledge in seconds, right in your terminal.
- 💰 **Cost Effective:** Leverages OpenRouter's free tier for unlimited research.

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/photon.git
   cd photon
   ```

2. **Install Go dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build the binary:**
   ```sh
   go build -o ptn ./cmd
   ```

4. **(Optional) Move to your PATH:**
   ```sh
   mv ptn /usr/local/bin/
   ```

## Environment Setup

Create a `.env` file in the `configs/` directory with your OpenRouter API key:

```
PHOTON_OPEN_ROUTER_KEY=your_openrouter_api_key_here
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

## Why Photon?
- **Zero Context Switching**: Research without leaving your terminal workflow
- **Instant Knowledge**: Get facts and insights in seconds, not minutes
- **Developer Optimized**: Perfect for code comments, documentation, or quick learning
- **Free & Fast**: Powered by OpenRouter's free tier with Mistral 3.1 performance
- **Clean Output**: No ads, no clutter, just the information you need

---

**Photon** — Packets of pure knowledge at light speed, right in your terminal. Powered by OpenRouter. 