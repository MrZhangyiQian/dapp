// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract AdvancedRomanConverter {
    // 存储罗马字符到数值的映射
    mapping(bytes1 => uint16) private _romanValues;
    
    // 事件记录转换
    event RomanConverted(string roman, uint256 value);
    
    constructor() {
        _romanValues['I'] = 1;
        _romanValues['V'] = 5;
        _romanValues['X'] = 10;
        _romanValues['L'] = 50;
        _romanValues['C'] = 100;
        _romanValues['D'] = 500;
        _romanValues['M'] = 1000;
    }
    
    /**
     * @dev 将罗马数字字符串转换为整数
     */
    function romanToInt(string memory roman) public returns (uint256) {
        bytes memory romanBytes = bytes(roman);
        uint256 length = romanBytes.length;
        
        // 空字符串检查
        require(length > 0, "Empty string");
        
        uint256 result = 0;
        uint256 i = 0;
        
        while (i < length) {
            bytes1 currentChar = romanBytes[i];
            require(_romanValues[currentChar] > 0, "Invalid Roman character");
            
            uint256 currentValue = _romanValues[currentChar];
            uint256 nextValue = 0;
            
            // 如果有下一个字符
            if (i < length - 1) {
                bytes1 nextChar = romanBytes[i + 1];
                if (_romanValues[nextChar] > 0) {
                    nextValue = _romanValues[nextChar];
                }
            }
            
            // 处理减法表示法
            if (nextValue > currentValue) {
                // 验证减法组合是否合法
                result += (nextValue - currentValue);
                i += 2; // 跳过一个字符
            } else {
                result += currentValue;
                i += 1;
            }
        }
        
        emit RomanConverted(roman, result);
        return result;
    }
    
    /**
     * @dev 批量转换罗马数字
     */
    function batchConvert(string[] memory romans) public returns (uint256[] memory) {
        uint256[] memory results = new uint256[](romans.length);
        
        for (uint256 i = 0; i < romans.length; i++) {
            results[i] = romanToInt(romans[i]);
        }
        
        return results;
    }
}