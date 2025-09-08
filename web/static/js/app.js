// Main application class for MyCoin
class MyCoinApp {
    constructor() {
        this.currentWallet = null;
        this.refreshInterval = null;
        this.currentPage = 1;
        this.itemsPerPage = 4;
        this.totalTransactions = 0;
        this.allTransactions = [];
        this.init();
    }

    init() {
        this.setupTabNavigation();
        this.setupSendFormListeners();
        this.initializeFormDefaults();
        this.loadWalletFromStorage();
        this.loadDashboard();
        this.startAutoRefresh();
    }

    setupTabNavigation() {
        const navItems = document.querySelectorAll('.nav-item');
        const tabContents = document.querySelectorAll('.tab-content');

        navItems.forEach(item => {
            item.addEventListener('click', () => {
                const tabId = item.getAttribute('data-tab');
                
                // Update navigation
                navItems.forEach(nav => nav.classList.remove('active'));
                item.classList.add('active');
                
                // Update tab content
                tabContents.forEach(tab => tab.classList.remove('active'));
                document.getElementById(tabId).classList.add('active');
                
                // Load tab-specific data
                this.loadTabData(tabId);
            });
        });
    }

    setupSendFormListeners() {
        // Add event listeners for send form to update preview in real-time
        const amountInput = document.getElementById('send-amount');
        const feeInput = document.getElementById('send-fee');
        
        if (amountInput) {
            amountInput.addEventListener('input', () => this.updateTransactionPreview());
        }
        if (feeInput) {
            feeInput.addEventListener('input', () => this.updateTransactionPreview());
        }
        
        // Add event listeners for import wallet radio buttons
        const importRadios = document.querySelectorAll('input[name="import-type"]');
        importRadios.forEach(radio => {
            radio.addEventListener('change', () => this.toggleImportFields());
        });
    }

    initializeFormDefaults() {
        // Set default values from CONFIG
        const sendFeeElement = document.getElementById('send-fee');
        if (sendFeeElement && !sendFeeElement.value) {
            sendFeeElement.value = CONFIG.DEFAULT_FEE;
        }

        // Set placeholder values that reference CONFIG
        const stakeAmountElement = document.getElementById('stake-amount');
        if (stakeAmountElement) {
            stakeAmountElement.placeholder = `Tối thiểu ${CONFIG.MIN_STAKE} MYC`;
            stakeAmountElement.min = CONFIG.MIN_STAKE;
        }
    }
    
    toggleImportFields() {
        const importType = document.querySelector('input[name="import-type"]:checked').value;
        const privateKeyInput = document.getElementById('private-key-input');
        const passphraseInput = document.getElementById('passphrase-input');
        
        if (importType === 'private-key') {
            privateKeyInput.style.display = 'block';
            passphraseInput.style.display = 'none';
            passphraseInput.value = '';
        } else {
            privateKeyInput.style.display = 'none';
            passphraseInput.style.display = 'block';
            privateKeyInput.value = '';
        }
    }

    loadTabData(tabId) {
        switch(tabId) {
            case 'dashboard':
                this.loadDashboard();
                break;
            case 'wallet':
                this.loadWalletData();
                break;
            case 'send':
                this.loadSendData();
                break;
            case 'history':
                this.loadTransactionHistory();
                break;
            case 'mining':
                this.loadMiningData();
                break;
            case 'staking':
                this.loadStakingData();
                break;
        }
    }

    // Dashboard functionality
    async loadDashboard() {
        try {
            Utils.showLoading(true);
            
            const [blockchainInfo, validators] = await Promise.all([
                api.getBlockchainInfo(),
                api.getValidators()
            ]);

            this.updateBlockchainStats(blockchainInfo);
            this.updateValidatorStats(validators);
            
            if (this.currentWallet) {
                const [balance, recentTransactions] = await Promise.all([
                    api.getBalance(this.currentWallet.address),
                    api.getTransactionHistory(this.currentWallet.address)
                ]);
                this.updateWalletInfo(balance);
                this.updateRecentTransactions(recentTransactions);
            } else {
                this.updateRecentTransactions({ transactions: [] });
            }

        } catch (error) {
            api.handleError(error, 'load dashboard');
        } finally {
            Utils.showLoading(false);
        }
    }

