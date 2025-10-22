// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract StringReverser {
    function reverseString(string memory s) public pure returns (string memory) {
        bytes memory byteS = bytes(s);
        uint256 stringLength = byteS.length;
        bytes memory reversedS = new bytes(stringLength);

        for (uint256 i = 0; i < stringLength; i++) {
            reversedS[i] = byteS[stringLength-1-i];
        }

        return string(reversedS);

    }
}