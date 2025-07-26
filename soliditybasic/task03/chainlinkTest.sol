// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";

contract PriceConsumerV3 {
    AggregatorV3Interface internal priceFeed;

    // Sepolia 测试网 ETH/USD 预言机地址
    constructor() {
        priceFeed = AggregatorV3Interface(0x694AA1769357215DE4FAC081bf1f309aDC325306);
    }
    

    function getLatestPrice() public view returns (int) {
        (, int price, , , ) = priceFeed.latestRoundData();
        // 价格精度为 8 位小数，需转换（例如 2000 USD = 2000_00000000）
        return price / 1e8; // 返回整数形式的 USD 价格（如 2000）
    }
}