    updateBlockchainStats(info) {
        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        updateElement('chain-length', info.chain_length || 0);
        updateElement('pending-transactions', info.pending_transactions || 0);
        updateElement('mining-reward', `${info.mining_reward || 0} MYC`);
        updateElement('difficulty', info.difficulty || 0);
        updateElement('chain-valid', info.is_valid ? 'Hợp lệ' : 'Không hợp lệ');
    }

    updateValidatorStats(validators) {
        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        const validatorCount = validators.validators ? validators.validators.length : 0;
        updateElement('validator-count', validatorCount);
        
        if (validators.validators && validators.validators.length > 0) {
            const totalStaked = validators.validators.reduce((sum, v) => sum + (v.staked_amount || 0), 0);
            updateElement('total-staked', `${totalStaked.toFixed(2)} MYC`);
        }
    }

    updateWalletInfo(balance) {
        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        updateElement('wallet-balance', `${balance.balance.toFixed(2)} MYC`);
    }

    updateRecentTransactions(history) {
        const container = document.getElementById('recent-transactions-list');
        if (!container) return;

        container.innerHTML = '';

        const transactions = history.transactions || [];
        
        if (transactions.length === 0) {
            container.innerHTML = `
                <div class="no-transactions">
                    <i class="fas fa-inbox"></i>
                    <p>Chưa có giao dịch nào</p>
                    <small>Giao dịch của bạn sẽ xuất hiện ở đây</small>
                </div>
            `;
            return;
        }

        // Show only the last 5 transactions
        const recentTransactions = transactions.slice(-5).reverse();

        recentTransactions.forEach(tx => {
            const transactionElement = document.createElement('div');
            transactionElement.className = 'transaction-item';
            
            const isIncoming = this.currentWallet && tx.to === this.currentWallet.address;
            const isOutgoing = this.currentWallet && tx.from === this.currentWallet.address;
            
            let amountClass = '';
            let amountPrefix = '';
            let transactionType = '';
            
            if (isIncoming) {
                amountClass = 'amount-in';
                amountPrefix = '+';
                transactionType = 'Nhận';
            } else if (isOutgoing) {
                amountClass = 'amount-out';
                amountPrefix = '-';
                transactionType = 'Gửi';
            } else {
                transactionType = 'Khác';
            }

            const timestamp = tx.timestamp ? Utils.formatDate(tx.timestamp) : 'N/A';
            const fromAddr = tx.from || 'genesis';
            const toAddr = tx.to || 'N/A';

            transactionElement.innerHTML = `
                <div class="transaction-info">
                    <div class="transaction-type-badge ${amountClass}">
                        <i class="fas ${isIncoming ? 'fa-arrow-down' : 'fa-arrow-up'}"></i>
                        ${transactionType}
                    </div>
                    <div class="transaction-details">
                        <div class="transaction-addresses">
                            <small>Từ: <code>${Utils.formatAddress(fromAddr)}</code></small>
                            <small>Đến: <code>${Utils.formatAddress(toAddr)}</code></small>
                        </div>
                        <div class="transaction-hash">
                            <small>Hash: <code>${Utils.formatHash(tx.hash, 16)}</code></small>
                        </div>
                    </div>
                </div>
                <div class="transaction-amount-info">
                    <div class="transaction-amount ${amountClass}">
                        ${amountPrefix}${tx.amount.toFixed(2)} MYC
                    </div>
                    <div class="transaction-time">
                        <small>${timestamp}</small>
                    </div>
                </div>
            `;
            
            container.appendChild(transactionElement);
        });
    }

    // Wallet functionality
    loadWalletData() {
        // Show/hide wallet info section based on whether user has a wallet
        const walletInfoSection = document.querySelector('.wallet-info');
        const noWalletMessage = document.querySelector('.no-wallet-message');
        
        if (this.currentWallet) {
            if (walletInfoSection) walletInfoSection.style.display = 'block';
            if (noWalletMessage) noWalletMessage.style.display = 'none';
            const updateInputValue = (id, value) => {
                const element = document.getElementById(id);
                if (element) {
                    element.value = value;
                }
            };

            updateInputValue('wallet-address', this.currentWallet.address);
            updateInputValue('wallet-private-key', this.currentWallet.private_key);
            updateInputValue('wallet-public-key', this.currentWallet.public_key);
            if (this.currentWallet.passphrase) {
                updateInputValue('wallet-passphrase', this.currentWallet.passphrase);
            }
            this.loadWalletBalance();
        } else {
            // No wallet - show message
            if (walletInfoSection) walletInfoSection.style.display = 'none';
            if (noWalletMessage) noWalletMessage.style.display = 'block';
        }
    }

