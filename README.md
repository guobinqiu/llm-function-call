# LLM Function Call Demo

## 不同模型的 Function Call 调用方式

| 模型/厂商                           | 调用方式字段                                                         | 支持的结构             | 示例                                 |
| ----------------------------------- | -------------------------------------------------------------------- | ---------------------- | ------------------------------------ |
| **OpenAI GPT（如 gpt-4, gpt-3.5）** | `functions` + `function_call`（旧）<br>`tools` + `tool_choice`（新） | ✅ 支持 `function` 类型 | ✅ 官方推荐使用 `tools`（兼容新特性） |
| **DeepSeek（如 deepseek-chat）**    | `tools` + `tool_choice`                                              | ✅ 支持 `function` 类型 | ✅ DeepSeek 用的是 OpenAI 新标准      |
| **Claude（Anthropic）**             | ❌ 不兼容 OpenAI 的 Function Call                                     | ❌ 完全不同结构         | ❌ 无法直接对齐                       |
| **Gemini（Google）**                | ❌ 不兼容 OpenAI 标准                                                 | ❌                      | ❌ 不支持 function call               |

## 代码举例

- deepseek 通过 `tools` + `tool_choice` 调用
- openai 通过 `functions` + `function_call` 调用

## 效果

<img width="382" alt="image" src="https://github.com/user-attachments/assets/57b355ea-f8e8-4e46-b0e4-9dfbb8b95a44" />

## 参考

> https://api-docs.deepseek.com/zh-cn/guides/function_calling

> https://platform.openai.com/docs/guides/function-calling?api-mode=responses

