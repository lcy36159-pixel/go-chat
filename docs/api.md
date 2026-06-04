# go-chat API 說明文件

Base URL: `http://localhost:8080`

除 `POST /register` 與 `POST /login` 外，所有 API（含 WebSocket 握手）都需要 JWT 驗證：

`Authorization: Bearer <token>`

伺服器需設定環境變數：

- `JWT_SECRET`：JWT 簽章密鑰（必填，至少 32 字元）
- `JWT_TTL_HOURS`：token 有效小時數（選填，預設 24）

---

## 目錄

| 方法 | 路徑              | 說明                 |
| ---- | ----------------- | -------------------- |
| POST | `/register`       | 註冊新使用者         |
| POST | `/login`          | 使用者登入並取得 JWT |
| POST | `/chats/private`  | 建立私人聊天室       |
| POST | `/chats/group`    | 建立群組聊天室       |
| GET  | `/chats`          | 取得使用者聊天室列表 |
| POST | `/chats/:id/read` | 標記已讀             |
| GET  | `/messages`       | 取得聊天訊息（分頁） |
| GET  | `/ws`             | WebSocket 即時通訊   |

---

## 1. 註冊

**`POST /register`**

建立新使用者帳號。

### 驗證

此 API 為公開端點，不需要 JWT。

### Request Body（JSON）

| 欄位       | 類型   | 必填 | 說明     |
| ---------- | ------ | ---- | -------- |
| `account`  | string | ✅    | 登入帳號（唯一） |
| `username` | string | ✅    | 顯示名稱 |
| `password` | string | ✅    | 密碼（至少 6 碼） |

```json
{
  "account": "alice001",
  "username": "alice",
  "password": "secret123"
}
```

### Response

**成功 `201 Created`**

```json
{
  "user_id": 1
}
```

**錯誤**

| HTTP 狀態                   | 說明                                       |
| --------------------------- | ------------------------------------------ |
| `400 Bad Request`           | body 格式錯誤、帳密缺失、或密碼長度不足 6 碼 |
| `409 Conflict`              | `account` 已存在                           |
| `500 Internal Server Error` | 註冊失敗                                   |

---

## 2. 登入

**`POST /login`**

使用帳號密碼登入，成功後回傳 JWT token。

### 驗證

此 API 為公開端點，不需要 JWT。

### Request Body（JSON）

| 欄位       | 類型   | 必填 | 說明     |
| ---------- | ------ | ---- | -------- |
| `account`  | string | ✅    | 登入帳號 |
| `password` | string | ✅    | 密碼     |

```json
{
  "account": "alice001",
  "password": "secret123"
}
```

### Response

**成功 `200 OK`**

```json
{
  "token": "<jwt_token>"
}
```

**錯誤**

| HTTP 狀態                   | 說明                            |
| --------------------------- | ------------------------------- |
| `400 Bad Request`           | body 格式錯誤或帳密缺失         |
| `401 Unauthorized`          | account 或密碼錯誤              |
| `500 Internal Server Error` | 登入失敗                        |

---

## 3. 建立私人聊天室

**`POST /chats/private`**

建立兩人私人聊天室。若已存在相同成員的私人聊天室，直接回傳既有的 `chat_id`。

### 驗證

發起者身分由 JWT 的 `user_id` 決定。

### Request Body（JSON）

| 欄位             | 類型 | 必填 | 說明            |
| ---------------- | ---- | ---- | --------------- |
| `target_user_id` | uint | ✅    | 對方的使用者 ID |

```json
{
  "target_user_id": 2
}
```

### Response

**成功 `200 OK`**

```json
{
  "chat_id": 10
}
```

**錯誤**

| HTTP 狀態                   | 說明                             |
| --------------------------- | -------------------------------- |
| `400 Bad Request`           | request body 格式錯誤            |
| `401 Unauthorized`          | token 缺失或無效                 |
| `500 Internal Server Error` | 建立失敗（例：與自己建立聊天室） |

---

## 4. 建立群組聊天室

**`POST /chats/group`**

建立多人群組聊天室。發起者會自動被加入群組。

