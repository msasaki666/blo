# blo

## 準備
```shell
cp dot.env .env
```

## コマンド

### ユーザー

#### 作成
```shell
curl -X POST -F "email=test@example.com" -F "password=test" http://localhost:8080/users/
```

#### トークン取得
```shell
curl -X POST -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"test"}' http://localhost:8080/users/login
```

#### 認証用ヘッダー
```shell
"Authorization: Bearer my_token"
```

### 投稿

#### 作成
```shell
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE2NDQyMzY4MTYsIm9yaWdfaWF0IjoxNjQ0MTUwNDE2fQ.mruQXEnX7IDiuK1U3nN98-tTQBSxjmV4F6qyCpiynI4" -F "title=test_title" -F "content=test_content" http://localhost:8080/posts/
```

#### 更新
```shell
curl -X PUT -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE2NDQyMzY4MTYsIm9yaWdfaWF0IjoxNjQ0MTUwNDE2fQ.mruQXEnX7IDiuK1U3nN98-tTQBSxjmV4F6qyCpiynI4" -F "title=update_title" -F "content=update_content" http://localhost:8080/posts/1
```
