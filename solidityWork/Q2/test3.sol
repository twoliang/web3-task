//任务描述
### ✅ 作业3：编写一个讨饭合约
任务目标
1. 使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
2. 记录每个捐赠者的地址和捐赠金额。
3. 允许合约所有者提取所有捐赠的资金。

任务步骤
1. 编写合约
  - 创建一个名为 BeggingContract 的合约。
  - 合约应包含以下功能：
  - 一个 mapping 来记录每个捐赠者的捐赠金额。
  - 一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
  - 一个 withdraw 函数，允许合约所有者提取所有资金。
  - 一个 getDonation 函数，允许查询某个地址的捐赠金额。
  - 使用 payable 修饰符和 address.transfer 实现支付和提款。
2. 部署合约
  - 在 Remix IDE 中编译合约。
  - 部署合约到 Goerli 或 Sepolia 测试网。
3. 测试合约
  - 使用 MetaMask 向合约发送以太币，测试 donate 功能。
  - 调用 withdraw 函数，测试合约所有者是否可以提取资金。
  - 调用 getDonation 函数，查询某个地址的捐赠金额。

任务要求
1. 合约代码：
  - 使用 mapping 记录捐赠者的地址和金额。
  - 使用 payable 修饰符实现 donate 和 withdraw 函数。
  - 使用 onlyOwner 修饰符限制 withdraw 函数只能由合约所有者调用。
2. 测试网部署：
  - 合约必须部署到 Goerli 或 Sepolia 测试网。
3. 功能测试：
  - 确保 donate、withdraw 和 getDonation 函数正常工作。

提交内容
1. 合约代码：提交 Solidity 合约文件（如 BeggingContract.sol）。
2. 合约地址：提交部署到测试网的合约地址。
3. 测试截图：提交在 Remix 或 Etherscan 上测试合约的截图。

//------------------------------------------------------------------
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title BeggingContract
 * @dev 允许用户向合约捐赠以太币，记录捐赠信息，并支持所有者提取资金
 */
contract BeggingContract {
    address public owner;
    mapping(address => uint256) public donations; // 记录每个捐赠者的捐赠金额

    constructor() {
        owner = msg.sender; // 合约部署者自动成为所有者
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    /**
     * @dev 用户捐赠以太币，记录到捐赠者地址
     * 注意：调用时需附带以太币（如通过 MetaMask 发送）
     */
    function donate() public payable {
        donations[msg.sender] += msg.value; // 累加捐赠金额
    }

    /**
     * @dev 合约所有者提取全部资金
     * 使用 address.transfer 安全转账（自动处理 2300 gas 限制）
     */
    function withdraw() public onlyOwner {
        uint256 amount = address(this).balance;
        owner.transfer(amount); // 提取合约全部余额
    }

    /**
     * @dev 查询指定地址的捐赠金额
     * @param donor 捐赠者地址
     * @return 捐赠金额（单位：wei）
     */
    function getDonation(address donor) public view returns (uint256) {
        return donations[donor];
    }
}
