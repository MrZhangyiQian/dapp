// scripts/deploy.js
const { ethers, upgrades } = require("hardhat"); // 引入 upgrades
const fs = require("fs");

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying contracts with account:", deployer.address);

  // 1. 部署 NFT 合约
  const MyNFT = await ethers.getContractFactory("MyNFT");
  const nft = await MyNFT.deploy();
  await nft.waitForDeployment();
  console.log("MyNFT deployed to:", nft.target);

  // 2. 铸造 NFT
  const tokenId = 0;
  const mintTx = await nft.mint(deployer.address);
  await mintTx.wait();
  console.log(`NFT ID ${tokenId} 已铸造给 ${deployer.address}`);
  console.log("NFT当前所有者:", await nft.ownerOf(tokenId));

  // 3. 部署工厂合约
  const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
  const factory = await AuctionFactory.deploy();
  await factory.waitForDeployment();
  console.log("AuctionFactory deployed to:", factory.target);

  // 4. 使用工厂创建拍卖代理合约（使用升级插件）
  const SEPOLIA_FEED = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";
  const duration = 3600;

  console.log("通过工厂创建拍卖...");
  const auctionProxy = await upgrades.deployProxy(
    await ethers.getContractFactory("Auction"),
    [
      deployer.address,     // _seller
      nft.target,           // _nftContract
      tokenId,              // _tokenId
      duration,             // _duration
      SEPOLIA_FEED          // _priceFeed
    ],
    { initializer: "initialize" }
  );
  await auctionProxy.waitForDeployment();
  const auctionProxyAddress = await auctionProxy.getAddress();
  console.log("Auction Proxy deployed to:", auctionProxyAddress);

  // 5. 将 NFT 转移到代理合约中
  console.log("将NFT转移到拍卖合约...");
  const transferTx = await nft["safeTransferFrom(address,address,uint256)"](
    deployer.address,
    auctionProxyAddress,
    tokenId
  );
  await transferTx.wait();
  console.log("NFT已转移至拍卖合约");

  // 6. 获取拍卖合约实例（不再需要手动初始化）
  const auction = await ethers.getContractAt("Auction", auctionProxyAddress);
  
  // 7. 查询拍卖信息
  const auctionInfo = await auction.auction();

  console.log("拍卖卖家:", auctionInfo.seller);
  console.log("NFT合约地址:", auctionInfo.nftContract);
  console.log("Token ID:", auctionInfo.tokenId.toString());
  console.log("开始时间:", new Date(Number(auctionInfo.startTime) * 1000));
  console.log("结束时间:", new Date(Number(auctionInfo.endTime) * 1000));
  console.log("NFT当前所有者:", await nft.ownerOf(tokenId));

  // 8. 保存部署信息到文件
  const deploymentInfo = {
    network: network.name,
    nftAddress: nft.target,
    factoryAddress: factory.target,
    proxyAddress: auctionProxyAddress,
    timestamp: Date.now()
  };

  if (!fs.existsSync("deployments")) {
    fs.mkdirSync("deployments");
  }

  fs.writeFileSync(
    `deployments/${network.name}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("部署完成！");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });