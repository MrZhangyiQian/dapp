// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract mergeNums{

    function merge(uint[] calldata arr1,uint[] calldata arr2) public pure returns (uint256[] memory merged){
        if (arr1.length == 0) return arr2;
        if (arr2.length == 0) return arr1;

        merged = new uint256[](arr1.length + arr2.length);

        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;
        while(i < arr1.length && j < arr2.length){
            if (arr1[i] < arr2[j]{
                merged[k++] = arr1[i++];
            })else {
                merged[k++] = arr2[j++];
            }
        }

        while(i < arr1.length){
            merged[k++] = arr1[i++];
        }
        
        while (j < arr2.length) {
            merged[k++] = arr2[j++];
        }

        return merged;
    }
}