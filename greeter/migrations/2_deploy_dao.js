var Greeter = artifacts.require("Greeter");

module.exports = function (deployer) {
	// deployment steps
	deployer.deploy(Greeter, "0x1234567890");
};