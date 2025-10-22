// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanNumerals {
    // 将映射设为内部，以便在内部函数中使用
    mapping (bytes1 => uint256) internal romanIntMap;
    
    constructor() {
        romanIntMap["I"] = 1;
        romanIntMap["V"] = 5;
        romanIntMap["X"] = 10;
        romanIntMap["L"] = 50;
        romanIntMap["C"] = 100;
        romanIntMap["D"] = 500;
        romanIntMap["M"] = 1000;
    }
    
    function romanToInt(string memory romanString) public view returns (uint256) {
        bytes memory romanByte = bytes(romanString);
        uint256 len = romanByte.length;
        
        // 空字符串检查
        if (len == 0) {
            return 0;
        }
        
        uint256 result = romanIntMap[bytes1(romanByte[len - 1])]; // 初始化为最右边字符的值
        
        for (uint256 i = len - 1; i > 0; i--) {
            uint256 current = romanIntMap[bytes1(romanByte[i - 1])];
            uint256 next = romanIntMap[bytes1(romanByte[i])];
            
            if (current < next) {
                result -= current;
            } else {
                result += current;
            }
        }
        
        return result;
        

    }
}