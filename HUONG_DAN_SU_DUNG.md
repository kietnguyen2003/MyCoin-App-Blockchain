# 📱 HƯỚNG DẪN SỬ DỤNG MYCOIN

## 🚀 Giới thiệu

MyCoin là một hệ thống tiền điện tử blockchain hoàn chỉnh sử dụng thuật toán **Proof of Stake (PoS)**. Hệ thống bao gồm:

- 🔐 **Quản lý ví** (Wallet Management)
- 💸 **Gửi/nhận coin** (Send/Receive)
- 📊 **Xem lịch sử giao dịch** (Transaction History)
- ⛏️ **Mining & Staking** (Consensus)
- 📈 **Dashboard thống kê** (Statistics)

## 🛠️ Cài đặt và khởi chạy

### 1. Khởi động Backend (Go Server)

```bash
# Di chuyển vào thư mục dự án
cd MyCoin

# Khởi động server
go run cmd/main.go
```

Server sẽ chạy tại: `http://localhost:8080`

### 2. Truy cập Web Interface

Mở trình duyệt và truy cập: `http://localhost:8080`

## 📋 Hướng dẫn sử dụng từng chức năng

### 🏠 1. DASHBOARD

**Mục đích**: Xem tổng quan thông tin blockchain và ví cá nhân

**Cách sử dụng**:
- Dashboard tự động load khi mở website
- Hiển thị 4 thông số chính:
  - 💰 **Số dư ví**: Số MYC trong ví hiện tại
  - 🧱 **Số khối**: Tổng số block trong blockchain
  - ⏳ **Giao dịch chờ**: Số transaction chưa được xác nhận
  - 🪙 **Coin đã stake**: Số MYC đang stake làm validator

- **Giao dịch gần đây**: 5 transaction mới nhất của ví
- **Thông tin blockchain**: Mining reward, consensus type, network status
- **Auto refresh**: Tự động cập nhật mỗi 30 giây

### 💼 2. QUẢN LÝ VÍ (WALLET)

#### 2.1 Tạo ví mới

**Bước 1**: Click tab **"Ví (Wallet)"**
**Bước 2**: Click nút **"Tạo ví mới"**
**Bước 3**: Click **"Tạo ví ngẫu nhiên"**

**Kết quả**: Hệ thống tạo ví với:
- ✅ **Address**: Địa chỉ ví công khai
- ✅ **Public Key**: Khóa công khai
- ✅ **Private Key**: Khóa riêng tư (BẢO MẬT)
- ✅ **Passphrase**: Cụm từ khôi phục 12 từ
- ✅ **1000 MYC miễn phí** từ hệ thống

#### 2.2 Import ví từ Private Key

**Bước 1**: Click **"Import ví"**
**Bước 2**: Chọn **"Private Key"**
**Bước 3**: Nhập private key vào ô
**Bước 4**: Click **"Import ví"**

#### 2.3 Import ví từ Passphrase

**Bước 1**: Click **"Import ví"**
**Bước 2**: Chọn **"Passphrase"**
**Bước 3**: Nhập 12 từ passphrase
**Bước 4**: Click **"Import ví"**

#### 2.4 Bảo mật thông tin ví

- 👁️ **Ẩn/hiện** Private Key và Passphrase
- 📋 **Copy** thông tin vào clipboard
- 💾 **Export Backup**: Tải file .json backup
- ⚠️ **Lưu ý bảo mật**: KHÔNG chia sẻ Private Key với ai

### 💸 3. GỬI COIN

**Mục đích**: Chuyển MYC cho người khác

**Bước 1**: Click tab **"Gửi Coin"**
**Bước 2**: Điền thông tin:
- **Từ ví**: Địa chỉ ví gửi (auto-fill từ ví hiện tại)
- **Private Key**: Private key của ví gửi (auto-fill)
- **Đến ví**: Địa chỉ ví nhận
- **Số lượng**: Số MYC muốn gửi
- **Phí giao dịch**: Phí transaction (mặc định 0.01 MYC)

**Bước 3**: Kiểm tra **"Xem trước giao dịch"**
**Bước 4**: Click **"Gửi giao dịch"**

**Lưu ý**: 
- Giao dịch sẽ vào **pending pool** và được xử lý trong 15-30 giây
- Auto-mining sẽ tự động tạo block chứa giao dịch

### 📊 4. LỊCH SỬ GIAO DỊCH

**Mục đích**: Xem tất cả giao dịch của một địa chỉ ví

#### 4.1 Tìm kiếm giao dịch

**Bước 1**: Click tab **"Lịch Sử"**
**Bước 2**: Nhập địa chỉ ví vào ô tìm kiếm
**Bước 3**: Click **"Tìm kiếm"**

**Kết quả hiển thị**:
- 📋 **Bảng giao dịch** với các cột:
  - **Hash**: Mã hash giao dịch (rút gọn)
  - **Khối**: Số block chứa giao dịch
  - **Thời gian**: Thời điểm thực hiện
  - **Từ**: Địa chỉ gửi
  - **Đến**: Địa chỉ nhận
  - **Số lượng**: Số MYC (🟢xanh=nhận, 🔴đỏ=gửi)
  - **Phí**: Phí giao dịch
  - **Trạng thái**: Đã xác nhận

- 📈 **Thống kê tổng hợp**:
  - Tổng số giao dịch
  - Tổng số MYC đã gửi
  - Tổng số MYC đã nhận

- 📄 **Pagination**: Mỗi trang hiển thị 4 giao dịch

### ⛏️ 5. MINING & CONSENSUS

