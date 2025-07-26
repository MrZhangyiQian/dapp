
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract RomaNumToInt {
    function revertWithString(string memory roman) public pure returns (uint256) {
        bytes memory romanBytes = bytes(roman);
        uint length = romanBytes.length;
        require(length > 0, "Empty string");
        uint256 result = 0;
        for (uint i = 0; i < length; i++){
            uint currentValue = charToValue(romanBytes[i]);
            if (i < length - 1){
                uint nextValue = charToValue(romanBytes[i+1]);
                if (currentValue < nextValue) {
                    result += nextValue - currentValue;
                    i++;
                    continue;
             }
         }
         result += currentValue;
        }
        return result;
    }

    function charToValue(bytes1 c) private pure returns (uint256){
         if (c == 'I') return 1;
        if (c == 'V') return 5;
        if (c == 'X') return 10;
        if (c == 'L') return 50;
        if (c == 'C') return 100;
        if (c == 'D') return 500;
        if (c == 'M') return 1000;
        // 如果遇到无效字符则抛出错误
        revert("Invalid Roman character");
    }
       
}
