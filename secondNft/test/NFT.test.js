const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MyNFT", function () {
  let nft;
  let owner, user;

  beforeEach(async function () {
    [owner, user] = await ethers.getSigners();
    
    const NFT = await ethers.getContractFactory("MyNFT");
    nft = await NFT.deploy();
    await nft.waitForDeployment();
  });

  it("Should have correct name and symbol", async function () {
    expect(await nft.name()).to.equal("MyNFT");
    expect(await nft.symbol()).to.equal("MNFT");
  });

  it("Should mint NFTs correctly", async function () {
    // 铸造给owner
    const tokenId0 = await nft.mint(owner.address);
    expect(await nft.ownerOf(0)).to.equal(owner.address);
    
    // 铸造给user
    const tokenId1 = await nft.safeMint(user.address);
    expect(await nft.ownerOf(1)).to.equal(user.address);
  });

  it("Should track token counter", async function () {
    // 初始计数为0
    // 第一次铸造
    await nft.mint(owner.address);
    // 第二次铸造
    await nft.mint(owner.address);
    
    // 总供应量应该是2
    expect(await nft.totalSupply()).to.equal(2);
    // 余额应该是2
    expect(await nft.balanceOf(owner.address)).to.equal(2);
  });
});