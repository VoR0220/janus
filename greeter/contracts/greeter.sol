pragma solidity >= 0.6.12;

contract Greeter {
	event ListGreeting(address indexed greeter, string message);
    function greet(string memory message) public {
        emit ListGreeting(msg.sender, message);
    }
}