// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IDeposit.sol";

contract Oracle is IOracle, OwnableUpgradeable, UUPSUpgradeable {
    IDeposit public deposit;
    IPayment public payment;

    function initialize() public initializer {
        __Ownable_init(msg.sender);
    }

    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    function verifyServiceActive(address user) external view returns (bool) {
        if (address(deposit) == address(0)) return false;
        (bool ok, bytes memory data) = address(deposit).staticcall(
            abi.encodeWithSignature("getLockExpiry(address)", user)
        );
        if (!ok) return false;
        (uint256 lockExpiry) = abi.decode(data, (uint256));
        return lockExpiry > 0 && block.timestamp < lockExpiry;
    }

    function monthlySettlement(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata amounts
    ) external onlyOwner {
        require(users.length == operatorIds.length, "Length mismatch");
        require(users.length == amounts.length, "Length mismatch");
        for (uint256 i = 0; i < users.length; i++) {
            uint256 amount = amounts[i];
            if (amount > 0) {
                payment.createBill(users[i], operatorIds[i], amount);
                emit UsageDataSubmitted(users[i], operatorIds[i], amount);
            }
        }
        if (address(payment) != address(0)) {
            for (uint256 i = 0; i < users.length; i++) {
                IPayment.Bill[] memory unpaidBills = IPayment(payment).getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    IPayment(payment).applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }
    }

    function submitUsage(address user, uint256, uint256, uint256) external view {
        require(user != address(0), "Invalid user");
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
