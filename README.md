# dig

**dig** is a high-performance, developer-friendly terminal research tool that delivers concise, structured summaries on any topic you ask. Just type a query like:

```
dig "black holes"
```

and get an instant, beautifully formatted report—complete with source links for every key piece of information. Ideal for students, developers, and curious minds who want fast, reliable knowledge without leaving the terminal.

## Features
- ⚡ **High Performance:** Fast LLM-backed research with minimal wait.
- 🖥️ **Beautiful TUI:** Clean, modern terminal interface with spinner and clear output.
- 🔗 **Source Links:** Every summary includes links to credible sources.
- 🧑‍💻 **Developer Friendly:** Easy to install, configure, and use.
- 🚀 **Instant Results:** Get structured, actionable knowledge in seconds.

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/dig.git
   cd dig
   ```

2. **Install Go dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build the binary:**
   ```sh
   go build -o dig
   ```

4. **(Optional) Move to your PATH:**
   ```sh
   mv dig /usr/local/bin/
   ```

## Environment Setup

Create a `.env` file in the project root with your API key(s):

```
OPENAI_API_KEY=your_openai_api_key_here
```

> **Note:** You can get an OpenAI API key from [platform.openai.com](https://platform.openai.com/account/api-keys)

## Usage

Just run:

```
./dig "your research topic here"
```

Example:
```
./dig "quantum computing"
```

## Why dig?
- No more tab switching or information overload.
- Get the facts, key points, and sources—instantly.
- Perfect for quick research, code comments, or learning on the fly.

---

**dig** — Fast, reliable knowledge, right in your terminal. 