### 驗證

發起者身分由 JWT 的 `user_id` 決定。

### Request Body（JSON）

| 欄位       | 類型   | 必填 | 說明                                           |
| ---------- | ------ | ---- | ---------------------------------------------- |
| `name`     | string | ✅    | 群組名稱                                       |
| `user_ids` | []uint | ✅    | 要加入的使用者 ID 陣列（可重複，系統自動去重） |

```json
{
  "name": "專案討論群",
  "user_ids": [2, 3, 4]
}
```

### Response

**成功 `201 Created`**

```json
{
  "chat_id": 11
}
```

**錯誤**

| HTTP 狀態                   | 說明                                        |
| --------------------------- | ------------------------------------------- |
| `400 Bad Request`           | `name` 為空或 body 格式錯誤                  |
| `401 Unauthorized`          | token 缺失或無效                              |
| `500 Internal Server Error` | 資料庫建立失敗                              |

---

## 5. 取得使用者聊天室列表

**`GET /chats`**

取得指定使用者的所有聊天室，包含最後一則訊息與未讀數。依最後訊息時間降序排列。

### 驗證

使用者身分由 JWT 的 `user_id` 決定。

### Response

**成功 `200 OK`**

```json
[
  {
    "ChatID": 10,
    "Name": "Alice",
    "LastMessage": "今天幾點開會？",
    "UpdatedAt": "2026-04-27T14:00:00Z",
    "UnreadCount": 3
  },
  {
    "ChatID": 11,
    "Name": "專案討論群",
    "LastMessage": "已上線！",
    "UpdatedAt": "2026-04-27T13:30:00Z",
    "UnreadCount": 0
  }
]
```

**欄位說明**

| 欄位          | 類型              | 說明                                   |
| ------------- | ----------------- | -------------------------------------- |
| `ChatID`      | uint              | 聊天室 ID                              |
| `Name`        | string            | 私人聊天顯示對方名稱；群組顯示群組名稱 |
| `LastMessage` | string            | 最後一則訊息內容                       |
| `UpdatedAt`   | string (ISO 8601) | 最後一則訊息時間                       |
| `UnreadCount` | int               | 未讀訊息數（不含自己發出的訊息）       |

**錯誤**

| HTTP 狀態                   | 說明           |
| --------------------------- | -------------- |
| `401 Unauthorized`          | token 缺失或無效 |
| `500 Internal Server Error` | 查詢失敗       |

---

## 6. 標記已讀

**`POST /chats/:id/read`**

更新指定使用者在某聊天室的最後已讀訊息 ID，進而重置未讀計數。

### Path 參數

| 名稱 | 類型 | 必填 | 說明      |
| ---- | ---- | ---- | --------- |
| `id` | uint | ✅    | 聊天室 ID |

### 驗證

使用者身分由 JWT 的 `user_id` 決定。

### Request Body（JSON）

| 欄位                   | 類型 | 必填 | 說明                                  |
| ---------------------- | ---- | ---- | ------------------------------------- |
| `last_read_message_id` | uint | ✅    | 最後已讀的訊息 ID（必須屬於該聊天室） |

```json
{
  "last_read_message_id": 42
}
```

### Response

**成功 `200 OK`**

```json
{
  "ok": true
}
```

**錯誤**

| HTTP 狀態                   | 說明                                                    |
| --------------------------- | ------------------------------------------------------- |
| `400 Bad Request`           | `id` 格式錯誤，或 `last_read_message_id` 缺少           |
| `401 Unauthorized`          | token 缺失或無效                                          |
| `403 Forbidden`             | 使用者不是該聊天室成員，或訊息不屬於該聊天室            |
| `500 Internal Server Error` | 資料庫更新失敗                                          |

---

## 7. 取得聊天訊息（分頁）

**`GET /messages`**

取得指定聊天室的訊息，每次最多回傳 20 筆，依 ID 降序（最新在前）。使用 `last_id` 做游標分頁。

### Query 參數

