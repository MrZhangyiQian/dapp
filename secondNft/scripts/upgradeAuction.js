// scripts/upgradeAuction.js
const { ethers, upgrades } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();

    // 获取之前部署的代理地址（可以从 deployments/${network}.json 中读取）
    const proxyAddress = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"; // 替换为实际的 Proxy 地址

    console.log("Upgrading Auction Proxy at:", proxyAddress);

    // 获取新版本合约的 Factory
    const AuctionV2 = await ethers.getContractFactory("AuctionV2");

    // 执行升级
    const upgraded = await upgrades.upgradeProxy(proxyAddress, AuctionV2);
    await upgraded.waitForDeployment();
    console.log("Auction 升级完成！新实现地址:", await upgrades.erc1967.getImplementationAddress(proxyAddress));
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });