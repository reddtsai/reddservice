# ReddService

[![Build Status](https://github.com/reddtsai/reddservice/actions/workflows/test.yml/badge.svg)](https://github.com/reddtsai/reddservice/actions)

從零開始構建微服務，完整實現設計、開發、測試、部署與維運的全過程。

## Document

這裡沒有說明此專案的文件，而是描述個人在開發和維護過程中，準備那些文件，或使用哪些工具來生成文件。

### Swagger

用於描述和記錄 API

### Mermaid

- Sequence Diagram：描述系統內部不同組件之間的交互；展示用戶與系統之間的交互過程。[Example](https://github.com/reddtsai/reddservice/tree/main/docs)

## Test

### Unit Test

每個單元以 Golang 的 package 進行劃分，並在開發和持續集成（CI）階段進行測試。

> 針對單元測試，增加 -tags unittest。

單元測試過程中，使用模擬(Mock)方式與單元外相依物件互動。

> 由 `go generate` 產生 mock file。generate 還不支持 generic，請必免使用。

### Integration Test

沙盒

## CI/CD

### CI

Containerization，透過 github workflow 產生 APP docker container image

### CD

Argo

## Environment

### Docker

由 docker-compose 建立整個系統，主要提供本地開發使用。

### K8s

由 K3d 建立整個系統，主要提供測試使用。

> 由 k3d 建立 k8s cluster，並由 Rancher 管理 cluster，Argo 管理上版。
