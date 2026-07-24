# 新增 .codegraph/.gitignore

- **日期**：2026-06-17
- **提交**：`10adc6fb` chore: 新增.gitignore
- **类型**：工程配置

## 变更内容

新增 `.codegraph/.gitignore`（16 行），将 codegraph 索引工具产生的本地缓存/索引文件排除在版本控制之外。

## 注意

`.codegraph/` 是本地代码索引工具的工作目录，属于本地开发产物，后续合并回上游时该文件可以保留也可以按需移除，不影响业务代码。