    async loadWalletBalance() {
        if (!this.currentWallet) return;
        
        try {
            const balance = await api.getBalance(this.currentWallet.address);
            const element = document.getElementById('current-wallet-balance');
            if (element) {
                element.textContent = `${balance.balance.toFixed(2)} MYC`;
            }
        } catch (error) {
            api.handleError(error, 'load wallet balance');
        }
    }

    showCreateWallet() {
        document.getElementById('create-wallet-form').style.display = 'block';
        document.getElementById('import-wallet-form').style.display = 'none';
    }

    showImportWallet() {
        document.getElementById('create-wallet-form').style.display = 'none';
        document.getElementById('import-wallet-form').style.display = 'block';
        this.toggleImportFields(); // Set initial state
    }

    async createWallet() {
        try {
            Utils.showLoading(true);
            
            const wallet = await api.createWallet();
            
            this.currentWallet = wallet;
            this.saveWalletToStorage();
            
            Utils.showToast(CONFIG.SUCCESS.WALLET_CREATED, 'success');
            
            // Hide form and refresh wallet data
            document.getElementById('create-wallet-form').style.display = 'none';
            this.loadWalletData();
            
        } catch (error) {
            api.handleError(error, 'create wallet');
        } finally {
            Utils.showLoading(false);
        }
    }

    async importWallet() {
        // Check which import method is selected
        const importType = document.querySelector('input[name="import-type"]:checked').value;
        const privateKey = importType === 'private-key' ? 
            document.getElementById('private-key-input').value.trim() : '';
        const passphrase = importType === 'passphrase' ? 
            document.getElementById('passphrase-input').value.trim() : '';

        if (!privateKey && !passphrase) {
            Utils.showToast(CONFIG.ERRORS.IMPORT_WALLET_REQUIRED, 'error');
            return;
        }

        try {
            Utils.showLoading(true);
            
            const wallet = await api.importWallet(privateKey, passphrase);
            
            this.currentWallet = wallet;
            this.saveWalletToStorage();
            
            Utils.showToast(CONFIG.SUCCESS.WALLET_IMPORTED, 'success');
            
            // Clear form and hide it
            document.getElementById('private-key-input').value = '';
            document.getElementById('passphrase-input').value = '';
            document.getElementById('import-wallet-form').style.display = 'none';
            
            this.loadWalletData();
            
        } catch (error) {
            api.handleError(error, 'import wallet');
        } finally {
            Utils.showLoading(false);
        }
    }

    // Send transaction functionality
    loadSendData() {
        if (this.currentWallet) {
            const fromElement = document.getElementById('send-from-address');
            const privateKeyElement = document.getElementById('send-private-key');
            if (fromElement) {
                fromElement.value = this.currentWallet.address;
            }
            if (privateKeyElement) {
                privateKeyElement.value = this.currentWallet.private_key;
            }
            this.loadWalletBalance();
            this.updateTransactionPreview();
        }
    }

    async sendTransaction() {
        const from = document.getElementById('send-from-address').value.trim();
        const to = document.getElementById('send-to-address').value.trim();
        const amount = parseFloat(document.getElementById('send-amount').value);
        const fee = parseFloat(document.getElementById('send-fee').value) || CONFIG.DEFAULT_FEE;

        if (!from || !to || !amount || amount <= 0) {
            Utils.showToast(CONFIG.ERRORS.SEND_TRANSACTION_REQUIRED, 'error');
            return;
        }

        if (!this.currentWallet) {
            Utils.showToast(CONFIG.ERRORS.NO_WALLET, 'error');
            return;
        }

        try {
            api.validateAddress(from);
            api.validateAddress(to);
            api.validateAmount(amount);

            Utils.showLoading(true);

            const result = await api.sendTransaction(from, to, amount, fee, this.currentWallet.private_key);

            Utils.showToast(CONFIG.SUCCESS.TRANSACTION_SENT, 'success');
            
            // Clear form
            document.getElementById('send-to-address').value = '';
            document.getElementById('send-amount').value = '';
            document.getElementById('send-fee').value = CONFIG.DEFAULT_FEE;
            this.updateTransactionPreview();
            
            this.loadWalletBalance();

        } catch (error) {
            api.handleError(error, 'send transaction');
        } finally {
            Utils.showLoading(false);
        }
    }

