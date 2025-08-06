// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./Auction.sol";

//创建拍卖合约的工厂合约
contract AuctionFactory {
    // 存储所有创建的拍卖合约地址
    address[] public auctions;
    // 用于记录某个地址是否为拍卖合约
    mapping(address => bool) public isAuction;
    
    // 当拍卖合约被创建时触发，记录拍卖合约地址、卖家地址和 NFT 合约地址
    event AuctionCreated(address indexed auction, address indexed seller, address indexed nftContract);
    
    // 创建一个新的拍卖合约
    function createAuction(
       address _nftContract,
       uint256 _tokenId,
       uint256 _duration,
       address _priceFeed
    ) external returns (address) {
        // 部署一个新的 Auction 逻辑合约。
       Auction auctionLogic = new Auction();
       auctionLogic.initialize(msg.sender, _nftContract, _tokenId, _duration, _priceFeed);
        // 触发事件，记录拍卖合约信息。
        emit AuctionCreated(address(auctionLogic), msg.sender, _nftContract);
        return address(auctionLogic);
}
    
    function getAllAuctions() external view returns (address[] memory) {
        return auctions;
    }
}