# 📚 HƯỚNG DẪN SỬ DỤNG MYCOIN BLOCKCHAIN

## 🎯 Giới thiệu
**MyCoin** là nền tảng blockchain cryptocurrency sử dụng thuật toán **Proof of Stake (PoS)**. Hướng dẫn này sẽ giúp bạn sử dụng tất cả tính năng từ cơ bản đến nâng cao.

---

## 🚀 KHỞI ĐỘNG HỆ THỐNG

### 1️⃣ Cài đặt và chạy server
```bash
# Clone project
git clone https://github.com/kietnguyen2003/MyCoin-App-Blockchain.git
cd MyCoinApp

# Cài đặt dependencies
go mod tidy

# Chạy server
go run cmd/main.go
```

### 2️⃣ Truy cập Web Interface
- Mở trình duyệt: **http://localhost:8080**
- Giao diện sẽ hiển thị với 6 tabs chính

---

## 💼 QUẢN LÝ VÍ (WALLET)

### 🔐 Tạo ví mới

#### **Bước 1: Vào tab "Ví (Wallet)"**
![Wallet Tab](docs/images/wallet-tab.png)

#### **Bước 2: Tạo ví mới**
1. Click **"Tạo ví mới"**
2. Hệ thống sẽ generate:
   - **Address**: Địa chỉ ví công khai
   - **Private Key**: Khóa riêng (GIỮ BÍ MẬT!)
   - **Public Key**: Khóa công khai
3. **Quan trọng**: Lưu Private Key an toàn!

```
✅ Ví được tạo thành công!
📧 Address: 0x1a2b3c4d5e6f...
🔑 Private Key: abc123def456... (GIỮ BÍ MẬT!)
💰 Số dư ban đầu: 100.00 MYC
```

#### **Bước 3: Backup ví**
1. Copy Private Key vào notepad
2. Click **"Export Backup"** để tải file JSON
3. Lưu file backup ở nơi an toàn

### 📥 Import ví có sẵn

#### **Cách 1: Import bằng Private Key**
1. Click **"Import ví"**
2. Chọn **"Private Key"**
3. Paste private key vào ô input
4. Click **"Import Wallet"**

#### **Cách 2: Import bằng Passphrase**
1. Click **"Import ví"**
2. Chọn **"Passphrase"**
3. Nhập passphrase (ví dụ: "my secret wallet 2024")
4. Click **"Import Wallet"**

```
✅ Import ví thành công!
📧 Address được khôi phục: 0x7f8e9d0c1b2a...
💰 Số dư hiện tại: 250.50 MYC
```

---

## 📊 DASHBOARD - TỔNG QUAN

### 🏠 Trang chính hiển thị:

#### **Thống kê ví cá nhân**
- 💰 **Số dư hiện tại**: Tổng MYC trong ví
- 📈 **Tổng gửi**: Tổng MYC đã gửi
- 📉 **Tổng nhận**: Tổng MYC đã nhận

#### **Thống kê blockchain**
- 🔗 **Chiều dài chain**: Số blocks hiện tại
- ⏳ **Giao dịch chờ**: Pending transactions
- 💎 **Phần thưởng mining**: Reward cho validators
- ✅ **Trạng thái chain**: Valid/Invalid

#### **Thống kê validator network**
- 👥 **Số validators**: Tổng validators active
- 💰 **Tổng staked**: Tổng MYC được stake
- 📊 **Thống kê staking**: Min stake, rewards, penalties

#### **Giao dịch gần đây**
- 📋 Hiển thị 5 transactions gần nhất
- 🔄 Auto-refresh mỗi 30 giây

---

## 💸 GỬI GIAO DỊCH

### 📤 Quy trình gửi MYC

#### **Bước 1: Vào tab "Gửi Coin"**
![Send Tab](docs/images/send-tab.png)

#### **Bước 2: Điền thông tin**
```
📧 Từ ví: 0x1a2b3c... (tự động điền)
🔑 Private Key: abc123... (tự động điền)
📨 Đến ví: 0x7f8e9d... (nhập địa chỉ người nhận)
💰 Số lượng: 50.5 (MYC muốn gửi)
⚡ Phí giao dịch: 0.01 (mặc định)
```

#### **Bước 3: Xem preview**
```
💰 Số tiền gửi: 50.50 MYC
⚡ Phí giao dịch: 0.01 MYC
📊 Tổng cần trừ: 50.51 MYC
```

