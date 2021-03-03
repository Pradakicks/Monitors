const {
	targetMonitor
} = require('./sites/targetEfficent');
const {
	newEggMonitor
} = require('./sites/newEgg')
const {
	gameStopMonitor
} = require('./sites/gameStop');
const delay = require('delay');
require('newrelic');
var skuBank = []

const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip: true,
});

function SKUADD(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.channel.type === 'dm' && message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				//    fetch('')
				console.log(site)
				console.log(SKU)
				console.log(content)
				if (SKU.length > 1 && site.length > 1) {
					if (site.toUpperCase() == 'TARGET') {
						skuBank.push({
							sku: SKU,
							site: 'TARGET',
							stop: false
						})
						let monitor = new targetMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: SKU,
							site: 'NEWEGG',
							stop: false
						})
						let monitor = new newEggMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: SKU,
							site: 'GAMESTOP',
							stop: false
						})
						let monitor = new gameStopMonitor(SKU.toString())
						monitor.task()
					}
				}
				console.log(skuBank)
			}
		});
	} catch (error) {
		console.log(error);
	}

}

function findCommand(clients, triggerText, replyText) {
	clients.on('message', message => {
		if (message.channel.type === 'dm' && message.content.toLowerCase() === triggerText.toLowerCase()) {
			message.author.send(replyText);
		}
	});

}

function deleteSku(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.channel.type === 'dm' && message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				console.log(site)
				console.log(SKU)
				console.log(content)
				let index = skuBank.findIndex(e => e.sku == SKU)
				// console.log(index)
				// console.log(skuBank[index])
				skuBank[index].stop = true;
				// console.log(skuBank)
				(async ()=>{
					await delay(100000)
					skuBank.splice(index, 1)
				})()
				// console.log(skuBank)
				//    fetch('')

			}
		});
	} catch (error) {
		console.log(error);
	}
}


module.exports = {
	SKUADD,
	findCommand,
	deleteSku,
	skuBank
}