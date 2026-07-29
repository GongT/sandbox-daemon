---
description: "Use when writing or modifying Go code in this repository. Enforces modern Go usage, project language assumptions, and localization rules for comments and error messages."
name: "Go Modern Style with Chinese Errors"
applyTo: "**/*.go"
---
# Go Project Instruction

This repository is a Go project. Use English for all assistant explanations and instruction text, while following the code rules below.

## Scope

- Applies to all Go source files in this repository.
- Prefer consistency with existing repository architecture, naming, and package boundaries.

## Language And Version Rules

- Treat Go as the default implementation language unless the task explicitly requires another language.
- Prefer modern Go syntax and standard-library capabilities available to the project version.
- When multiple approaches are valid, choose the clearer and more idiomatic modern Go approach.
- Use `require` and `assert` packages for testing, avoid directly calling `t.Xxx`.
- Use `errors` package for error creation and wrapping.

## Modern syntax examples

```go
// correct
for i := range 100 {}

// avoid
for i := 0; i < 100; i++ {}
```

## Comment And Error Message Localization

- Write code comments in Simplified Chinese.
- Write developer-facing and user-facing and testing error messages in Simplified Chinese.
- Keep logging and error text concise, actionable, and context-rich.

## Practical Exceptions

- Do not translate third-party error strings when wrapping; preserve original errors with `%w` and add Chinese context around them.
- Keep protocol fields, external API keys, environment variable names, and machine-readable constants unchanged.

## Example Patterns

```go
// 正确: 注释使用中文，错误上下文使用中文，并保留原始错误
if err := doWork(); err != nil {
    return errors.Errorf("执行任务失败: %w", err)
}

// 正确: 错误信息清晰且可定位
return errors.Errorf("读取配置文件失败，路径=%s: %w", path, err)
```

```go
// 避免: 英文注释和英文错误信息
// load config from file
return errors.Errorf("failed to load config: %w", err)
```