| 名稱      | 類型 | 必填 | 說明                                                            |
| --------- | ---- | ---- | --------------------------------------------------------------- |
| `chat_id` | uint | ✅    | 聊天室 ID                                                       |
| `last_id` | uint | ❌    | 上一頁最後一筆訊息的 ID（用於載入更早訊息），預設為 0（取最新） |

### Response

**成功 `200 OK`**

```json
[
  {
    "ID": 55,
    "ChatID": 10,
    "SenderID": 3,
    "Type": "text",
    "Content": "明天見！",
    "CreatedAt": "2026-04-27T14:10:00Z",
    "DeletedAt": null
  },
  {
    "ID": 54,
    "ChatID": 10,
    "SenderID": 1,
    "Type": "text",
    "Content": "好的",
    "CreatedAt": "2026-04-27T14:09:00Z",
    "DeletedAt": null
  }
]
```

**欄位說明**

| 欄位        | 類型              | 說明                               |
| ----------- | ----------------- | ---------------------------------- |
| `ID`        | uint              | 訊息 ID                            |
| `ChatID`    | uint              | 所屬聊天室 ID                      |
| `SenderID`  | uint \| null      | 發送者 ID；若訊息已被刪除設為 null |
| `Type`      | string            | 訊息類型，目前固定為 `"text"`      |
| `Content`   | string            | 訊息內容                           |
| `CreatedAt` | string (ISO 8601) | 發送時間                           |
| `DeletedAt` | string \| null    | 軟刪除時間（非 null 代表已收回）   |

**分頁範例**

```
# 第一頁（最新 20 筆）
GET /messages?chat_id=10

# 第二頁（取比 ID=36 更早的訊息）
GET /messages?chat_id=10&last_id=36
```

**錯誤**

| HTTP 狀態                   | 說明                            |
| --------------------------- | ------------------------------- |
| `400 Bad Request`           | `chat_id` 或 `last_id` 格式錯誤 |
| `500 Internal Server Error` | 查詢失敗                        |

---

## 8. WebSocket 即時通訊

**`GET /ws`**

建立 WebSocket 長連線，用於即時發送與接收訊息。

### 連線升級（HTTP → WebSocket）

使用標準 `Upgrade: websocket` 握手，框架採用 [gorilla/websocket](https://github.com/gorilla/websocket)。

### 驗證

連線使用者身分由 JWT 的 `user_id` 決定。

**連線範例**

`ws://localhost:8080/ws`

---

### 發送訊息（Client → Server）

連線建立後，Client 發送 JSON 格式文字訊息：

```json
{
  "chat_id": 10,
  "content": "Hello!"
}
```

| 欄位      | 類型   | 必填 | 說明                                  |
| --------- | ------ | ---- | ------------------------------------- |
| `chat_id` | uint   | ✅    | 目標聊天室 ID（必須是該使用者的成員） |
| `content` | string | ✅    | 訊息文字內容                          |

> ⚠️ 若使用者不是該 `chat_id` 的成員，訊息會被靜默丟棄（不斷線）。

---

### 接收訊息（Server → Client）

訊息廣播給同一聊天室的所有在線成員，格式為 `Message` 物件的 JSON：

```json
{
  "ID": 56,
  "ChatID": 10,
  "SenderID": 1,
  "Type": "text",
  "Content": "Hello!",
  "CreatedAt": "2026-04-27T14:15:00Z",
  "DeletedAt": null
}
```

| 欄位        | 類型              | 說明                        |
| ----------- | ----------------- | --------------------------- |
| `ID`        | uint              | 訊息 ID（已存入資料庫）     |
| `ChatID`    | uint              | 聊天室 ID                   |
| `SenderID`  | uint              | 發送者 ID                   |
| `Type`      | string            | 固定為 `"text"`             |
| `Content`   | string            | 訊息內容                    |
| `CreatedAt` | string (ISO 8601) | 發送時間                    |
| `DeletedAt` | null              | 軟刪除欄位（新訊息為 null） |

---

## 錯誤格式

所有 HTTP API 的錯誤回應皆使用統一 JSON 格式：

```json
{
  "error": "錯誤說明"
}
```
