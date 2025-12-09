// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MyERC20 {
    // 状态变量
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;
    
    // 映射：存储账户余额
    mapping(address => uint256) public balanceOf;
    // 映射：存储授权信息，allowance[授权人][被授权人] = 授权额度
    mapping(address => mapping(address => uint256)) public allowance;
    
    // 事件：转账事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    // 事件：授权事件
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // 存储部署者地址，用于 mint 权限控制
    address public deployer;
    
    // 合并后的构造函数：初始化代币信息和部署者地址
    constructor(string memory _name, string memory _symbol, uint8 _decimals) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        deployer = msg.sender; // 同时初始化部署者地址
    }
    
    // mint 函数：合约所有者增发代币
    function mint(address to, uint256 amount) public {
        require(msg.sender == deployer, "Only deployer can mint");
        totalSupply += amount;
        balanceOf[to] += amount;
        emit Transfer(address(0), to, amount);
    }
    
    // transfer 函数：转账
    function transfer(address to, uint256 amount) public returns (bool) {
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        emit Transfer(msg.sender, to, amount);
        return true;
    }
    
    // approve 函数：授权
    function approve(address spender, uint256 amount) public returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }
    
    // transferFrom 函数：代扣转账
    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        require(balanceOf[from] >= amount, "Insufficient balance");
        require(allowance[from][msg.sender] >= amount, "Insufficient allowance");
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        allowance[from][msg.sender] -= amount;
        emit Transfer(from, to, amount);
        return true;
    }
}
//待测试（晚上看视频测试功能）
// 任务描述
作业 1：ERC20 代币
任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
合约包含以下标准 ERC20 功能：
balanceOf：查询账户余额。
transfer：转账。
approve 和 transferFrom：授权和代扣转账。
使用 event 记录转账和授权操作。
提供 mint 函数，允许合约所有者增发代币。
提示：
使用 mapping 存储账户余额和授权信息。
使用 event 定义 Transfer 和 Approval 事件。
部署到sepolia 测试网，导入到自己的钱包