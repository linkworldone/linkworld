// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract UserRegistry is Ownable, ReentrancyGuard {
    

    struct UserInfo {
        address wallet;
        string email;
        uint256 registerTime;
        bool isRegistered;
        string ogDataRoot;      
    }

    mapping(address => UserInfo) public users;
    
    uint256 public totalUsers;

    event UserRegistered(
        address indexed wallet,
        string email,
        uint256 timestamp,
        string ogDataRoot
    );

    event UserUpdated(
        address indexed wallet,
        string newEmail,
        string newOgDataRoot
    );

    
    error AlreadyRegistered();
    error InvalidAddress();
    error EmptyEmail();


    constructor() Ownable(msg.sender) {}

 
    function registerUser(
        string calldata _email,
        string calldata _ogDataRoot   
    ) external nonReentrant {
        if (users[msg.sender].isRegistered) revert AlreadyRegistered();
        if (bytes(_email).length == 0) revert EmptyEmail();
        if (msg.sender == address(0)) revert InvalidAddress();

        users[msg.sender] = UserInfo({
            wallet: msg.sender,
            email: _email,
            registerTime: block.timestamp,
            isRegistered: true,
            ogDataRoot: _ogDataRoot
        });

        totalUsers++;

        emit UserRegistered(msg.sender, _email, block.timestamp, _ogDataRoot);
    }

  
    function updateUserInfo(
        string calldata _newEmail,
        string calldata _newOgDataRoot
    ) external {
        if (!users[msg.sender].isRegistered) revert("User not registered");

        users[msg.sender].email = _newEmail;
        users[msg.sender].ogDataRoot = _newOgDataRoot;

        emit UserUpdated(msg.sender, _newEmail, _newOgDataRoot);
    }

  
    function getUserInfo(address _wallet) external view returns (UserInfo memory) {
        return users[_wallet];
    }

    function isUserRegistered(address _wallet) external view returns (bool) {
        return users[_wallet].isRegistered;
    }
}
