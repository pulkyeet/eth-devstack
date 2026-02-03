// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Token2178 is ERC20, ERC20Burnable, ERC20Pausable, Ownable {
    uint256 public constant MAX_SUPPLY = 1_000_000 * 10 **18;

    // maps an address to unix timestamp. till then, it would be locked (no transfers)
    mapping(address => uint256) public lockedUntil;

    event TokensLocked(address indexed account, uint256 until);

    // setup the token
    constructor() ERC20("TK2178", "T21") Ownable(msg.sender) {
        _mint(msg.sender, 100_000_000 * 10 ** 18);
    }

    //give tokens to an account
    function mint(address to, uint256 amount) public onlyOwner {
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        _mint(to, amount);
    }

    //pauses token transfers (inherits)
    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    // address tokens locked
    function lockTokens(address account, uint256 until) public onlyOwner {
        require(until > block.timestamp, "Must be future timestamp");

        lockedUntil[account] = until;

        emit TokensLocked(account, until);
    }

    function _update(address from, address to, uint256 value) internal override(ERC20, ERC20Pausable) {
        if (from!=address(0)) {
            require(block.timestamp >= lockedUntil[from], "Tokens are locked.");
        }

        super._update(from, to, value);
    }
}