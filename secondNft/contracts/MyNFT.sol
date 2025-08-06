// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

// MyNFT.sol 是一个基于 ERC-721 标准的 非同质化代币（NFT）合约，使用了 OpenZeppelin 提供的安全且经过审计的 ERC721 基础合约。
contract MyNFT is ERC721 {
    // tokenId 的计数器，记录当前已铸造的 NFT 数量，并作为下一个 Token ID 的基础值
    uint256 private _tokenIdCounter;

    constructor() ERC721("MyNFT", "MNFT") {
        _tokenIdCounter = 0;
    }
    // 铸造一个新的 NFT 并将其发送给指定地址
    // mint() 就是用来“生成”一个新 NFT 并送出去的函数。
    // “我要打印一张独一无二的数字藏品，并把它放进某个人的钱包。”
    function mint(address to) public returns (uint256) {
        uint256 tokenId = _tokenIdCounter++;
        _mint(to, tokenId);
        return tokenId;
    }
}
