# Tabby AI Coding Assistant Analysis

## Overview
Tabby is an open-source, self-hosted AI coding assistant. It provides features like code completion, chat, and search.

## Features & APIs

### Code Completion (`/v1/completions`)
- **Request**: `CompletionRequest`
  - `language`: optional string
  - `segments`: prefix, suffix, filepath, git_url, declarations, relevant_snippets
  - `temperature`, `seed`, `debug_options`
- **Response**: `CompletionResponse`
  - `id`: string
  - `choices`: array of `{index, text}`

### Chat (`/v1/chat/completions`)
- Follows OpenAI-compatible chat completion API.

### Code Search (Internal/v1beta)
- `CodeSearchQuery`: `git_url`, `filepath`, `language`, `content`
- `CodeSearchHit`: `scores` (rrf, bm25, embedding), `doc` (body, filepath, git_url, language, start_line)

### Events (`/v1/events`)
- Logs events like `view`, `select`, `dismiss` for completions.

## Implementation Details in HyperHarness
- **Parity Tools**: Implement `tabby_completion` and `tabby_chat` in `tools/tabby_parity.go`.
- **AST Indexing**: Tabby uses tree-sitter for parsing and indexing. HyperHarness's `internal/repograph` already has basic indexing; it should be extended to support more languages and detailed symbol extraction (functions, classes, etc.) to match Tabby's retrieval quality.
- **Repository Context**: Tabby connects to git repositories and indexes them. HyperHarness should use its `repograph` to provide similar context to the LLM during completion.
