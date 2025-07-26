// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract SimpleNFT is ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;
    
    constructor() 
        ERC721("SimpleNFT", "SNFT")
        Ownable(msg.sender) 
    {}
    
    // 铸造函数 
    function mintNFT(address recipient, string memory tokenURI) public onlyOwner returns (uint256) {
        uint256 newItemId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        _safeMint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);
        return newItemId;
    }
    
    // 更新元数据（仅所有者或被授权用户可以调用）
    function updateTokenURI(uint256 tokenId, string memory newURI) public {
        address owner = ownerOf(tokenId);
        require(
            msg.sender == owner || 
            isApprovedForAll(owner, msg.sender) || 
            getApproved(tokenId) == msg.sender,
            "Not owner or approved"
        );
        _setTokenURI(tokenId, newURI);
    }
    
    // 锁定元数据
    function lockMetadata(uint256 tokenId) public view {
        address owner = ownerOf(tokenId);
        require(
            msg.sender == owner,
            "Only token owner can lock"
        );
        // 测试部署地址 https://testnets.opensea.io/assets/sepolia/0xe6220d218d7fb528908a92d1ad9929ca5a979452/0/
    }
}