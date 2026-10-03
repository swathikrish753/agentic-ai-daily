# 🤖 Agentic AI — Daily

> Building a real mental model of agentic AI — **one topic a day, from scratch, in Go.**
> Explained like I'm five, then taken deep enough to crack the interview.

![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-local%20LLMs-black)
![Learning in public](https://img.shields.io/badge/learning-in%20public-brightgreen)
![One day](https://img.shields.io/badge/cadence-1%20topic%2Fday-blue)

---

## 📖 About

Most agentic-AI material is Python-first and jumps straight to frameworks that
hide what's actually happening. This repo does the opposite: **hand-rolled, in Go,
nothing hidden.** Every concept is built up from zero — the agent loop, tool use,
memory, planning, RAG, guardrails — so the "magic" turns into plain engineering.

Each day = one small, runnable program + the intuition behind it.

## 🛠 Stack

- **Go** — the whole thing, by hand (no agent framework)
- **Ollama** — open-source LLMs running locally (no API keys, no cloud)
- **MCP** *(coming)* — the Model Context Protocol for wiring up tools

## 🧭 Philosophy

- **ELI5 first.** If I can't explain it to a five-year-old, I don't understand it yet.
- **Run everything.** No concept counts until the code executes on my machine.
- **Slow and steady.** One topic a day, no moving on until it clicks.
- **Interview-ready.** Every day ends with how I'd answer it out loud.

## 📅 Progress

| Day | Topic | The idea in one line |
|-----|-------|----------------------|
| 01  | What is an AI agent, really? | brain + tools + a loop → **think → do → see** |
| 02  | The feedback loop | the agent **reacts** to what it sees — it writes its own plan |
| 03  | Talking to a real LLM | an LLM is just **text in → text out**; the model talks, the code acts |
| 04  | Tool use | the LLM **picks** a tool; our code runs it — *model proposes, code disposes* |

*…updated daily.*

## 🗺 Roadmap

- [x] Tool use — let the LLM *pick* which tool to run
- [ ] Structured output — reliable JSON instead of fuzzy text
- [ ] Memory — short-term (context) vs. long-term
- [ ] Planning — breaking a big goal into steps
- [ ] RAG — giving the agent knowledge it wasn't trained on
- [ ] Guardrails & approval gates — safe autonomous action
- [ ] MCP, multi-agent, evaluation, observability — the real-world layer

## ▶️ Run it locally

```bash
# 1. start a local LLM brain (leave running in its own tab)
ollama serve
ollama pull llama3.2:3b

# 2. run any day
go run ./day04
```

---

*A daily build-in-public log on the road to mastering agentic systems.* 🚀
