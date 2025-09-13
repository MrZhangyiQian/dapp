// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ERC1820Registry {
    mapping(address => mapping(bytes32 => address)) private interfaces;
    
    event InterfaceSet(address indexed account, bytes32 indexed interfaceHash, address implementer);
    
    function setInterfaceImplementer(address account, bytes32 interfaceHash, address implementer) external {
        interfaces[account][interfaceHash] = implementer;
        emit InterfaceSet(account, interfaceHash, implementer);
    }
    
    function getInterfaceImplementer(address account, bytes32 interfaceHash) external view returns (address) {
        return interfaces[account][interfaceHash];
    }
}