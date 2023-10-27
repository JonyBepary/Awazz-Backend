# Awazz-Backend

import repo:
```
    "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
```

proto file compile command:
```
cd internal/model
protoc -I=. --go_out=. ./*.proto

```

```
Awazz-Backend
├─ .devcontainer
│  └─ devcontainer.json
├─ .git
├─ .github
│  └─ workflows
│     └─ go.yml
├─ .gitignore
├─ LICENSE.md
├─ README.md
├─ api
│  ├─ AI API.postman_collection.json
│  ├─ Insomnia_2023-10-26.json
│  └─ Social Media.postman_collection.json
├─ api.go
├─ awazz.go
├─ build
│  └─ Dockerfile
├─ configs
│  ├─ app.conf
│  ├─ auth.conf
│  └─ db.conf
├─ docs
│  └─ docs.md
├─ go.mod
├─ go.sum
├─ internal
│  ├─ durable
│  │  ├─ database.go
│  │  └─ database_test.go
│  ├─ hander
│  ├─ middlewares
│  │  ├─ authentication.go
│  │  ├─ constraint.go
│  │  ├─ context.go
│  │  ├─ logger.go
│  │  └─ state.go
│  └─ model
│     ├─ comments.proto
│     ├─ follows.proto
│     ├─ instance.proto
│     ├─ like.proto
│     ├─ messages.proto
│     ├─ notifications.proto
│     ├─ person.proto
│     └─ post.proto
└─ pkg
   └─ utils.go

```
