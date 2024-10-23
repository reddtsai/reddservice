註冊用戶

```mermaid
sequenceDiagram
    actor C as Client
    participant G as Gateway
    participant A as Auth
    C->>+G: POST /v1/sign-up
    G->>+A: rpc SignUp
    A->>-G: OK
    G->>-C: OK
```
