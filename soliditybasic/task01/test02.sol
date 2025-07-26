// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract RevertString {
    function revertWithString(string memory _input) public pure returns (string memory) {
       bytes memory inputBytes = bytes(_input);
        uint256 length = inputBytes.length;
        if (length == 0) {
            return _input;
        }
        bytes memory reversedBytes = new bytes(length);
        // 使用两个指针，i从0开始，j从末尾开始
        for (uint i = 0; i < length; i++) {
            // j = length - 1 - i
            reversedBytes[i] = inputBytes[length - 1 - i];
        }
        return string(reversedBytes);
    }
}