pragma solidity ^0.8.9;

contract Bank
{
    int balance;

    constructor() public
    {
        balance = 1;
    }

    function getBalance() view public returns(int)
    {
        return balance;
    }
    function withdraw(int ammount) public
    {
        balance = balance - ammount;
    }
    function deposit(int ammount) public
    {
        balance = balance + ammount;
    }
}
