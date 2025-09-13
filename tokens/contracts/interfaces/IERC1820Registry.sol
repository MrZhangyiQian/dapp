// IERC1820Registry.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC1820Registry {
    function setInterfaceImplementer(address account, bytes32 interfaceHash, address implementer) external;
    function getInterfaceImplementer(address account, bytes32 interfaceHash) external view returns (address);
}