    updateTransactionPreview() {
        const amount = parseFloat(document.getElementById('send-amount').value) || 0;
        const fee = parseFloat(document.getElementById('send-fee').value) || CONFIG.DEFAULT_FEE;
        const total = amount + fee;

        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        updateElement('preview-amount', `${amount.toFixed(2)} MYC`);
        updateElement('preview-fee', `${fee.toFixed(2)} MYC`);
        updateElement('preview-total', `${total.toFixed(2)} MYC`);
    }

    // Transaction history functionality
    async loadTransactionHistory() {
        const addressElement = document.getElementById('search-address');
        if (addressElement) {
            const address = addressElement.value.trim();
            if (!address && this.currentWallet) {
                addressElement.value = this.currentWallet.address;
                // Auto-search for current wallet
                this.searchTransactions();
            }
        }
    }

    async searchTransactions() {
        const address = document.getElementById('search-address').value.trim();

        if (!address) {
            Utils.showToast(CONFIG.ERRORS.SEARCH_ADDRESS_REQUIRED, 'error');
            return;
        }

        try {
            api.validateAddress(address);
            Utils.showLoading(true);

            const history = await api.getTransactionHistory(address);
            this.allTransactions = history.transactions || [];
            this.totalTransactions = this.allTransactions.length;
            this.currentPage = 1;
            this.displayTransactionHistory();

        } catch (error) {
            api.handleError(error, 'search transactions');
        } finally {
            Utils.showLoading(false);
        }
    }

    displayTransactionHistory() {
        const tbody = document.getElementById('transaction-history-body');
        if (!tbody) return;
        
        tbody.innerHTML = '';

        if (!this.allTransactions || this.allTransactions.length === 0) {
            tbody.innerHTML = '<tr><td colspan="8" class="no-data">Không có giao dịch nào</td></tr>';
            this.updatePagination();
            return;
        }

        // Calculate pagination
        const startIndex = (this.currentPage - 1) * this.itemsPerPage;
        const endIndex = Math.min(startIndex + this.itemsPerPage, this.allTransactions.length);
        const paginatedTransactions = this.allTransactions.slice(startIndex, endIndex);

        const currentAddress = document.getElementById('search-address').value.trim();
        
        paginatedTransactions.forEach(tx => {
            const row = document.createElement('tr');
            const timestamp = tx.timestamp ? Utils.formatDate(tx.timestamp) : 'N/A';
            
            // Debug: log transaction object to see structure
            console.log('Transaction object:', tx);
            
            // Fix status logic - if transaction has been retrieved from blockchain.GetTransactionHistory(), 
            // it means it's already in a block (confirmed). Only pending transactions are in PendingTransactions array
            const status = 'Đã xác nhận';
            const statusClass = 'confirmed';
            
            // Determine transaction type and apply color
            const isIncoming = tx.to === currentAddress;
            const isOutgoing = tx.from === currentAddress;
            
            let amountClass = '';
            let amountPrefix = '';
            if (isIncoming) {
                amountClass = 'amount-in';
                amountPrefix = '+';
            } else if (isOutgoing) {
                amountClass = 'amount-out';
                amountPrefix = '-';
            }
            
            row.innerHTML = `
                <td><code class="hash" title="${tx.hash}">${Utils.formatHash(tx.hash, 12)}</code></td>
                <td><span class="block-number">Block #${tx.block_index !== undefined ? tx.block_index : 'N/A'}</span></td>
                <td>${timestamp}</td>
                <td><code class="address" title="${tx.from}">${Utils.formatAddress(tx.from||"genesis")}</code></td>
                <td><code class="address" title="${tx.to}">${Utils.formatAddress(tx.to)}</code></td>
                <td class="amount ${amountClass}">${amountPrefix}${tx.amount.toFixed(2)} MYC</td>
                <td class="fee">${tx.fee.toFixed(4)} MYC</td>
                <td><span class="status ${statusClass}">${status}</span></td>
            `;
            tbody.appendChild(row);
        });
        
        // Update pagination and stats
        this.updatePagination();
        this.updateTransactionStats(this.allTransactions);
    }

