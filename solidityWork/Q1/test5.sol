// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BinarySearch {
    function binarySearch(uint256[] memory nums, uint256 target) public pure returns (uint256) {
        uint256 len = nums.length;
        uint256 head = 0;
        uint256 tail = len - 1;

        while (tail >= head) {
            uint256 mid = head + (tail - head) / 2; // 防止溢出
            if (nums[mid] > target) {
                tail = mid - 1;
            } else if (nums[mid] < target) {
                head = mid + 1;
            } else {
                return mid; // 返回目标值的索引
            }
        }

        // 如果目标值不存在，返回一个特殊值
        return type(uint256).max;
    }
}