// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleERC20{
     // 代币基本信息
     string public constant name = "Simple ERC20 Token";
     string public constant symbol = "SIM";
     uint8 public constant decimals = 18;

     // 状态变量
     uint256 public totalSupply;
     address public owner;

     // 余额存储
     mapping(address => uint256) private _balances;
     // 授权存储
     mapping(address => mapping(address => uint256)) private _allowances;

     // 事件定义
     event Transfer(address indexed from,address indexed to,uint256 value);
     event Approval(address indexed owner,address indexed spender,uint256 value);

     // 构造函数
     constructor(uint256 _initialSupply){
          owner = msg.sender;
          _mint(msg.sender,_initialSupply * (10 ** uint256(decimals)));
     }

     function _mint(address to, uint256 amount) private {
          require(to != address(0),"TO should not be zero");
          totalSupply += amount;
          _balances[to] += amount;
          emit Transfer(address(0), to, amount);
     }

     // 查询余额
     function balanceOf(address account) public view returns (uint256){
          return _balances[account];
     }

     // 转账功能
     function transfer(address to, uint256 amount) public returns(bool){
          require(to != address(0),"TO should not be zero");
          require(_balances[msg.sender] >= amount,"Insufficient balance!");
          _balances[msg.sender] -= amount;
          _balances[to] += amount;
          emit Transfer(msg.sender, to, amount);
          return true;
     }

     // 授权功能
     function approve(address spender,uint256 amount) public returns(bool){
          _allowances[msg.sender][spender] = amount;
          emit Approval(msg.sender, spender, amount);
          return true;
     }

     // 查询授权额度
     function allowance(address ownerAddress, address spender) public view returns (uint256){
          return _allowances[ownerAddress][spender];
     }

     // 授权转账功能
     function transferFrom(address from, address to,uint256 amount) public returns(bool){
          require(to != address(0),"ERC20: transfer to the zero address");
          require(_balances[from] >= amount,"ERC20: insufficient balance");
          require(_allowances[from][msg.sender] >= amount,"ERC20: transfer amount exceeds allowance");
          _balances[from] -= amount;
          _balances[to] += amount;
          _allowances[from][msg.sender] -= amount;

          emit Transfer(from, to, amount);
          return true;
     }

     // 增发代币
     function mint(address to, uint256 amount) public{
          require(msg.sender == owner,"Only owner can mint");
          _mint(to,amount);
     }

     // 转移所有权
     function transferOwnership(address newOwner) public {
          require(msg.sender == owner,"Only the original owner can change the owner");
          owner = newOwner;
     }
}