    updateTransactionStats(transactions) {
        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        updateElement('total-transactions', transactions.length);
        
        const currentAddress = document.getElementById('search-address').value.trim();
        let totalSent = 0;
        let totalReceived = 0;
        
        transactions.forEach(tx => {
            if (tx.from === currentAddress) {
                totalSent += tx.amount + tx.fee;
            }
            if (tx.to === currentAddress) {
                totalReceived += tx.amount;
            }
        });
        
        updateElement('total-sent', `${totalSent.toFixed(2)} MYC`);
        updateElement('total-received', `${totalReceived.toFixed(2)} MYC`);
    }

    updatePagination() {
        const paginationControls = document.getElementById('pagination-controls');
        const totalPages = Math.ceil(this.totalTransactions / this.itemsPerPage);
        
        if (this.totalTransactions <= this.itemsPerPage) {
            paginationControls.style.display = 'none';
            return;
        }

        paginationControls.style.display = 'flex';
        
        const startIndex = (this.currentPage - 1) * this.itemsPerPage + 1;
        const endIndex = Math.min(this.currentPage * this.itemsPerPage, this.totalTransactions);
        
        document.getElementById('showing-from').textContent = startIndex;
        document.getElementById('showing-to').textContent = endIndex;
        document.getElementById('total-transactions-count').textContent = this.totalTransactions;
        document.getElementById('current-page').textContent = this.currentPage;
        document.getElementById('total-pages').textContent = totalPages;
        
        document.getElementById('prev-page').disabled = this.currentPage === 1;
        document.getElementById('next-page').disabled = this.currentPage === totalPages;
    }

    changePage(direction) {
        const totalPages = Math.ceil(this.totalTransactions / this.itemsPerPage);
        
        if (direction > 0 && this.currentPage < totalPages) {
            this.currentPage++;
        } else if (direction < 0 && this.currentPage > 1) {
            this.currentPage--;
        }
        
        this.displayTransactionHistory();
    }

    // Mining functionality
    async loadMiningData() {
        if (this.currentWallet) {
            document.getElementById('miner-address').value = this.currentWallet.address;
        }
        
        try {
            const [info, validators] = await Promise.all([
                api.getBlockchainInfo(),
                api.getValidators()
            ]);
            
            const updateElement = (id, content) => {
                const element = document.getElementById(id);
                if (element) {
                    element.textContent = content;
                }
            };
            
            // Update consensus info
            updateElement('current-consensus', 'Proof of Stake (PoS)');
            updateElement('current-mining-reward', `${info.mining_reward || 10.00} MYC`);
            updateElement('total-validators', validators.validators ? validators.validators.length : 0);
            updateElement('mining-pending-tx', info.pending_transactions || 0);
            
            // Update mining stats
            if (info.latest_block) {
                updateElement('latest-block-hash', info.latest_block.hash ? 
                    info.latest_block.hash.substring(0, 16) + '...' : 'N/A');
                updateElement('latest-block-time', info.latest_block.timestamp ? 
                    Utils.formatDate(info.latest_block.timestamp) : 'N/A');
            } else {
                updateElement('latest-block-hash', 'N/A');
                updateElement('latest-block-time', 'N/A');
            }
            
        } catch (error) {
            console.error('Error loading mining data:', error);
        }
    }

    async mineBlock() {
        const minerAddress = document.getElementById('miner-address').value.trim();

        if (!minerAddress) {
            Utils.showToast(CONFIG.ERRORS.MINER_ADDRESS_REQUIRED, 'error');
            return;
        }

        try {
            api.validateAddress(minerAddress);
            Utils.showLoading(true);

            const result = await api.mineBlock(minerAddress);
            
            Utils.showToast(CONFIG.SUCCESS.BLOCK_MINED, 'success');
            
            this.loadMiningData();
            this.loadWalletBalance();

        } catch (error) {
            api.handleError(error, 'mine block');
        } finally {
            Utils.showLoading(false);
        }
    }

    // Staking functionality
    async loadStakingData() {
        try {
            Utils.showLoading(true);
            
            const [validators, stakingInfo] = await Promise.all([
                api.getValidators(),
                api.getStakingInfo()
            ]);

            this.displayValidators(validators.validators);
            this.displayStakingInfo(stakingInfo);

            if (this.currentWallet) {
                const updateInputValue = (id, value) => {
                    const element = document.getElementById(id);
                    if (element) {
                        element.value = value;
                    }
                };
                
                updateInputValue('stake-address', this.currentWallet.address);
                updateInputValue('stake-private-key', this.currentWallet.private_key);
                updateInputValue('unstake-address', this.currentWallet.address);
                updateInputValue('unstake-private-key', this.currentWallet.private_key);
                
                // Check if current wallet is a validator
                const isValidator = validators.validators.some(v => v.address === this.currentWallet.address);
                if (isValidator) {
                    const validatorInfo = await api.getValidatorInfo(this.currentWallet.address);
                    this.displayValidatorInfo(validatorInfo);
                }
            }

        } catch (error) {
            api.handleError(error, 'load staking data');
        } finally {
            Utils.showLoading(false);
        }
    }