#### **Bước 4: Xác nhận gửi**
1. Kiểm tra thông tin
2. Click **"Gửi giao dịch"**
3. Đợi confirmation

```
✅ Giao dịch thành công!
🆔 Hash: 0xabcd1234...
📊 Status: Pending → Confirmed
⏱️ Thời gian: ~10-30 giây
```

### ⚠️ Lưu ý quan trọng
- ✅ Kiểm tra địa chỉ người nhận
- ✅ Đảm bảo đủ số dư (amount + fee)
- ✅ Private key phải đúng
- ⚠️ Giao dịch không thể hoàn tác!

---

## 📈 LỊCH SỬ GIAO DỊCH

### 🔍 Tìm kiếm giao dịch

#### **Bước 1: Vào tab "Lịch Sử"**
![History Tab](docs/images/history-tab.png)

#### **Bước 2: Nhập địa chỉ tìm kiếm**
```
🔍 Nhập địa chỉ: 0x1a2b3c4d5e6f...
🔄 Auto-search cho ví hiện tại
```

#### **Bước 3: Xem bảng giao dịch**

| 🆔 Hash | 🔗 Block | ⏰ Thời gian | 📤 From | 📥 To | 💰 Amount | ⚡ Fee | 📊 Status |
|---------|----------|---------------|---------|-------|-----------|---------|-----------|
| 0xabc... | #125 | 2024-01-15 14:30 | 0x1a2b... | 0x7f8e... | +50.50 | 0.01 | ✅ Confirmed |
| 0xdef... | #124 | 2024-01-15 14:25 | 0x7f8e... | 0x1a2b... | -25.00 | 0.01 | ✅ Confirmed |

### 📊 Thông tin hiển thị:
- **🆔 Hash**: Transaction ID (click để copy)
- **🔗 Block**: Block number chứa transaction
- **⏰ Thời gian**: Timestamp của giao dịch
- **📤 From**: Địa chỉ người gửi
- **📥 To**: Địa chỉ người nhận
- **💰 Amount**: Số tiền (+nhận/-gửi)
- **⚡ Fee**: Phí giao dịch
- **📊 Status**: Trạng thái (Confirmed/Pending)

### 📋 Tính năng bổ sung:
- **📄 Pagination**: 4 giao dịch/trang
- **📊 Thống kê**: Tổng sent/received/transactions
- **🔄 Auto-refresh**: Cập nhật real-time

---

## ⛏️ MINING BLOCKS

### 🏭 Cách thức Mining trong PoS

#### **Hiểu về Proof of Stake**
- 🎯 **Chọn validator**: Hệ thống chọn ngẫu nhiên weighted by stake
- ⚖️ **Stake nhiều = Cơ hội cao hơn** mine được block
- 🕐 **Cooldown**: Mỗi validator phải đợi 60s giữa các blocks
- 💰 **Reward**: 5.0 MYC cho mỗi block mined

#### **Bước 1: Vào tab "Mining"**
![Mining Tab](docs/images/mining-tab.png)

#### **Bước 2: Nhập địa chỉ miner**
```
👤 Miner Address: 0x1a2b3c... (địa chỉ validator của bạn)
```

#### **Bước 3: Click "Mine Block"**
```
🎯 Hệ thống kiểm tra:
✅ Bạn có phải validator không?
✅ Bạn có được chọn mine block này không?
✅ Bạn có trong cooldown period không?
```

#### **Kết quả có thể:**

**✅ Thành công:**
```
🎉 Block mined thành công!
📦 Block Hash: 0xdef456...
📊 Block Index: #126
💰 Reward: 5.0 MYC
📋 Transactions: 3
```

**❌ Thất bại:**
```
⚠️ Bạn không phải validator được chọn để tạo block!
💡 Tip: Hãy thử lại sau, hoặc stake nhiều MYC hơn
```

### 🔍 Thông tin mining hiển thị:
- **🏭 Consensus**: Proof of Stake (PoS)
- **💰 Reward hiện tại**: 5.0 MYC/block
- **👥 Tổng validators**: Số validators active
- **⏳ Pending transactions**: Giao dịch chờ được mine
- **📦 Block gần nhất**: Hash và thời gian

---

## 🏦 STAKING SYSTEM

### 💰 Trở thành Validator