**Mục đích**: Tạo block mới và hiểu về hệ thống PoS

#### 5.1 Thông tin Consensus

- **Loại Consensus**: Proof of Stake (PoS)
- **Block Reward**: 50 MYC cho validator tạo block
- **Auto-mining**: Tự động tạo block mỗi 15 giây
- **Số Validators**: Số người đang stake coin

#### 5.2 Tạo block thủ công

**Bước 1**: Click tab **"Mining"**
**Bước 2**: Nhập địa chỉ ví nhận phần thưởng
**Bước 3**: Click **"Tạo khối thủ công"**

**Kết quả**: 
- Tạo block mới chứa tất cả giao dịch pending
- Nhận 50 MYC block reward
- Block được thêm vào blockchain

### 🪙 6. STAKING

**Mục đích**: Stake coin để trở thành validator và nhận phần thưởng

#### 6.1 Thông tin Staking

- **Minimum Stake**: 100 MYC tối thiểu
- **Reward Rate**: 10% năm + block rewards
- **Max Validators**: Tối đa 100 validator

#### 6.2 Stake coins

**Bước 1**: Click tab **"Staking"**
**Bước 2**: Điền thông tin:
- **Địa chỉ ví**: Ví muốn stake (auto-fill)
- **Private Key**: Private key (auto-fill)
- **Số lượng stake**: Tối thiểu 100 MYC

**Bước 3**: Click **"Stake Coins"**

**Kết quả**:
- Trở thành validator
- Có cơ hội được chọn tạo block
- Nhận block reward (50 MYC) khi tạo block thành công

#### 6.3 Unstake coins

**Bước 1**: Điền địa chỉ validator muốn unstake
**Bước 2**: Nhập private key
**Bước 3**: Click **"Unstake Coins"**

**Kết quả**: 
- Không còn là validator
- Nhận lại số MYC đã stake (trừ penalty nếu có)

#### 6.4 Xem danh sách Validators

Bảng hiển thị:
- **Địa chỉ**: Address của validator
- **Số coin stake**: Số MYC đã stake
- **Thời gian stake**: Khi nào stake
- **Trạng thái**: Hoạt động/Không hoạt động

## 🔧 Các tính năng bổ sung

### Auto-refresh Dashboard
- Tự động cập nhật dashboard mỗi 30 giây
- Hiển thị thông tin real-time

### Toast Notifications
- Thông báo thành công/lỗi cho mọi thao tác
- Auto-dismiss sau 3-5 giây

### Loading States
- Hiển thị spinner khi đang xử lý request
- Prevent double-click trong khi loading

### Responsive Design
- Tương thích mobile/tablet/desktop
- UI tự động điều chỉnh theo kích thước màn hình

### Local Storage
- Lưu thông tin ví trong browser
- Tự động load ví khi reload trang

## ⚠️ Lưu ý quan trọng

### 🔐 Bảo mật
- **Private Key** và **Passphrase** là thông tin CỰC KỲ QUAN TRỌNG
- KHÔNG chia sẻ với bất kỳ ai
- Lưu trữ ở nơi an toàn, tạo backup
- Mất Private Key = mất toàn bộ tài sản

### 💰 Giao dịch
- Luôn kiểm tra địa chỉ ví nhận trước khi gửi
- Phí giao dịch tối thiểu 0.01 MYC
- Giao dịch không thể hoàn tác
- Auto-mining xử lý giao dịch trong 15-30 giây

### 🪙 Staking
- Stake tối thiểu 100 MYC
- Unstake có thể có penalty nếu vi phạm
- Validator có cơ hội được chọn ngẫu nhiên dựa trên số coin stake

## 🆘 Khắc phục sự cố

### Không kết nối được server
```
Lỗi: Cannot connect to server
Giải pháp: Kiểm tra server Go đang chạy tại port 8080
```

### Private key không hợp lệ
```
Lỗi: Invalid private key
Giải pháp: Kiểm tra private key có đúng format hex không
```

### Số dư không đủ
```
Lỗi: Insufficient balance
Giải pháp: Kiểm tra số dư ví có đủ để gửi (amount + fee)
```

### Giao dịch chưa xuất hiện
```
Vấn đề: Giao dịch đã gửi nhưng chưa thấy trong lịch sử
Giải pháp: Đợi 15-30 giây để auto-mining xử lý
```

## 🎯 Demo Scenarios

### Scenario 1: Tạo ví và nhận coin miễn phí
1. Tạo ví mới → Nhận 1000 MYC
2. Check Dashboard → Thấy số dư 1000 MYC
3. Copy địa chỉ ví để nhận thêm coin từ người khác

### Scenario 2: Gửi coin cho bạn bè
1. Có 2 ví: A (1000 MYC) và B (1000 MYC)
2. Ví A gửi 100 MYC cho ví B
3. Check lịch sử: A thấy -100 MYC, B thấy +100 MYC
4. Check số dư: A còn ~899 MYC, B có ~1100 MYC

### Scenario 3: Trở thành validator
1. Có ít nhất 100 MYC trong ví
2. Stake 100 MYC → Trở thành validator
3. Được auto-mining chọn → Nhận 50 MYC reward
4. Unstake → Nhận lại 100 MYC + profits

## 📞 Hỗ trợ

- **Repository**: [GitHub MyCoin](link-to-repo)
- **Issues**: Báo lỗi tại GitHub Issues
- **Documentation**: File README.md
- **Video Demo**: [Link video](link-to-video)

---

🎉 **Chúc bạn sử dụng MyCoin thành công!** 

*Made with ❤️ using Go, JavaScript, and Blockchain Technology*