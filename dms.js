const {
	targetMonitor
} = require('./sites/target');
const {
	newEggMonitor
} = require('./sites/newEgg')
const {
	gameStopMonitor
} = require('./sites/gameStop');
const { bestBuyMonitor } = require('./sites/bestBuy')
const {MessageAttachment } = require('discord.js');
const { amdMonitor } = require('./sites/amd')
const { amdSiteMonitor } = require('./sites/amdSite')
const { walmartMonitor } = require('./sites/walmart')
const delay = require('delay');
const fs = require('fs').promises
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
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				let pricerange = ''
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				if(content.includes('[')){
					pricerange = content.split('[')[1].split(']')[0]

				}
				//    fetch('')
				console.log(site)
				console.log(SKU)
				console.log(content)
				console.log(pricerange)
				
				if (SKU.length > 1 && site.length > 1) {
					let isContinue = true
					skuBank.map(e => {
						if(e.sku == SKU){
							console.log('Duplicate Found')
							message.channel.send(`${SKU} is already present in monitor`)
							isContinue = false
						}
					})
					if(isContinue){
					if (site.toUpperCase() == 'TARGET') {
						skuBank.push({
							sku: SKU,
							site: 'TARGET',
							stop: false,
                          name: ""
						})
						let monitor = new targetMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: SKU,
							site: 'NEWEGG',
							stop: false,
                          name: ""
						})
						let monitor = new newEggMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: SKU,
							site: 'GAMESTOP',
							stop: false,
                          name: ""
						})
						let monitor = new gameStopMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMD') {
						skuBank.push({
							sku: SKU,
							site: 'AMD',
							stop: false,
                          name: ""
						})
						let monitor = new amdMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMDSITE') {
						skuBank.push({
							sku: SKU,
							site: 'AMDSITE',
							stop: false,
                          name: ""
						})
						let monitor = new amdSiteMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'WALMART') {
						skuBank.push({
							sku: SKU,
							site: 'WALMART',
							stop: false,
                          name: ""
						})
						let monitor = new walmartMonitor(SKU.toString(), pricerange)
						monitor.task()
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: SKU,
							site: 'BESTBUY',
							stop: false,
                          name: ""
						})
						let monitor = new bestBuyMonitor(SKU.toString())
						monitor.task()
					}
				console.log(skuBank)
				message.channel.send(`${SKU} Added to ${site}`)
					}
				
			}
		}});
	} catch (error) {
		console.log(error);
	}

}

function findCommand(clients, triggerText, replyText) {
	clients.on('message', message => {
		if (message.content.toLowerCase() === triggerText.toLowerCase()) {
			message.author.send(replyText);
		}
	});

}

function deleteSku(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				console.log(site)
				console.log(`SKU - ${SKU}`)
				console.log(content)
				let index = skuBank.findIndex(e => e.sku == SKU)
				// console.log(index)
				// console.log(skuBank[index])
				skuBank[index].stop = true;
				// console.log(skuBank)
				(async () => {
					await delay(10000)
					skuBank.splice(index, 1)
				})()
				
				// console.log(skuBank)
				function replaceWithTheCapitalLetter(values){
				return values.charAt(0).toUpperCase() + values.slice(1);
				}
				message.channel.send(`${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`)
				return;
				//    fetch('')

			}
		});
	} catch (error) {
		console.log(error);
	}
}

function checkBank (clients, triggerText, replyText){
	
		clients.on('message', async (message) => {
			try {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				let string = (skuBank)
				let skuString = ''
				if(skuBank.length != 0){
					for (let i = 0; i < skuBank.length; i++){
					let status = "Terminated"
						if(!skuBank[i].stop){
						status = "Running"
					}
					await fs.appendFile(`monitorBank-${message.author.username}.txt`, `Site : ${skuBank[i].site} | ${skuBank[i].sku} | ${status} | ${skuBank[i].name} \n` , err => {
						if(err) message.content.send('Error While Creating Text Document')
								else console.log("File Sent")
				})
				}
				let attachment = new MessageAttachment(`monitorBank-${message.author.username}.txt`);
				message.channel.send(attachment)
				message.author.send('Attachment Successfully Fetched and Sent')
				await delay(2500)
				await fs.unlink(`monitorBank-${message.author.username}.txt`, err =>{
					if(err) console.log('Error doing the unthinkable')
				})
				} else {
					message.channel.send('Monitor Bank is empty')
				}
		
				
			}
			} catch (error) {
		console.log(error);
		message.channel.send('Error checking Bank')
	}
		});
	
}

function massAdd (clients, triggerText, replyText){
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				
				let string = message.content
				const content = message.content;
				const site = content.split(' ')[1].split('|')[0]
				console.log(site.toUpperCase())
				
			

			//	const SKU = content.split(' ')[2];
			//	console.log(site)
				let g  = string.split('\n')
			//	console.log(g)
			for(let i = 0; i < g.length; i++){		
				if(!g[i].toUpperCase().includes('!MASSADD')){
					let isContinue = true
					let SKU
					let pricerange = ''
					if(g[i].includes('[')){
					pricerange = g[i].split('[')[1].split(']')[0]
					SKU = g[i].split(' ')[0]
					}
					SKU = g[i]
					console.log(g[i])
					console.log(site.toUpperCase())
						skuBank.map(e => {
						if(e.sku == SKU) {
							console.log('Duplicate Found')
							message.channel.send(`${SKU} is already present in monitor`)
							isContinue = false
						}
					})
					if(isContinue){
							if (site.toUpperCase() == 'TARGET') {
						skuBank.push({
							sku: g[i],
							site: 'TARGET',
							stop: false,
                          name: ""
						})
						let monitor = new targetMonitor(g[i].toString())
						monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: g[i],
							site: 'NEWEGG',
							stop: false,
                          name: ""
						})
						let monitor = new newEggMonitor(g[i].toString())
						monitor.task()
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: g[i],
							site: 'GAMESTOP',
							stop: false,
                          name: ""
						})
						let monitor = new gameStopMonitor(g[i].toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMD') {
						skuBank.push({
							sku: g[i],
							site: 'AMD',
							stop: false,
                          name: ""
						})
						let monitor = new amdMonitor(g[i].toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMDSITE') {
						skuBank.push({
							sku: g[i],
							site: 'AMDSITE',
							stop: false,
                          name: ""
						})
						let monitor = new amdSiteMonitor(g[i].toString())
						monitor.task()
					} else if (site.toUpperCase() == 'WALMART') {
						skuBank.push({
							sku: g[i],
							site: 'WALMART',
							stop: false,
                          name: ""
						})
						let monitor = new walmartMonitor(g[i].toString())
						monitor.task()
						await delay(30000)
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: g[i],
							site: 'BESTBUY',
							stop: false,
                          name: ""
						})
						let monitor = new bestBuyMonitor(g[i].toString())
						monitor.task()
						
					}
					}
				

				}
				await delay(30000)
			}
			message.channel.send("SKUS Added")
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
	checkBank,
	skuBank,
	massAdd
}