    displayValidators(validators) {
        const tbody = document.getElementById('validators-body');
        if (!tbody) return;
        
        tbody.innerHTML = '';

        if (!validators || validators.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="no-data">Không có validator nào</td></tr>';
            return;
        }

        validators.forEach(validator => {
            const row = document.createElement('tr');
            const stakeTime = validator.join_time ? Utils.formatDate(validator.join_time) : 'N/A';
            
            row.innerHTML = `
                <td><code class="address" title="${validator.address}">${Utils.formatAddress(validator.address)}</code></td>
                <td class="amount">${validator.staked_amount.toFixed(2)} MYC</td>
                <td>${stakeTime}</td>
                <td><span class="status ${validator.is_active ? 'active' : 'inactive'}">${validator.is_active ? 'Hoạt động' : 'Không hoạt động'}</span></td>
            `;
            tbody.appendChild(row);
        });
    }

    displayStakingInfo(info) {
        const updateElement = (id, content) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = content;
            }
        };

        updateElement('total-staked-amount', `${info.total_staked.toFixed(2)} MYC`);
        updateElement('active-validators', info.active_validators);
        updateElement('min-stake-amount', `${info.min_stake_amount.toFixed(2)} MYC`);
    }

    displayValidatorInfo(info) {
        const container = document.getElementById('validator-info');
        if (!container) {
            console.log('validator-info element not found');
            return;
        }
        container.innerHTML = `
            <div class="validator-details">
                <h4>Thông tin Validator của bạn:</h4>
                <p><strong>Địa chỉ:</strong> <code>${info.address}</code></p>
                <p><strong>Số coin đã stake:</strong> ${info.staked_amount.toFixed(2)} MYC</p>
                <p><strong>Voting Power:</strong> ${(info.voting_power * 100).toFixed(2)}%</p>
                <p><strong>Trạng thái:</strong> <span class="status ${info.is_active ? 'active' : 'inactive'}">${info.is_active ? 'Hoạt động' : 'Không hoạt động'}</span></p>
            </div>
        `;
        container.style.display = 'block';
    }

    async stakeCoins() {
        console.log('stakeCoins called');
        const address = document.getElementById('stake-address').value.trim();
        const amount = parseFloat(document.getElementById('stake-amount').value);
        console.log('Address:', address, 'Amount:', amount);

        if (!address || !amount || amount <= 0) {
            console.log('Validation failed - missing address or invalid amount');
            Utils.showToast(CONFIG.ERRORS.STAKE_REQUIRED, 'error');
            return;
        }

        if (amount < CONFIG.MIN_STAKE) {
            Utils.showToast(`Số lượng stake tối thiểu là ${CONFIG.MIN_STAKE} MYC`, 'error');
            return;
        }

        if (!this.currentWallet) {
            Utils.showToast(CONFIG.ERRORS.NO_WALLET, 'error');
            return;
        }

        try {
            api.validateAddress(address);
            api.validateAmount(amount);

            Utils.showLoading(true);

            const response = await api.stakeCoins(address, this.currentWallet.private_key, amount);

            Utils.showToast(CONFIG.SUCCESS.COINS_STAKED, 'success');
            
            // Clear form
            document.getElementById('stake-amount').value = '';
            
            this.loadStakingData();
            this.loadWalletBalance();

        } catch (error) {
            api.handleError(error, 'stake coins');
        } finally {
            Utils.showLoading(false);
        }
    }

    async unstakeCoins() {
        const address = document.getElementById('unstake-address').value.trim();

        if (!address) {
            Utils.showToast(CONFIG.ERRORS.UNSTAKE_ADDRESS_REQUIRED, 'error');
            return;
        }

        if (!this.currentWallet) {
            Utils.showToast(CONFIG.ERRORS.NO_WALLET, 'error');
            return;
        }

        try {
            api.validateAddress(address);
            api.validatePrivateKey(this.currentWallet.private_key);

            Utils.showLoading(true);

            const response = await api.unstakeCoins(address, this.currentWallet.private_key);

            Utils.showToast(CONFIG.SUCCESS.COINS_UNSTAKED, 'success');

            this.loadStakingData();
            this.loadWalletBalance();

        } catch (error) {
            api.handleError(error, 'unstake coins');
        } finally {
            Utils.showLoading(false);
        }
    }

    // Auto-refresh functionality
    startAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
        }

        this.refreshInterval = setInterval(() => {
            const activeTab = document.querySelector('.tab-content.active').id;
            if (activeTab === 'dashboard') {
                this.loadDashboard();
            }
        }, CONFIG.REFRESH_INTERVAL);
    }

    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }

    // Wallet storage functionality
    saveWalletToStorage() {
        if (this.currentWallet) {
            localStorage.setItem(CONFIG.STORAGE_KEYS.CURRENT_WALLET, JSON.stringify(this.currentWallet));
        }
    }

    loadWalletFromStorage() {
        const stored = localStorage.getItem(CONFIG.STORAGE_KEYS.CURRENT_WALLET);
        if (stored) {
            try {
                this.currentWallet = JSON.parse(stored);
            } catch (error) {
                console.error('Error loading wallet from storage:', error);
                localStorage.removeItem(CONFIG.STORAGE_KEYS.CURRENT_WALLET);
            }
        }
    }
}

