// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IPayment {
    struct Bill {
        uint256 id;
        address user;
        uint256 operatorId;
        uint256 amount;       // 运营商费用
        uint256 platformFee;  // 平台手续费
        uint256 createdAt;
        bool isPaid;
    }

    event BillCreated(uint256 indexed billId, address indexed user, uint256 amount, uint256 platformFee);
    event BillPaid(uint256 indexed billId, address indexed user, uint256 totalAmount);

    function createBill(address user, uint256 operatorId, uint256 amount) external;
    function payBill(uint256 billId) external payable;
    function getUserBills(address user) external view returns (Bill[] memory);
    function getUnpaidBills(address user) external view returns (Bill[] memory);
    function autoSettle(address[] calldata users, uint256[] calldata operatorIds, uint256[] calldata amounts) external;
}
