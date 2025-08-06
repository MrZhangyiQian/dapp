const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("AuctionFactory", () => {
  let factory, nft;
  let deployer, seller;

  // Sepolia ETH/USD 预言机地址
  const ETH_USD_PRICE_FEED = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
  const AUCTION_DURATION = 3600; // 1小时

  before(async () => {
    [deployer, seller] = await ethers.getSigners();

    // 部署NFT合约
    const NFT = await ethers.getContractFactory("MyNFT");
    nft = await NFT.deploy();
    await nft.waitForDeployment();

    // 部署工厂合约
    const Factory = await ethers.getContractFactory("AuctionFactory");
    factory = await Factory.deploy();
    await factory.waitForDeployment();
  });

  it("Should create new auctions", async () => {
    // 铸造NFT并批准
    await nft.safeMint(seller.address, 0);
    await nft.connect(seller).approve(factory.target, 0);

    // 创建拍卖
    const tx = await factory.connect(seller).createAuction(
      nft.target,
      0,
      AUCTION_DURATION,
      ETH_USD_PRICE_FEED
    );

    const receipt = await tx.wait();
    const auctionAddress = await factory.auctions(0);

    // 验证事件
    expect(auctionAddress).to.be.properAddress;
    expect(await factory.isAuction(auctionAddress)).to.be.true;

    // 验证创建的拍卖数量
    const allAuctions = await factory.getAllAuctions();
    expect(allAuctions.length).to.equal(1);
    expect(allAuctions[0]).to.equal(auctionAddress);
  });
}); 