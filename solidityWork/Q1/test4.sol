// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanNumerals {
    uint[] private values = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
    string[] private symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];
    

    function toRoman(uint num) public view returns (string memory) {
        require(num > 0 && num <= 3999, "Number must be between 1 and 3999");
        
        string memory result = "";
        
        for (uint i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                result = string(abi.encodePacked(result, symbols[i]));
                num -= values[i];
            }
        }
        
        return result;
    }
}