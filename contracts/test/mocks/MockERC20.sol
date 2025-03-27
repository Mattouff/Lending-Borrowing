// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MockERC20
 * @author DeFi Lending Platform
 * @notice Mock ERC20 token for testing purposes
 * @dev Extends ERC20 with mint and burn capabilities for the owner
 */
contract MockERC20 is ERC20, Ownable {
    uint8 private _decimals;

    /**
     * @dev Constructor to create a new MockERC20 token
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param initialSupply The initial supply of tokens to mint
     * @param decimalsValue The number of decimals for the token
     * @param owner The address that will own this contract and be able to mint/burn
     */
    constructor(string memory name, string memory symbol, uint256 initialSupply, uint8 decimalsValue, address owner)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {
        _decimals = decimalsValue;

        // Mint initial supply
        if (initialSupply > 0) {
            _mint(owner, initialSupply);
        }

        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @dev Returns the number of decimals used for the token
     * @return The number of decimals
     */
    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    /**
     * @dev Mints tokens to an address
     * @param to The address to mint tokens to
     * @param amount The amount of tokens to mint
     */
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    /**
     * @dev Burns tokens from an address
     * @param from The address to burn tokens from
     * @param amount The amount of tokens to burn
     */
    function burn(address from, uint256 amount) external onlyOwner {
        _burn(from, amount);
    }

    /**
     * @dev Allows anyone to mint tokens to themselves (for testing)
     * @param amount The amount of tokens to mint
     */
    function faucet(uint256 amount) external {
        _mint(msg.sender, amount);
    }

    /**
     * @dev Creates a low level failure for testing purposes
     */
    function forceFailure() external pure {
        revert("MockERC20: forced failure");
    }

    /**
     * @dev Simulates a transfer that depletes all tokens for testing
     * @param to The recipient address
     * @param amount The amount to transfer
     */
    function transferWithLowBalance(address to, uint256 amount) external {
        _burn(msg.sender, balanceOf(msg.sender));
        emit Transfer(msg.sender, to, amount);
    }
}
