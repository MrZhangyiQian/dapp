// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IERC777Recipient.sol";
import "./interfaces/IERC1820Registry.sol";
import "./interfaces/IERC777.sol"; 

contract TokenReceiver is IERC777Recipient {
    IERC1820Registry public erc1820;
    bytes32 public constant TOKENS_RECIPIENT_INTERFACE_HASH = keccak256("ERC777TokensRecipient");
    
    event TokensReceived(
        address operator,
        address from,
        address to,
        uint256 amount,
        bytes userData,
        bytes operatorData
    );
    
    constructor(address erc1820RegistryAddress) {
        erc1820 = IERC1820Registry(erc1820RegistryAddress);
        // 注册接收接口
        erc1820.setInterfaceImplementer(address(this), TOKENS_RECIPIENT_INTERFACE_HASH, address(this));
    }
    
    function tokensReceived(
        address operator,
        address from,
        address to,
        uint256 amount,
        bytes calldata userData,
        bytes calldata operatorData
    ) external override {
        // 这里实现转账后的自动处理逻辑
        // 例如：将代币再转给另一个地址，或者执行其他合约调用
        
        emit TokensReceived(operator, from, to, amount, userData, operatorData);
        
        // 示例：自动将10%的代币转回发送者
        uint256 returnAmount = amount / 10;
        if (returnAmount > 0) {
            IERC777(msg.sender).send(from, returnAmount, "");
        }
    }
}