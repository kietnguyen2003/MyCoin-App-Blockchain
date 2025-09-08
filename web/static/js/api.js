// API client for MyCoin backend
class MyCoinAPI {
    constructor() {
        this.baseURL = CONFIG.API_BASE_URL;
        this.apiKey = this.getAPIKey();
    }
    
    // Get API key from storage or default
    getAPIKey() {
        return Utils.getStorage(CONFIG.STORAGE_KEYS.API_KEY, '');
    }
    
    // Set API key
    setAPIKey(apiKey) {
        this.apiKey = apiKey;
        Utils.setStorage(CONFIG.STORAGE_KEYS.API_KEY, apiKey);
    }
    
    // Generic HTTP request method
    async request(endpoint, options = {}) {
        const apiPath = this.apiKey ? `/${this.apiKey}` : '';
        const url = `${this.baseURL}/api${apiPath}${endpoint}`;
        
        const defaultOptions = {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        };
        
        const requestOptions = {
            ...defaultOptions,
            ...options,
            headers: {
                ...defaultOptions.headers,
                ...options.headers,
            },
        };
        
        
        try {
            const response = await fetch(url, requestOptions);
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP ${response.status}: ${errorText || response.statusText}`);
            }
            
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('API Request Error:', error);
            console.error('URL was:', url);
            throw error;
        }
    }
    
    // GET request
    async get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }
    
    // POST request
    async post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }
    
    // Wallet API methods
    async createWallet() {
        return this.post('/wallet/create', {});
    }
    
    async importWallet(privateKey = null, passphrase = null) {
        const data = {};
        if (privateKey) data.private_key = privateKey;
        if (passphrase) data.passphrase = passphrase;
        return this.post('/wallet/import', data);
    }
    
    async getBalance(address) {
        return this.get(`/wallet/balance/${address}`);
    }
    
    // Transaction API methods
    async sendTransaction(fromAddress, toAddress, amount, fee, privateKey) {
        const payload = {
            from: fromAddress,
            to: toAddress,
            amount: parseFloat(amount),
            fee: parseFloat(fee),
            private_key: privateKey
        };
        console.log('Sending transaction with payload:', payload);
        return this.post('/transaction/send', payload);
    }
    
    async getTransactionHistory(address) {
        return this.get(`/transaction/history/${address}`);
    }
    
    // Blockchain API methods
    async getBlockchainInfo() {
        return this.get('/blockchain/info');
    }
    
    async getAllBlocks() {
        return this.get('/blockchain/blocks');
    }
    
    async getBlock(index) {
        return this.get(`/blockchain/block/${index}`);
    }
    
    async mineBlock(minerAddress) {
        return this.post('/blockchain/mine', {
            miner_address: minerAddress
        });
    }
    
    // Staking API methods
    async stakeCoins(address, privateKey, amount) {
        return this.post('/staking/stake', {
            address: address,
            private_key: privateKey,
            amount: parseFloat(amount)
        });
    }
    
    async unstakeCoins(address, privateKey) {
        return this.post('/staking/unstake', {
            address: address,
            private_key: privateKey
        });
    }
    
    async getValidators() {
        return this.get('/staking/validators');
    }
    
    async getValidatorInfo(address) {
        return this.get(`/staking/validator/${address}`);
    }
    
    async getStakingInfo() {
        return this.get('/staking/info');
    }
    
    // Helper methods for common operations
    async getFullBlockchainData() {
        try {
            const [blockchainInfo, blocks] = await Promise.all([
                this.getBlockchainInfo(),
                this.getAllBlocks()
            ]);
            
            return {
                info: blockchainInfo,
                blocks: blocks
            };
        } catch (error) {
            console.error('Error fetching blockchain data:', error);
            throw error;
        }
    }
    
    async getWalletSummary(address) {
        try {
            const [balance, history] = await Promise.all([
                this.getBalance(address),
                this.getTransactionHistory(address)
            ]);
            
            // Calculate summary statistics
            let totalSent = 0;
            let totalReceived = 0;
            let totalFees = 0;
            
            history.transactions.forEach(tx => {
                if (tx.from === address) {
                    totalSent += tx.amount;
                    totalFees += tx.fee;
                } else if (tx.to === address) {
                    totalReceived += tx.amount;
                }
            });
            
            return {
                address: address,
                balance: balance.balance,
                totalTransactions: history.transactions.length,
                totalSent: totalSent,
                totalReceived: totalReceived,
                totalFees: totalFees,
                transactions: history.transactions
            };
        } catch (error) {
            console.error('Error fetching wallet summary:', error);
            throw error;
        }
    }
    
    // Validation helpers
    validateAddress(address) {
        console.log('Validating address:', address, 'Length:', address?.length);
        if (!Utils.isValidAddress(address)) {
            console.error('Address validation failed for:', address);
            throw new Error(CONFIG.ERRORS.INVALID_ADDRESS || 'Địa chỉ ví không hợp lệ');
        }
        console.log('Address validation passed');
    }
    
    validatePrivateKey(privateKey) {
        if (!Utils.isValidPrivateKey(privateKey)) {
            throw new Error('Private key không hợp lệ');
        }
    }
    
    validateAmount(amount) {
        if (!Utils.isValidAmount(amount)) {
            throw new Error(CONFIG.ERRORS.INVALID_AMOUNT);
        }
        if (parseFloat(amount) < CONFIG.MIN_AMOUNT) {
            throw new Error(`Số lượng tối thiểu là ${CONFIG.MIN_AMOUNT} MYC`);
        }
    }
    
    // Error handling helper
    handleError(error, operation = 'thao tác') {
        console.error(`Error during ${operation}:`, error);
        
        let message = error.message || CONFIG.ERRORS.NETWORK_ERROR;
        
        // Parse specific error messages from server
        if (message.includes('insufficient balance')) {
            message = CONFIG.ERRORS.INSUFFICIENT_BALANCE;
        } else if (message.includes('invalid address')) {
            message = CONFIG.ERRORS.INVALID_ADDRESS;
        } else if (message.includes('invalid private key')) {
            message = 'Private key không hợp lệ';
        } else if (message.includes('not found')) {
            message = CONFIG.ERRORS.WALLET_NOT_FOUND;
        }
        
        Utils.showToast(`Lỗi ${operation}: ${message}`, 'error');
        return message;
    }
}

// Create global API instance
const api = new MyCoinAPI();

// Auto-detect API key from URL or prompt user
document.addEventListener('DOMContentLoaded', () => {
    const params = Utils.getURLParams();
    if (params.api_key && !api.apiKey) {
        api.setAPIKey(params.api_key);
        // Remove from URL
        Utils.updateURL({ api_key: null });
    }
    
    // If no API key is set, try to get it from user or use empty string
    if (!api.apiKey) {
        // For demo purposes, we'll use empty string
        // In production, you might want to prompt the user
        api.setAPIKey('');
    }
});