// Global functions for HTML onclick handlers
function showCreateWallet() {
    app.showCreateWallet();
}

function showImportWallet() {
    app.showImportWallet();
}

function createWallet() {
    app.createWallet();
}

function importWallet() {
    app.importWallet();
}

function sendTransaction() {
    app.sendTransaction();
}

function searchTransactions() {
    app.searchTransactions();
}

function mineBlock() {
    app.mineBlock();
}

function stakeCoins() {
    app.stakeCoins();
}

function unstakeCoins() {
    app.unstakeCoins();
}

function refreshDashboard() {
    app.loadDashboard();
}

function changePage(direction) {
    app.changePage(direction);
}

function copyToClipboard(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.select();
        document.execCommand('copy');
        Utils.showToast('Đã copy vào clipboard!', 'success');
    }
}

function togglePasswordVisibility(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        const button = element.parentElement.querySelector('button i.fa-eye, button i.fa-eye-slash');
        if (element.type === 'password') {
            element.type = 'text';
            if (button) {
                button.className = 'fas fa-eye-slash';
            }
        } else {
            element.type = 'password';
            if (button) {
                button.className = 'fas fa-eye';
            }
        }
    }
}

// Initialize app when DOM is loaded
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new MyCoinApp();
});

// New wallet utility functions
function exportWalletBackup() {
    if (!app.currentWallet) {
        Utils.showToast('Không có ví để xuất backup!', 'error');
        return;
    }

    const backupData = {
        address: app.currentWallet.address,
        private_key: app.currentWallet.private_key,
        public_key: app.currentWallet.public_key,
        passphrase: app.currentWallet.passphrase,
        created_at: new Date().toISOString(),
        backup_version: '1.0'
    };

    const dataStr = JSON.stringify(backupData, null, 2);
    const dataBlob = new Blob([dataStr], {type: 'application/json'});
    
    const link = document.createElement('a');
    link.href = URL.createObjectURL(dataBlob);
    link.download = `mycoin-wallet-backup-${app.currentWallet.address.substring(0, 8)}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    Utils.showToast('Backup wallet đã được xuất thành công!', 'success');
}

function clearWalletData() {
    if (!app.currentWallet) {
        Utils.showToast('Không có ví để xóa!', 'error');
        return;
    }

    const confirmed = confirm(
        'Bạn có chắc chắn muốn xóa thông tin ví khỏi trình duyệt?\n\n' +
        'Lưu ý: Điều này chỉ xóa thông tin ví khỏi trình duyệt này, ' +
        'không ảnh hưởng đến blockchain. Bạn có thể import lại ví bằng private key hoặc passphrase.'
    );

    if (confirmed) {
        app.currentWallet = null;
        localStorage.removeItem('mycoin_wallet');
        
        // Hide wallet info
        document.getElementById('wallet-info').style.display = 'none';
        
        Utils.showToast('Đã xóa thông tin ví khỏi trình duyệt!', 'success');
        
        // Refresh wallet tab
        app.loadWalletData();
    }
}