// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EnhancedBeggingContract {
    // 合约所有者
    address payable public owner;
    
    // 记录每个地址的捐赠总额
    mapping(address => uint256) public donations;
    
    // 存储所有捐赠者地址（去重）
    address[] private donors;
    
    // 存储地址是否已经计入捐赠者列表
    mapping(address => bool) private isDonor;
    
    // 捐赠事件
    event DonationReceived(address indexed donor, uint256 amount);
    event FundsWithdrawn(address indexed owner, uint256 amount);
    
    // 确保只有合约所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only contract owner can call this function");
        _;
    }
    
    // 构造函数：设置合约所有者
    constructor() {
        owner = payable(msg.sender);
    }
    
    // 接收以太币的捐赠函数
    function donate() public payable {
        // 验证捐赠金额大于0
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        // 如果是第一次捐赠，加入捐赠者列表
        if (donations[msg.sender] == 0) {
            donors.push(msg.sender);
            isDonor[msg.sender] = true;
        }
        
        // 更新捐赠记录
        donations[msg.sender] += msg.value;
        
        // 触发捐赠事件
        emit DonationReceived(msg.sender, msg.value);
    }
    
    // 获取捐赠者总数（精确统计）
    function getTotalDonors() public view returns (uint256) {
        return donors.length;
    }
    
    // 获取捐赠者地址（分页）
    function getDonorsByPage(uint256 page, uint256 pageSize) public view returns (address[] memory) {
        require(page > 0, "Page number must be greater than 0");
        uint256 start = (page - 1) * pageSize;
        
        // 防止越界
        if (start >= donors.length) {
            return new address[](0);
        }
        
        uint256 end = start + pageSize;
        if (end > donors.length) {
            end = donors.length;
        }
        
        // 创建结果数组
        address[] memory result = new address[](end - start);
        for (uint256 i = start; i < end; i++) {
            result[i - start] = donors[i];
        }
        
        return result;
    }
    
    // 获取所有捐赠者地址（仅用于测试）
    function getAllDonors() public view returns (address[] memory) {
        return donors;
    }
    
    // 允许所有者提取所有资金
    function withdraw() public onlyOwner {
        // 获取合约当前余额
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds available to withdraw");
        
        // 转移资金给所有者
        owner.transfer(balance);
        
        // 触发提款事件
        emit FundsWithdrawn(owner, balance);
    }
    
    // 查询特定地址的捐赠金额
    function getDonation(address donor) public view returns (uint256) {
        return donations[donor];
    }
    
    // 获取合约总余额
    function getContractBalance() public view returns (uint256) {
        return address(this).balance;
    }
    
    // 销毁合约（只有所有者可以调用）
    function destroy() public onlyOwner {
        // 转移剩余资金给所有者
        withdraw();
        // 销毁合约
        selfdestruct(owner);
    }
    
    // 接收直接转账的回调函数
    receive() external payable {
        donate();
    }
}