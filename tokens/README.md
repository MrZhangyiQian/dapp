# Sample Hardhat Project

This project demonstrates a basic Hardhat use case. It comes with a sample contract, a test for that contract, and a script that deploys that contract.

Try running some of the following tasks:

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat run scripts/deploy.js


project/
├── contracts/                 # Solidity 合约文件
│   ├── ERC777Token.sol        # 主代币合约
│   ├── TokenOperator.sol      # 操作员合约
│   ├── TokenReceiver.sol       # 接收合约（带钩子功能）
│   ├── ERC1820Registry.sol     # 接口注册表
│   └── interfaces/             # 接口定义
│       ├── IERC777.sol
│       ├── IERC777Recipient.sol
│       ├── IERC777Sender.sol
│       └── IERC1820Registry.sol
├── deploy/                    # 部署脚本
│   └── 01_deploy.js
├── test/                      # 测试文件
│   └── erc777.test.js
├── scripts/                   # 实用脚本
│   └── deploy.js
└── hardhat.config.js          # Hardhat 配置
```