#### **Bước 1: Vào tab "Staking"**
![Staking Tab](docs/images/staking-tab.png)

#### **Bước 2: Stake coins**
```
👤 Địa chỉ: 0x1a2b3c... (tự động điền)
🔑 Private Key: abc123... (tự động điền) 
💰 Số lượng stake: 50.0 (tối thiểu 10.0 MYC)
```

#### **Bước 3: Xác nhận stake**
1. Kiểm tra đủ số dư
2. Click **"Stake Coins"**
3. Đợi confirmation

```
✅ Stake thành công!
👤 Trở thành validator: 0x1a2b3c...
💰 Đã stake: 50.0 MYC
📊 Voting Power: 12.5% (dựa trên tỷ lệ stake)
🎯 Có thể mine blocks ngay!
```

### 📊 Thông tin Validator

#### **Bảng validators hiện tại:**

| 👤 Address | 💰 Staked | ⏰ Join Time | 📊 Status |
|------------|-----------|--------------|-----------|
| 0x1a2b... | 50.0 MYC | 2024-01-15 | 🟢 Hoạt động |
| 0x7f8e... | 100.0 MYC | 2024-01-14 | 🟢 Hoạt động |
| 0x9c0d... | 25.0 MYC | 2024-01-13 | 🔴 Cooldown |

#### **Thống kê staking:**
- **💰 Tổng staked**: 175.0 MYC
- **👥 Active validators**: 3
- **📊 Min stake**: 10.0 MYC
- **💎 Block reward**: 5.0 MYC
- **⚡ Staking reward**: 5% annually

### 🚫 Unstake Coins

#### **Bước 1: Trong tab "Staking"**
```
👤 Địa chỉ unstake: 0x1a2b3c...
🔑 Private Key: abc123...
```

#### **Bước 2: Click "Unstake Coins"**
```
⚠️ Xác nhận unstake?
💰 Sẽ nhận lại: 50.0 MYC
📊 Mất quyền validator
```

#### **Bước 3: Confirm**
```
✅ Unstake thành công!
💰 Đã hoàn trả: 50.0 MYC vào ví
📊 Không còn là validator
🎯 Cần stake lại để mine blocks
```

---

## ⚡ TIPS & TRICKS

### 🎯 Tối ưu hóa Mining

#### **Tăng cơ hội mine:**
1. **💰 Stake nhiều hơn**: Voting power = stake/total_stake
2. **⏰ Tránh cooldown**: Đợi 60s giữa các blocks
3. **🔄 Mine thường xuyên**: Cơ hội được chọn cao hơn

#### **Chiến lược staking:**
```
💡 Chiến lược tốt:
- Stake 50-100 MYC cho cân bằng risk/reward
- Monitor network để thấy competition
- Unstake nếu quá nhiều validators
```

### 🔒 Bảo mật

#### **Bảo vệ Private Key:**
- ✅ **Không chia sẻ** private key
- ✅ **Backup multiple copies** ở nơi an toàn  
- ✅ **Sử dụng hardware wallet** nếu có
- ❌ **Không screenshot** private key
- ❌ **Không lưu trên cloud** không mã hóa

#### **Best practices:**
```
🔒 Security checklist:
- [x] Private key được backup
- [x] Không share cho ai
- [x] Test với số tiền nhỏ trước
- [x] Verify địa chỉ người nhận
- [x] Check network connection
```

### 🚨 Xử lý lỗi thường gặp

#### **❌ "Invalid address"**
```
🔍 Kiểm tra:
- Address có đủ 40+ characters
- Chỉ chứa số và chữ (0-9, a-f)
- Không có spaces hoặc ký tự đặc biệt
```

#### **❌ "Insufficient balance"**
```
💰 Giải pháp:
- Kiểm tra số dư trong Dashboard
- Đảm bảo đủ cho amount + fee
- Mine blocks để kiếm thêm MYC
```

#### **❌ "Validator not selected"**
```
🎯 Lý do:
- Hệ thống random chọn validator khác
- Bạn trong cooldown period (60s)
- Cần stake nhiều hơn để tăng cơ hội
```

#### **❌ "Mining timeout"**
```
⏱️ Nguyên nhân:
- Network chậm
- Quá nhiều pending transactions
- Hệ thống overload

🔧 Giải pháp:
- Thử lại sau 10-30 giây
- Refresh trang
- Kiểm tra network connection
```

