// Configuration for MyCoin Web App
const CONFIG = {
    // API Configuration
    API_BASE_URL: 'http://localhost:8080',
    API_KEY: '', // Will be loaded from server or localStorage
    
    // Network Configuration
    NETWORK: 'testnet',
    
    // UI Configuration
    REFRESH_INTERVAL: 30000, // 30 seconds
    TOAST_DURATION: 5000, // 5 seconds
    
    // Blockchain Configuration
    MINING_REWARD: 10.0,
    INITIAL_WALLET_BALANCE: 1000.0,
    MIN_STAKE: 100.0,
    
    // Transaction Configuration
    DEFAULT_FEE: 0.01,
    MIN_AMOUNT: 0.01,
    
    // Local Storage Keys
    STORAGE_KEYS: {
        CURRENT_WALLET: 'mycoin_current_wallet',
        API_KEY: 'mycoin_api_key',
        SETTINGS: 'mycoin_settings'
    },
    
    // Error Messages
    ERRORS: {
        NETWORK_ERROR: 'Lỗi kết nối mạng',
        INVALID_ADDRESS: 'Địa chỉ ví không hợp lệ',
        INSUFFICIENT_BALANCE: 'Số dư không đủ',
        INVALID_AMOUNT: 'Số lượng không hợp lệ',
        PRIVATE_KEY_REQUIRED: 'Vui lòng nhập Private Key',
        WALLET_NOT_FOUND: 'Không tìm thấy ví',
        SEND_TRANSACTION_REQUIRED: 'Vui lòng điền đầy đủ thông tin giao dịch',
        NO_WALLET: 'Vui lòng tạo hoặc import ví trước',
        MINER_ADDRESS_REQUIRED: 'Vui lòng nhập địa chỉ miner',
        SEARCH_ADDRESS_REQUIRED: 'Vui lòng nhập địa chỉ để tìm kiếm',
        STAKE_REQUIRED: 'Vui lòng nhập địa chỉ và số lượng để stake',
        UNSTAKE_ADDRESS_REQUIRED: 'Vui lòng nhập địa chỉ để unstake',
        IMPORT_WALLET_REQUIRED: 'Vui lòng nhập Private Key hoặc Passphrase'
    },
    
    // Success Messages
    SUCCESS: {
        WALLET_CREATED: 'Tạo ví thành công',
        WALLET_IMPORTED: 'Import ví thành công',
        TRANSACTION_SENT: 'Gửi giao dịch thành công',
        BLOCK_MINED: 'Tạo khối thành công',
        COINS_STAKED: 'Stake coin thành công',
        COINS_UNSTAKED: 'Unstake coin thành công'
    }
};

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = CONFIG;
}