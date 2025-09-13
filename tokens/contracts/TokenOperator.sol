// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IERC777.sol";

contract TokenOperator {
    IERC777 public token;
    
    constructor(address tokenAddress) {
        token = IERC777(tokenAddress);
    }
    
    function operatorSend(
        address sender,
        address recipient,
        uint256 amount,
        bytes calldata data,
        bytes calldata operatorData
    ) external {
        // 在实际应用中，这里应该添加权限控制
        token.operatorSend(sender, recipient, amount, data, operatorData);
    }
    
    function tokensReceived(
        address operator,
        address from,
        address to,
        uint256 amount,
        bytes calldata userData,
        bytes calldata operatorData
    ) external pure {
        // 操作员合约也可以实现接收钩子
        // 这里可以添加自定义逻辑
    }
}