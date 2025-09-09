// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./Auction.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AuctionFactory is Ownable {
    address[] public auctions;
    mapping(address => bool) public isAuction;

    event AuctionCreated(
        address indexed auction,
        address indexed seller,
        address indexed nftContract
    );

    constructor() Ownable(msg.sender) {}

    function createAuction(
        address _nftContract,
        uint256 _tokenId,
        uint256 _duration,
        address _priceFeed
    ) external onlyOwner returns (address) {
       
        Auction auctionLogic = new Auction();
        auctionLogic.initialize(
            msg.sender,
            _nftContract,
            _tokenId,
            _duration,
            _priceFeed
        );
        address auctionAddress = address(auctionLogic);
        auctions.push(auctionAddress);
        isAuction[auctionAddress] = true;
        emit AuctionCreated(address(auctionLogic), msg.sender, _nftContract);
        return address(auctionLogic);
    }

    function getAllAuctions() external view returns (address[] memory) {
        return auctions;
    }
}