---

## 🎮 DEMO SCENARIOS

### 🆕 Scenario 1: Người dùng mới bắt đầu
```
1. 🔐 Tạo ví mới → Nhận 100 MYC miễn phí
2. 📊 Check Dashboard → Xem số dư và thông tin
3. 💸 Gửi 10 MYC cho bạn → Test giao dịch
4. 📈 Xem lịch sử → Theo dõi transaction
```

### 💰 Scenario 2: Trở thành Validator
```
1. 🏦 Có ít nhất 10 MYC trong ví
2. 🪙 Stake 50 MYC → Trở thành validator
3. ⛏️ Mine blocks → Có cơ hội được chọn
4. 💎 Nhận rewards → 5 MYC mỗi block
5. 🚫 Unstake → Rút về khi cần
```

### 🔄 Scenario 3: Giao dịch giữa 2 ví
```
Ví A (100 MYC) → Ví B (100 MYC)
1. A gửi 25 MYC cho B
2. A: 100 - 25 - 0.01 = 74.99 MYC
3. B: 100 + 25 = 125 MYC  
4. Cả 2 đều xem được trong lịch sử
```

### 🌐 Scenario 4: Network với nhiều validators
```
1. 5 người stake: 10, 20, 30, 40, 50 MYC
2. Tổng stake: 150 MYC
3. Cơ hội được chọn:
   - 10 MYC: 6.7%
   - 20 MYC: 13.3% 
   - 30 MYC: 20%
   - 40 MYC: 26.7%
   - 50 MYC: 33.3%
```

---

## 🔧 ADVANCED FEATURES

### 🎛️ Configuration
File `web/static/js/config.js` chứa các settings:
```javascript
CONFIG = {
    API_BASE_URL: 'http://localhost:8080',
    REFRESH_INTERVAL: 30000, // 30s
    DEFAULT_FEE: 0.01,
    MIN_STAKE: 10.0,
    // ...
}
```

### 📱 Mobile Support
- ✅ Responsive design cho phone/tablet
- ✅ Touch-friendly buttons
- ✅ Optimized layouts
- ✅ Portrait/landscape support

### 🔄 Real-time Updates
- 📊 Dashboard auto-refresh mỗi 30s
- 🔔 Toast notifications
- ⚡ Loading states
- 🎯 Status indicators

### 💾 Data Persistence
- 🏦 Blockchain data → `blockchain.json`
- 💼 Wallet info → LocalStorage
- ⚙️ Settings → Browser storage
- 📋 Session data → Memory

---

## 📞 HỖ TRỢ

### 🆘 Khi cần trợ giúp:

#### **📋 Thông tin cần cung cấp:**
- 🖥️ Browser và version
- 📊 Console errors (F12)
- 📋 Steps reproduce lỗi
- 🖼️ Screenshots nếu cần

#### **📞 Liên hệ:**
- 📧 **Email**: ngkiet2611@gmail.com
- 🐛 **GitHub Issues**: [Report Bug](https://github.com/kietnguyen2003/MyCoin-App-Blockchain/issues)
- 💬 **Discussions**: [Ask Question](https://github.com/kietnguyen2003/MyCoin-App-Blockchain/discussions)

### 🔗 Resources hữu ích:
- 📚 **README.md**: Technical documentation
- 🏗️ **DE_BAI.md**: Project requirements
- 💻 **Source Code**: GitHub repository
- 🎥 **Demo Video**: Coming soon...

### 📖 Learning Resources:
- 🌐 **Blockchain Basics**: Ethereum docs
- 🔒 **PoS Consensus**: Proof of Stake explained
- 💻 **Go Programming**: Go documentation
- 🎨 **Web Development**: MDN Web Docs

---

<div align="center">

## 🎉 CHÚC BẠN SỬ DỤNG MYCOIN VUI VẺ!

**Made with ❤️ by KitDev**  
*Blockchain Developer & Crypto Enthusiast*

[![GitHub](https://img.shields.io/badge/GitHub-kietnguyen2003-blue?logo=github)](https://github.com/kietnguyen2003)  
[![Email](https://img.shields.io/badge/Email-ngkiet2611@gmail.com-red?logo=gmail)](mailto:ngkiet2611@gmail.com)

---

**⭐ Star the project if you found it helpful! ⭐**

</div>