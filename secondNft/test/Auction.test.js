const { ethers, upgrades } = require("hardhat");
const { expect } = require("chai");

describe("Auction Contract", function () {
    let deployer, bidder1, bidder2;
    let nft, auctionProxy;

    const tokenId = 0;
    const duration = 3600; // 1 hour
    const ethUsdFeed = "0x694AA1769357215DE4FAC081bf1f309aDC325306"; // Sepolia ETH/USD 预言机

    beforeEach(async function () {
        [deployer, bidder1, bidder2] = await ethers.getSigners();

        // 部署 NFT
        const MyNFT = await ethers.getContractFactory("MyNFT");
        nft = await MyNFT.deploy();
        await nft.waitForDeployment();

        // 铸造 NFT 给部署者
        const mintTx = await nft.mint(deployer.address);
        await mintTx.wait();

        // 部署 Mock Chainlink Aggregator
        const MockV3Aggregator = await ethers.getContractFactory("MockV3Aggregator");
        mockPriceFeed = await MockV3Aggregator.deploy(8, 2000e8); // $2000 per ETH
        await mockPriceFeed.waitForDeployment();

        // 获取 Auction 合约工厂
        const Auction = await ethers.getContractFactory("Auction");

        // 使用 Hardhat Upgrades 插件部署代理合约，并注入 Mock 预言机
        auctionProxy = await upgrades.deployProxy(Auction, [
            deployer.address,
            nft.target,
            tokenId,
            duration,
            mockPriceFeed.target, // 使用 mock 替代真实预言机
        ], { initializer: "initialize" });
        await auctionProxy.waitForDeployment();

        // 获取代理合约地址并授权
        const auctionAddress = await auctionProxy.getAddress();

        // 授权 NFT 给拍卖合约地址
        const proxyApproveTx = await nft.connect(deployer).approve(auctionAddress, tokenId);
        await proxyApproveTx.wait();

        // 将 NFT 转移到拍卖合约
        await nft.connect(deployer)["safeTransferFrom(address,address,uint256)"](deployer.address, auctionAddress, tokenId);
    });
    describe("Initialization", function () {
        it("should initialize auction correctly", async function () {
            const auctionInfo = await auctionProxy.auction();
            expect(auctionInfo.seller).to.equal(deployer.address);
            expect(auctionInfo.nftContract).to.equal(nft.target);
            expect(auctionInfo.tokenId).to.equal(tokenId);
            expect(await nft.ownerOf(tokenId)).to.equal(await auctionProxy.getAddress());
        });
    });

    describe("Bidding Logic", function () {
        it("should allow a bid and update highest bidder", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await auctionProxy.connect(bidder1).placeBid({ value: bidAmount });

            const auctionInfo = await auctionProxy.auction();
            expect(auctionInfo.highestBidder).to.equal(bidder1.address);
            expect(auctionInfo.highestBid).to.equal(bidAmount);
        });

        it("should revert if seller tries to bid", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await expect(
                auctionProxy.connect(deployer).placeBid({ value: bidAmount })
            ).to.be.revertedWith("Seller cannot bid");
        });

        it("should revert if bid is not higher than current", async function () {
            const bidAmount = ethers.parseEther("0.05");
            await auctionProxy.connect(bidder1).placeBid({ value: bidAmount });

            await expect(
                auctionProxy.connect(bidder2).placeBid({ value: bidAmount })
            ).to.be.revertedWith("Bid too low");
        });

        it("should revert if auction has ended", async function () {
            // 快进时间到拍卖结束后
            await network.provider.send("evm_increaseTime", [duration + 1]);
            await network.provider.send("evm_mine");

            await expect(
                auctionProxy.connect(bidder1).placeBid({ value: ethers.parseEther("0.1") })
            ).to.be.revertedWith("Auction ended");
        });
    });

    describe("Ending Auction", function () {
        it("should end auction and transfer NFT to winner", async function () {
            const bidAmount = ethers.parseEther("0.1");
            await auctionProxy.connect(bidder1).placeBid({ value: bidAmount });

            // 快进时间
            await network.provider.send("evm_increaseTime", [duration + 1]);
            await network.provider.send("evm_mine");

            await auctionProxy.endAuction();

            const owner = await nft.ownerOf(tokenId);
            expect(owner).to.equal(bidder1.address);
        });

        it("should return NFT if no bids", async function () {
            // 不出价直接结束
            await network.provider.send("evm_increaseTime", [duration + 1]);
            await network.provider.send("evm_mine");

            await auctionProxy.endAuction();

            const owner = await nft.ownerOf(tokenId);
            expect(owner).to.equal(deployer.address);
        });

        it("should revert if called before auction ends", async function () {
            await expect(
                auctionProxy.endAuction()
            ).to.be.revertedWith("Auction not ended");
        });

        it("should revert if called by unauthorized user", async function () {
            // 快进时间
            await network.provider.send("evm_increaseTime", [duration + 1]);
            await network.provider.send("evm_mine");

            await expect(
                auctionProxy.connect(bidder1).endAuction()
            ).to.be.revertedWith("Not authorized");
        });
    });

    describe("Price Calculation", function () {
        it("should get correct ETH/USD price", async function () {
            const price = await auctionProxy.getEthUsdPrice();
            expect(price).to.be.gt(0);
        });

        it("should calculate USD value of ETH correctly", async function () {
            const ethAmount = ethers.parseEther("0.5"); // BigInt
            const ethPrice = await auctionProxy.getEthUsdPrice(); // assume this is also a BigInt

            // 使用 ethers.utils.formatUnits 转换后再做乘法
            const ethAmountInEth = Number(ethers.formatEther(ethAmount));
            const ethPriceInUsd = Number(ethers.formatUnits(ethPrice, 8)); // 根据预言机精度调整 decimals
            const expected = ethAmountInEth * ethPriceInUsd;
        });
    });

    describe("UUPS Upgrade Support", function () {
        it("should upgrade contract successfully", async function () {
            const AuctionV2 = await ethers.getContractFactory("AuctionV2");
            // const upgradedProxy = await upgrades.upgradeProxy(await auctionProxy .getAddress(), AuctionV2);  
            // 等价与
            const upgradedProxy = await upgrades.upgradeProxy(await auctionProxy.getAddress(), AuctionV2.connect(deployer));
            await upgradedProxy.waitForDeployment();
            const auctionInfo = await upgradedProxy.auction();
            expect(auctionInfo.seller).to.equal(deployer.address);
        });

        it("should only allow owner to upgrade", async function () {
            const [deployer, bidder1] = await ethers.getSigners();
            const AuctionV2 = await ethers.getContractFactory("AuctionV2");

            // 非 owner 尝试升级 - 应该失败
            await expect(
                upgrades.upgradeProxy(await auctionProxy.getAddress(), AuctionV2.connect(bidder1))
            ).to.be.revertedWith("Ownable: caller is not the owner");
        });

    });
});
