// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract BinarySearch {
    function binarySearch(uint256[] calldata data, uint256  target) public pure returns(uint256){
        uint256 left = 0;
        uint256 len = data.length;
        uint256 right = len-1;
        if (len == 0) return type(uint256).max;
        while (left <= right ){
            uint256 mid = left + ((right - left) / 2);
            if(data[mid] == target){
                return mid;
            }else if (target < data[mid]) {
                right = mid -1;
            }else {
                left = mid + 1;
            }
        }
        return 0;
    }
}