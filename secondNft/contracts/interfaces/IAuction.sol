// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 拍卖合约接口
interface IAuction {
    function initialize(
        // 卖家
        address _seller,
        // NFT合约地址
        address _nftContract,
        // 代币ID
        uint256 _tokenId,
        // 拍卖持续时间
        uint256 _duration,
        // 价格预言机地址
        address _priceFeed
    ) external payable;
    // 参与竞价
    function placeBid() external payable;
    //  结束拍卖
    function endAuction() external;
    // 获取ETH对USD的当前价格
    function getEthUsdPrice() external view returns (uint256);
    // 将ETH金额转换为USD价值
    function getBidUsdValue(uint256 ethAmount) external view returns (uint256);

    // ✅ 添加 getter 函数声明
    function auction()
        external
        view
        returns (
            address seller,
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