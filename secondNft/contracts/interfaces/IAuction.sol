// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IAuction {
    function initialize(
        address _seller,
        address _nftContract,
        uint256 _tokenId,
        uint256 _duration,
        address _priceFeed
    ) external payable;
    
    function placeBid(uint256 _amount) external payable;
    function endAuction() external;
    function getEthUsdPrice() external view returns (uint256);
    function getBidUsdValue(uint256 ethAmount) external view returns (uint256);

    // ✅ 添加 getter 函数声明
    function auction()
        external
        view
        returns (
            address seller,
            address originalOwner,
            address nftContract,
            uint256 tokenId,
            uint256 startTime,
            uint256 endTime,
            address highestBidder,
            uint256 highestBid,
            uint256 highestBidUsd,
            bool ended
        ); 
}
