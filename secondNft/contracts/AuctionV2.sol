// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./Auction.sol";

contract AuctionV2 is Auction {
    // 添加新功能以区分版本
    function version() public pure returns (string memory) {
        return "V2";
    }
    
    // 保持权限检查
    function _authorizeUpgrade(address) internal override onlyOwner {}
}