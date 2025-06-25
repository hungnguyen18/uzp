# 📘 `uzp` – CLI Tool for Managing Secrets for Developers

---

## ✅ Mục tiêu dự án

`uzp` là một công cụ dòng lệnh (**CLI tool**) dùng để **lưu trữ và quản lý an toàn** các thông tin nhạy cảm như:

* API keys
* Access tokens
* Secrets cho `.env`
* Service credentials

---

## 🚀 MVP – Minimum Viable Product

### Tính năng bắt buộc cho MVP

| Lệnh                          | Mô tả                                                                 |
| ----------------------------- | --------------------------------------------------------------------- |
| `uzp init`                    | Khởi tạo vault lần đầu (tạo file mã hóa, yêu cầu master password)     |
| `uzp unlock`                  | Nhập master password để mở vault (và giữ unlock trong phiên hiện tại) |
| `uzp lock`                    | Khoá vault thủ công                                                   |
| `uzp add`                     | Thêm secret: `project`, `key`, `value`                                |
| `uzp get <project/key>`       | Truy xuất giá trị từ vault                                            |
| `uzp copy <project/key>`      | Copy giá trị vào clipboard, tự xóa sau TTL (default: 15s)             |
| `uzp list`                    | Xem danh sách các project/key hiện có                                 |
| `uzp search <keyword>`        | Tìm kiếm theo tên key hoặc project                                    |
| `uzp inject --project <name>` | Xuất các secret của project thành `.env` format                       |
| `uzp reset`                   | Xóa sạch vault (yêu cầu xác nhận kỹ)                                  |

---

## 🧱 Kiến trúc hệ thống

```
Terminal ⇄ uzp CLI (Go)
              ↓
         Vault Manager
              ↓
       Encrypted Storage (.uzp.vault)
              ↑
     Encryption Layer (AES-256-GCM)
              ↓
         File System (~/.uzp/)
```

---

## 🛠️ Tech Stack

### 🏗 Ngôn ngữ chính

* ✅ **Go** (Golang): Gọn nhẹ, cross-platform, build binary dễ dàng

### 🔐 Mã hóa

| Mục tiêu                      | Thư viện đề xuất                                         |
| ----------------------------- | -------------------------------------------------------- |
| AES-256-GCM                   | `golang.org/x/crypto` hoặc `crypto/aes`, `crypto/cipher` |
| Derive key từ master password | `scrypt` hoặc `bcrypt` hoặc `argon2`                     |

### 🗃 Lưu trữ

| Cách                                         | Ưu điểm                                             |
| -------------------------------------------- | --------------------------------------------------- |
| **Encrypted JSON File** (`~/.uzp/uzp.vault`) | Đơn giản, dễ migrate, dễ quản lý                    |
| ✳ SQLite + AES (nâng cao)                    | Tùy chọn về sau nếu dữ liệu lớn hoặc cần query mạnh |

> ⚠️ **MVP nên dùng Encrypted JSON File** (được mã hóa toàn bộ) để đơn giản, sau này có thể nâng cấp lên SQLite mã hóa (sử dụng `sqlcipher` hoặc `modernc.org/sqlite` với extension crypto).

### 📋 Clipboard support

| Hệ điều hành                                                                                       | Cách                |
| -------------------------------------------------------------------------------------------------- | ------------------- |
| macOS                                                                                              | `pbcopy`            |
| Linux                                                                                              | `xclip` hoặc `xsel` |
| Windows                                                                                            | `clip.exe`          |
| ➜ Dùng Go thư viện: [`atotto/clipboard`](https://github.com/atotto/clipboard) – hỗ trợ đa nền tảng |                     |

---

## 🧪 Cấu trúc file vault (sau khi giải mã)

```json
{
  "version": 1,
  "projects": {
    "myapp": {
      "firebase_key": "AIzaSyD...",
      "admin_secret": "abcdef..."
    },
    "aws": {
      "access_key": "AKIA...",
      "secret_key": "wxyz..."
    }
  }
}
```

> Cả file này sẽ được mã hóa bằng AES-256-GCM và lưu tại `~/.uzp/uzp.vault`

---

## 🛡️ Bảo mật

| Rủi ro         | Giải pháp                                                  |
| -------------- | ---------------------------------------------------------- |
| Dò mật khẩu    | Sử dụng `scrypt`/`argon2` để derive key, chậm và an toàn   |
| Lộ file vault  | File vault mã hóa toàn bộ bằng key sinh từ master password |
| Clipboard leak | Tự động xoá sau TTL (mặc định 15s)                         |
| Dump memory    | Hạn chế giữ plaintext trong RAM, xóa ngay sau dùng         |
| Brute force    | Giới hạn retry + delay tăng dần                            |

---

## 📁 Cấu trúc thư mục dự án (Go)

```
uzp/
├── cmd/               # CLI entrypoints
│   ├── root.go
│   ├── add.go
│   ├── get.go
│   ├── init.go
│   └── ...
├── internal/
│   ├── crypto/        # Mã hóa, giải mã
│   │   └── crypto.go
│   ├── storage/       # Đọc/ghi file vault
│   │   └── vault.go
│   └── utils/
│       └── clipboard.go
├── go.mod
├── main.go
└── README.md
```

> Dùng `cobra` để build CLI tool đa lệnh (`uzp add`, `uzp get`,...)

---

## 📌 Dev Roadmap sau MVP

| Giai đoạn | Tính năng                                                         |
| --------- | ----------------------------------------------------------------- |
| Post-MVP  | `uzp edit`, `uzp remove`, `uzp backup/restore`, UI hỗ trợ         |
| V2        | Hỗ trợ SQLite + mã hóa, plugin framework (e.g. `uzp export vue`)  |
| V3        | GPG integration, YubiKey unlock, cloud sync (optional), audit log |

---

## ✨ Gợi ý thêm

* **UX tối giản**: Luôn show confirm step, tránh lệnh phá hoại (`uzp reset`, `remove`)
* **Không phụ thuộc cloud**: All data is local and encrypted
* **Fast UX**: Load toàn bộ vault vào memory (giải mã 1 lần), thao tác cực nhanh
