# 🪙 MyCoin - Blockchain Cryptocurrency Platform

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active-green.svg)]()

**MyCoin** là một nền tảng tiền điện tử blockchain hoàn chỉnh được xây dựng bằng Go và JavaScript, sử dụng thuật toán **Proof of Stake (PoS)** để đảm bảo tính bảo mật và hiệu quả năng lượng.

## 📋 Mục lục

- [🌟 Tính năng chính](#-tính-năng-chính)
- [🏗️ Kiến trúc hệ thống](#️-kiến-trúc-hệ-thống)
- [🚀 Cài đặt và chạy](#-cài-đặt-và-chạy)
- [💻 Sử dụng](#-sử-dụng)
- [🔧 API Documentation](#-api-documentation)
- [📁 Cấu trúc dự án](#-cấu-trúc-dự-án)
- [🧪 Testing](#-testing)
- [🤝 Đóng góp](#-đóng-góp)
- [📄 License](#-license)

## 🌟 Tính năng chính

### 💼 Quản lý Ví (Wallet Management)
- ✅ Tạo ví mới với Private Key tự động
- ✅ Import ví bằng Private Key hoặc Passphrase
- ✅ Xem số dư và thông tin ví chi tiết
- ✅ Export backup ví an toàn

### 💸 Giao dịch (Transactions)
- ✅ Gửi MyCoin đến bất kỳ địa chỉ nào
- ✅ Phí giao dịch linh hoạt
- ✅ Preview giao dịch trước khi gửi
- ✅ Lịch sử giao dịch chi tiết với pagination

### ⛏️ Mining & Staking
- ✅ **Proof of Stake (PoS)** consensus algorithm
- ✅ Stake coin để trở thành validator
- ✅ Weighted random validator selection
- ✅ Block rewards cho validators
- ✅ Anti-monopoly với cooldown mechanism

### 📊 Dashboard & Analytics
- ✅ Real-time blockchain statistics
- ✅ Validator network information
- ✅ Transaction pool monitoring
- ✅ Auto-refresh mỗi 30 giây

### 🎨 Giao diện người dùng
- ✅ Web UI modern tương tự MyEtherWallet
- ✅ Transaction history giống Etherscan
- ✅ Responsive design cho mobile
- ✅ Dark/Light theme support

## 🏗️ Kiến trúc hệ thống

```
MyCoinApp/
├── 🖥️ Backend (Go)
│   ├── internal/
│   │   ├── api/          # REST API endpoints
│   │   ├── blockchain/   # Core blockchain logic
│   │   ├── consensus/    # PoS consensus mechanism
│   │   ├── wallet/       # Wallet cryptography
│   │   ├── pool/         # Transaction pool
│   │   └── models/       # Data structures
│   ├── config/           # Configuration management
│   └── cmd/              # Application entry point
│
├── 🌐 Frontend (HTML/CSS/JS)
│   ├── web/static/
│   │   ├── js/           # Application logic
│   │   ├── css/          # Styling
│   │   └── index.html    # Main UI
│
└── 📝 Data Storage
    └── blockchain.json   # Persistent blockchain data
```

## 🚀 Cài đặt và chạy

### Yêu cầu hệ thống
- **Go 1.21+** 
- **Git**
- **Web Browser** (Chrome, Firefox, Safari)

### 1️⃣ Clone repository
```bash
git clone https://github.com/yourusername/MyCoinApp.git
cd MyCoinApp
```

### 2️⃣ Cài đặt dependencies
```bash
go mod tidy
```

### 3️⃣ Cấu hình (tùy chọn)
Chỉnh sửa file `config/config.go` nếu cần:
```go
type Config struct {
    Port                 string  `default:":8080"`
    InitialWalletBalance float64 `default:"100.0"`
    // ... other configs
}
```

### 4️⃣ Chạy ứng dụng
```bash
# Development mode
go run cmd/main.go

# Hoặc build và chạy
go build -o mycoin cmd/main.go
./mycoin
```

### 5️⃣ Truy cập Web UI
Mở trình duyệt và truy cập: **http://localhost:8080**

## 💻 Sử dụng

### 🔐 Tạo ví đầu tiên
1. Vào tab **"Ví (Wallet)"**
2. Click **"Tạo ví mới"**
3. Lưu lại Private Key một cách an toàn
4. Ví sẽ được cấp **100 MYC** miễn phí

### 💰 Stake để trở thành Validator
1. Vào tab **"Staking"**
2. Nhập địa chỉ ví và số lượng stake (tối thiểu 10 MYC)
3. Nhập Private Key để xác thực
4. Click **"Stake Coins"**

### ⛏️ Mining Blocks
1. Vào tab **"Mining"**
2. Nhập địa chỉ validator
3. Click **"Mine Block"**
4. Hệ thống sẽ kiểm tra bạn có được chọn làm validator không

### 💸 Gửi giao dịch
1. Vào tab **"Gửi Coin"**
2. Nhập địa chỉ người nhận và số lượng
3. Thiết lập phí giao dịch (mặc định 0.01 MYC)
4. Xem preview và xác nhận

### 📊 Xem lịch sử giao dịch
1. Vào tab **"Lịch Sử"**
2. Nhập địa chỉ cần tìm kiếm
3. Xem bảng giao dịch chi tiết với pagination

## 🔧 API Documentation

### Wallet APIs
```http
POST /api/wallet/create
POST /api/wallet/import
GET  /api/wallet/balance/:address
```

### Transaction APIs
```http
POST /api/transaction/send
GET  /api/transaction/history/:address
```

### Blockchain APIs
```http
POST /api/blockchain/mine
GET  /api/blockchain/info
```

### Staking APIs
```http
POST /api/staking/stake
POST /api/staking/unstake
GET  /api/staking/validators
GET  /api/staking/validator/:address
GET  /api/staking/info
```

### Ví dụ API Call
```javascript
// Tạo ví mới
const response = await fetch('/api/wallet/create', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
});

// Gửi giao dịch
const txResponse = await fetch('/api/transaction/send', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        from: "0x123...",
        to: "0x456...",
        amount: 10.5,
        fee: 0.01,
        private_key: "abc123..."
    })
});
```

## 📁 Cấu trúc dự án

```
MyCoinApp/
├── 📂 cmd/
│   └── main.go                 # Entry point
├── 📂 config/
│   └── config.go              # App configuration
├── 📂 internal/
│   ├── 📂 api/
│   │   └── handle.go          # HTTP handlers
│   ├── 📂 blockchain/
│   │   ├── blockchain.go      # Core blockchain
│   │   └── block.go           # Block structure
│   ├── 📂 consensus/
│   │   └── pos.go             # Proof of Stake
│   ├── 📂 models/
│   │   └── models.go          # Data structures
│   ├── 📂 pool/
│   │   └── transaction.go     # TX pool
│   └── 📂 wallet/
│       └── wallet.go          # Cryptography
├── 📂 web/static/
│   ├── 📂 js/
│   │   ├── app.js             # Main application
│   │   ├── api.js             # API client
│   │   ├── config.js          # Frontend config
│   │   └── utils.js           # Utilities
│   ├── 📂 css/
│   │   └── style.css          # Styling
│   └── index.html             # Main UI
├── 📄 blockchain.json          # Blockchain data
├── 📄 go.mod                  # Go modules
├── 📄 go.sum                  # Dependencies
└── 📄 README.md              # Documentation
```

## 🔒 Bảo mật

- ✅ **Private Key Encryption**: ECDSA với secp256k1
- ✅ **Transaction Signing**: Cryptographic signatures
- ✅ **Input Validation**: Server-side validation mọi endpoint
- ✅ **CORS Protection**: Configured CORS headers
- ✅ **Rate Limiting**: Cooldown mechanism cho validators

## ⚡ Performance

- ✅ **Concurrent Processing**: Goroutines cho block mining
- ✅ **Caching**: In-memory balance caching
- ✅ **Pagination**: Transaction history pagination
- ✅ **Auto-refresh**: Efficient real-time updates

## 🛠️ Development

### Thêm feature mới
1. Fork repository
2. Tạo feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Mở Pull Request

### Code style
- Tuân theo Go conventions
- Comment code rõ ràng
- Error handling đầy đủ
- Unit tests cho functions quan trọng

## 🤝 Đóng góp

Chúng tôi hoan nghênh mọi đóng góp! Vui lòng đọc [CONTRIBUTING.md](CONTRIBUTING.md) để biết thêm chi tiết.

### Contributors
- 👤 **KitDev** - *Initial work* - [@kietnguyen2003](https://github.com/kietnguyen2003)

## 📄 License

Dự án này được cấp phép dưới MIT License - xem file [LICENSE](LICENSE) để biết chi tiết.

## 📞 Liên hệ

- 📧 Email: ngkiet2611@gmail.com
- 🐛 Issues: [GitHub Issues](https://github.com/kietnguyen2003/MyCoin-App-Blockchain/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/kietnguyen2003/MyCoin-App-Blockchain/discussions)

## 🚀 Roadmap

### Version 1.1
- [ ] Multi-signature wallets
- [ ] Smart contracts support
- [ ] Mobile app (React Native)

### Version 1.2  
- [ ] Cross-chain bridges
- [ ] NFT marketplace
- [ ] Advanced analytics

---

<div align="center">

**⭐ Nếu dự án hữu ích, hãy cho chúng tôi một star! ⭐**

Made with ❤️ by KitDev

</div>