# Sample Hardhat Project

This project demonstrates a basic Hardhat use case. It comes with a sample contract, a test for that contract, and a Hardhat Ignition module that deploys that contract.

// 部署合约
Pull 代码 终端命令
1、安装依赖
npm install --save-dev @openzeppelin/contracts-upgradeable @openzeppelin/hardhat-upgrades @nomicfoundation/hardhat-toolbox @chainlink/contracts dotenv
 npm install @chainlink/contracts@0.8 --save-dev
2、增加配置文件.env
SPEOLIA_URL=
PRIVATE_KEY=
ETHERSCAN_API_KEY=


  test result:

npx hardhat node
npx hardhat run scripts/deploy.js --network localhost
npx hardhat run scripts/upgradeAuction.js --network localhost
npx hardhat test test/Auction.test.js

 1、 PS D: npx hardhat run scripts/deploy.js --network localhost
[dotenv@17.0.1] injecting env (3) from .env – [tip] encrypt with dotenvx: https://dotenvx.com
拍卖卖家: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
NFT合约地址: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Token ID: 0
开始时间: 2025-07-04T07:25:41.000Z
结束时间: 2025-07-04T08:25:41.000Z
NFT当前所有者: 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9
部署完成！



2、PS D: npx hardhat run scripts/upgradeAuction.js --network localhost
[dotenv@17.0.1] injecting env (3) from .env – [tip] encrypt with dotenvx: https://dotenvx.com
[dotenv@17.0.1] injecting env (0) from .env – [tip] encrypt with dotenvx: https://dotenvx.com
Upgrading Auction Proxy at: 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9
Auction 升级完成！新实现地址: 0x0165878A594ca255338adfa4d48449f69242Eb8F

 3、 PS D: npx hardhat test test/Auction.test.js
[dotenv@17.0.1] injecting env (3) from .env – [tip] encrypt with dotenvx: https://dotenvx.com
Warning: Unused function parameter. Remove or comment out the variable name to silence this warning.
  --> contracts/AuctionFactory.sol:15:8:
   |
15 |        uint256 _tokenId,
   |        ^^^^^^^^^^^^^^^^


Warning: Unused function parameter. Remove or comment out the variable name to silence this warning.
  --> contracts/AuctionFactory.sol:16:8:
   |
16 |        uint256 _duration,
   |        ^^^^^^^^^^^^^^^^^


Warning: Unused function parameter. Remove or comment out the variable name to silence this warning.
  --> contracts/AuctionFactory.sol:17:8:
   |
17 |        address _priceFeed
   |        ^^^^^^^^^^^^^^^^^^


Compiled 3 Solidity files successfully (evm target: paris).
duplicate definition - Upgraded(address)
duplicate definition - Upgraded(address)


  Auction Contract
    Initialization
      ✔ should initialize auction correctly
    Bidding Logic
      ✔ should allow a bid and update highest bidder
      ✔ should revert if seller tries to bid
      ✔ should revert if bid is not higher than current
      ✔ should revert if auction has ended
    Ending Auction
      ✔ should end auction and transfer NFT to winner
      ✔ should return NFT if no bids
      ✔ should revert if called before auction ends
      ✔ should revert if called by unauthorized user
    Price Calculation
      ✔ should get correct ETH/USD price
      ✔ should calculate USD value of ETH correctly
    UUPS Upgrade Support
      ✔ should upgrade contract successfully
Proxy Owner: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Bidder1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
AuctionV2 Implementation: 0xfbC22278A96299D91d41C453234d97b4F5Eb9B2d
      ✔ should only allow owner to upgrade


  13 passing (703ms)

··············································································································
|  Solidity and Network Configuration                                                                        │
·························|·················|···············|·················|································
|  Solidity: 0.8.28      ·  Optim: false   ·  Runs: 200    ·  viaIR: false   ·     Block: 30,000,000 gas     │
·························|·················|···············|·················|································
|  Methods                                                                                                   │
·························|·················|···············|·················|················|···············
|  Contracts / Methods   ·  Min            ·  Max          ·  Avg            ·  # calls       ·  usd (avg)   │
·························|·················|···············|·················|················|···············
|  Auction               ·                                                                                   │
·························|·················|···············|·················|················|···············
|      endAuction        ·         97,216  ·      107,880  ·        102,548  ·             2  ·           -  │
·························|·················|···············|·················|················|···············
|      placeBid          ·              -  ·            -  ·        221,132  ·             3  ·           -  │
·························|·················|···············|·················|················|···············
|  AuctionV2             ·                                                                                   │
·························|·················|···············|·················|················|···············
|      upgradeToAndCall  ·              -  ·            -  ·         38,511  ·             1  ·           -  │
·························|·················|···············|·················|················|···············
|  MyNFT                 ·                                                                                   │
·························|·················|···············|·················|················|···············
|      approve           ·         49,009  ·       49,021  ·         49,020  ·            26  ·           -  │
·························|·················|···············|·················|················|···············
|      mint              ·              -  ·            -  ·         91,307  ·            26  ·           -  │
·························|·················|···············|·················|················|···············
|      safeTransferFrom  ·         64,604  ·       64,616  ·         64,615  ·            13  ·           -  │
·························|·················|···············|·················|················|···············
|  Deployments                             ·                                 ·  % of limit    ·              │
·························|·················|···············|·················|················|···············
|  Auction               ·              -  ·            -  ·      2,529,288  ·         8.4 %  ·           -  │
·························|·················|···············|·················|················|···············
|  AuctionV2             ·              -  ·            -  ·      2,553,827  ·         8.5 %  ·           -  │
·························|·················|···············|·················|················|···············
|  MockV3Aggregator      ·              -  ·            -  ·        408,698  ·         1.4 %  ·           -  │
·························|·················|···············|·················|················|···············
|  MyNFT                 ·              -  ·            -  ·      1,788,695  ·           6 %  ·           -  │
·························|·················|···············|·················|················|···············
|  Key                                                                                                       │
··············································································································
|  ◯  Execution gas for this method does not include intrinsic gas overhead                                  │
··············································································································
|  △  Cost was non-zero but below the precision setting for the currency display (see options)               │
··············································································································
|  Toolchain:  hardhat                                                                                       │
··············································································································
