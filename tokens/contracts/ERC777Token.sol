// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IERC777.sol";
import "./interfaces/IERC777Recipient.sol";
import "./interfaces/IERC777Sender.sol";
import "./interfaces/IERC1820Registry.sol";

contract ERC777Token is IERC777 {
    string public constant name = "ERC777Token";
    string public constant symbol = "ERC777";
    uint8 public constant decimals = 18;
    uint256 public totalSupply;
    
    IERC1820Registry public erc1820;
    bytes32 public constant TOKENS_SENDER_INTERFACE_HASH = keccak256("ERC777TokensSender");
    bytes32 public constant TOKENS_RECIPIENT_INTERFACE_HASH = keccak256("ERC777TokensRecipient");
    
    mapping(address => uint256) private balances;
    mapping(address => mapping(address => bool)) private operators;
    
    // 修复：添加下划线前缀避免命名冲突
    address[] private _defaultOperators;
    mapping(address => bool) private _isDefaultOperator;
    mapping(address => mapping(address => bool)) private _revokedDefaultOperators;
    
    constructor(uint256 initialSupply, address[] memory initialDefaultOperators, address erc1820RegistryAddress) {
        erc1820 = IERC1820Registry(erc1820RegistryAddress);
        totalSupply = initialSupply;
        balances[msg.sender] = initialSupply;
        _defaultOperators = initialDefaultOperators;
        
        for (uint256 i = 0; i < initialDefaultOperators.length; i++) {
            _isDefaultOperator[initialDefaultOperators[i]] = true;
        }
        
        // 注册代币接口
        erc1820.setInterfaceImplementer(address(this), keccak256("ERC777Token"), address(this));
        erc1820.setInterfaceImplementer(address(this), keccak256("ERC20Token"), address(this));
    }
    
    function balanceOf(address holder) public view override returns (uint256) {
        return balances[holder];
    }
    
    function send(address recipient, uint256 amount, bytes calldata data) external override {
        _send(msg.sender, msg.sender, recipient, amount, data, "", true);
    }
    
    function operatorSend(
        address sender,
        address recipient,
        uint256 amount,
        bytes calldata data,
        bytes calldata operatorData
    ) external override {
        require(isOperatorFor(msg.sender, sender), "ERC777: caller is not operator");
        _send(msg.sender, sender, recipient, amount, data, operatorData, true);
    }
    
    function burn(uint256 amount, bytes calldata data) external override {
        _burn(msg.sender, msg.sender, amount, data, "");
    }
    
    function operatorBurn(
        address account,
        uint256 amount,
        bytes calldata data,
        bytes calldata operatorData
    ) external override {
        require(isOperatorFor(msg.sender, account), "ERC777: caller is not operator");
        _burn(msg.sender, account, amount, data, operatorData);
    }
    
    function authorizeOperator(address operator) external override {
        operators[msg.sender][operator] = true;
        emit AuthorizedOperator(operator, msg.sender);
    }
    
    function revokeOperator(address operator) external override {
        operators[msg.sender][operator] = false;
        emit RevokedOperator(operator, msg.sender);
    }
    
    // 修复：返回私有变量
    function defaultOperators() public view override returns (address[] memory) {
        return _defaultOperators;
    }
    
    function isOperatorFor(address operator, address tokenHolder) public view override returns (bool) {
        return (
            tokenHolder == operator ||
            operators[tokenHolder][operator] ||
            (_isDefaultOperator[operator] && !_revokedDefaultOperators[tokenHolder][operator])
        );
    }
    
    function _send(
        address operator,
        address from,
        address to,
        uint256 amount,
        bytes memory userData,
        bytes memory operatorData,
        bool requireReception
    ) private {
        require(from != address(0), "ERC777: send from zero address");
        require(to != address(0), "ERC777: send to zero address");
        require(amount <= balances[from], "ERC777: insufficient balance");
        
        balances[from] -= amount;
        balances[to] += amount;
        
        emit Sent(operator, from, to, amount, userData, operatorData);
        
        if (requireReception) {
            address implementer = erc1820.getInterfaceImplementer(to, TOKENS_RECIPIENT_INTERFACE_HASH);
            if (implementer != address(0)) {
                IERC777Recipient(implementer).tokensReceived(operator, from, to, amount, userData, operatorData);
            }
        }
    }
    
    function _burn(
        address operator,
        address from,
        uint256 amount,
        bytes memory data,
        bytes memory operatorData
    ) private {
        require(from != address(0), "ERC777: burn from zero address");
        require(amount <= balances[from], "ERC777: insufficient balance");
        
        balances[from] -= amount;
        totalSupply -= amount;
        
        emit Burned(operator, from, amount, data, operatorData);
    }
    
    // ERC20兼容方法
    function transfer(address recipient, uint256 amount) external returns (bool) {
        _send(msg.sender, msg.sender, recipient, amount, "", "", true);
        return true;
    }
}