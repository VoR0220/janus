const greeter = artifacts.require("Greeter");


contract("Greeter", async accounts => {

	it("has been deployed", async () => {
		const greeterDeployed = await greeter.deployed();
		assert(greeterDeployed, "contract has been deployed");
		console.log("Address is: ", greeterDeployed.address);
	})
	let greeterInstance = new web3.eth.Contract(greeter.abi, greeter.address);
	it("should handle events properly", async () => {
		await greeterInstance.methods.greet("this is a log").send({ from: accounts[0] });
		await greeterInstance.methods.greet("this is discarded").send({ from: accounts[0] });
		let expectedResponse1 = { greeter: accounts[0], message: "this is a log" };
		let expectedResponse2 = { greeter: accounts[1], message: "this is discarded" };
		let pastEvents = await greeterInstance.getPastEvents('ListGreeting');
		assert.equal(pastEvents.length, 2, "expected 2 past events, got " + pastEvents.length)
		assert.equal(pastEvents[0].returnValues, expectedResponse1, "incorrect response received in first return value");
		assert.equal(pastEvents[1].returnValues, expectedResponse2, "incorrect response received in second return value");
		let filteredEvents = await greeterInstance.getPastEvents('ListGreeting', { filter: { greeter: accounts[1] } });
		assert.equal(filteredEvents[0].returnValues, expectedResponse2, "expected to get a filtered event of discarded logs but got wrong value");
	})
})
