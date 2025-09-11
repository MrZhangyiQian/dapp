// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "./interfaces/IAuction.sol";
import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
// 改为这个
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract Auction is
    IAuction,
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    // 记录出价信息
    struct Bid {
        // 出价人地址
        address bidder;
        // 金额
        uint256 amount;
        // 时间
        uint256 timestamp;
        // 美元价值
        uint256 usdValue;
    }
    // 拍卖详细信息
    struct AuctionItem {
        // 卖家
        address seller;
        // NFT 地址
        address nftContract;
        // ID
        uint256 tokenId;
        // 起止时间
        uint256 startTime;
        // 结束时间
        uint256 endTime;
        // 最高出价者
        address highestBidder;
        // 最高出价
        uint256 highestBid;
        // 最高出价（USD）
        uint256 highestBidUsd;
        // 是否结束
        bool ended;
    }

    uint256 public constant AUCTION_DURATION = 3600; // 1 hour
    // 拍卖信息
    AuctionItem public auction;
    // 记录出价人信息，Bid[]动态数据考虑多次出价情况
    mapping(address => Bid[]) public bids;
    // 通过调用这个预言机，合约可以获取当前 ETH 对 USD 的价格。
    AggregatorV3Interface internal ethPriceFeed;
    mapping(address => AggregatorV3Interface) internal erc20PriceFeeds;

    // 拍卖创建时触发
    event AuctionCreated(
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId,
        uint256 startTime,
        uint256 endTime
    );
    // 用户出价时触发
    event BidPlaced(address indexed bidder, uint256 amount, uint256 usdValue);
    // 拍卖结束时触发
    event AuctionEnded(address indexed winner, uint256 amount, uint256 tokenId);

    // 拍卖无人出价时触发
    event AuctionCancelled();
    // 合约升级后触发
    event Upgraded(address indexed implementation);

    // 重写IAcution 方法，initializer函数修改器确保此函数只能被调用一次，避免重复初始化

    function initialize(
        address _seller,
        address _nftContract,
        uint256 _tokenId,
        uint256 _duration,
        address _priceFeed
    ) public payable override initializer {
        // 初始化合约拥有者为 _seller
        __Ownable_init(_seller);
        // 初始化 UUPS 升级相关逻辑（如权限控制）
        __UUPSUpgradeable_init();

        ethPriceFeed = AggregatorV3Interface(_priceFeed);

        auction = AuctionItem({
            seller: _seller,
            nftContract: _nftContract,
            tokenId: _tokenId,
            startTime: block.timestamp,
            endTime: block.timestamp + _duration,
            highestBidder: address(0),
            highestBid: 0,
            highestBidUsd: 0,
            ended: false
        });

        emit AuctionCreated(
            _seller,
            _nftContract,
            _tokenId,
            block.timestamp,
            block.timestamp + _duration
        );
    }

    // 用户出价
    function placeBid() external payable {
        // 验证拍卖是否未结束
        require(block.timestamp <= auction.endTime, "Auction ended");
        // 验证出价必须高于当前最高价
        require(msg.value > auction.highestBid, "Bid too low");
        // 不是卖家
        require(msg.sender != auction.seller, "Seller cannot bid");

        // 获取ETH/USD价格
        (, int256 price, , , ) = ethPriceFeed.latestRoundData();
        //  ETH 对应的美元价值
        uint256 usdValue = (msg.value * uint256(price)) / 1e18;

        // 如果有前一个最高出价者，返还其 ETH
        if (auction.highestBidder != address(0)) {
            payable(auction.highestBidder).transfer(auction.highestBid);
        }
        // 更新最高出价
        auction.highestBidder = msg.sender;
        auction.highestBid = msg.value;
        auction.highestBidUsd = usdValue;

        // 保存出价记录
        bids[msg.sender].push(
            Bid({
                bidder: msg.sender,
                amount: msg.value,
                timestamp: block.timestamp,
                usdValue: usdValue
            })
        );

        emit BidPlaced(msg.sender, msg.value, usdValue);
    }

    // 拍卖结束后调用
    function endAuction() external {
        // 确保拍卖已结束
        require(block.timestamp >= auction.endTime, "Auction not ended");
        // 确保拍卖未结束
        require(!auction.ended, "Auction already ended");
        // 确保拍卖者或合约拥有者调用
        require(
            msg.sender == auction.seller || msg.sender == owner(),
            "Not authorized"
        );

        auction.ended = true;
        // 有出价人
        if (auction.highestBidder != address(0)) {
            // 转移NFT给出价最高者
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );

            // 转移资金给卖家
            payable(auction.seller).transfer(auction.highestBid);

            emit AuctionEnded(
                auction.highestBidder,
                auction.highestBid,
                auction.tokenId
            );
        } else {
            // 如果没有出价，将NFT返还给卖家
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
            // 发送NFT的取消通知
            emit AuctionCancelled();
        }
    }

    function getEthUsdPrice() public view returns (uint256) {
        (, int256 price, , , ) = ethPriceFeed.latestRoundData();
        return uint256(price);
    }

    function getBidUsdValue(uint256 ethAmount) public view returns (uint256) {
        (, int256 price, , , ) = ethPriceFeed.latestRoundData();
        return (ethAmount * uint256(price)) / 1e18;
    }

    // UUPS升级授权
    // 确保只有合约所有者可以授权升级。
    function _authorizeUpgrade(address) internal virtual override onlyOwner {}

    // 为了安全禁止在拍卖结束前转移NFT
    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public pure returns (bytes4) {
        return this.onERC721Received.selector;
    }
    // 一个 NFT 到拍卖合约中，设定默认拍卖时间
    // _ NFT 合约地址
    // _ NFT ID
    function bindNFT(address _nftContract, uint256 _tokenId) external payable {
        require(auction.nftContract == address(0), "NFT already bound");
        require(
            IERC721(_nftContract).ownerOf(_tokenId) == address(this),
            "NFT not in contract"
        );

        auction.nftContract = _nftContract;
        auction.tokenId = _tokenId;
        auction.seller = msg.sender;
        auction.startTime = block.timestamp;
        auction.endTime = block.timestamp + AUCTION_DURATION; // 可以从常量中获取
    }
}
