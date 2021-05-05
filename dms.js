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
				let kw = []
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				if(content.includes('[') && site.toUpperCase() !== "TARGETNEW"){
					pricerange = content.split('[')[1].split(']')[0]
				//	SKU = content.split('')

				}
				if(site.toUpperCase() == "TARGETNEW"){
					let kwArray = content.split('[')[1].split(']')[0]
					console.log(kwArray)
					let eachItem = kwArray.split(',')
					eachItem.map(e =>{
						kw.push(e)
					})
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
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Target",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/target`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: SKU,
							site: 'NEWEGG',
							stop: false,
                          	name: ""
						})

						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "New Egg",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
									skuName: await getSku(SKU, await getProxies())
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/newEgg`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: SKU,
							site: 'GAMESTOP',
							stop: false,
                          name: ""
						})
						let monitor = new gameStopMonitor(SKU.toString())
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
						console.log(pricerange)
							let currentBody = {
								  	site: "Walmart",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1])
							}
						if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/walmart`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						

						// let monitor = new walmartMonitor(SKU.toString(), pricerange)
						// monitor.task()
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: SKU,
							site: 'BESTBUY',
							stop: false,
                          name: ""
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Best Buy",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
						if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
							console.log(err)
						})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/bestBuy`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
	
					} else if (site.toUpperCase() == 'BIGLOTS') {
						skuBank.push({
							sku: SKU,
							site: 'BIGLOTS',
							stop: false,
                          name: ""
						})
						console.log(pricerange)
							let currentBody = {
								  	site: "Big Lots",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1])
							}
						if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/bigLots`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						

						// let monitor = new walmartMonitor(SKU.toString(), pricerange)
						// monitor.task()
					} else if (site.toUpperCase() == 'TARGETNEW') {
						skuBank.push({
							sku: SKU,
							site: 'TARGETNEW',
							stop: false,
                          name: ""
						})
						console.log(kw)
							let currentBody = {
									  endpoint : SKU,
									  keywords : kw
							}
							console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/targetNew`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						

						// let monitor = new walmartMonitor(SKU.toString(), pricerange)
						// monitor.task()
					} else if (site.toUpperCase() == 'ACADEMY') {
						skuBank.push({
							sku: SKU,
							site: 'ACADEMY',
							stop: false,
                          name: ""
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Academy",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
								console.log(err)
							})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/academy`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
					} else if (site.toUpperCase() == 'AMD') {
						skuBank.push({
							sku: SKU,
							site: 'AMD',
							stop: false,
                          name: ""
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Amd",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
								console.log(err)
							})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/amd`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
					}
				console.log(skuBank)
				fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
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
					fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				});
				// console.log(skuBank)
				(async () => {
					const fs = require('fs')
					await delay(10000)
					skuBank.splice(index, 1)
					await delay(2000)
					await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
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
					SKU = g[i]
					if(g[i].includes('[')){
					pricerange = g[i].split('[')[1].split(']')[0]
					SKU = g[i].split(' ')[0]
					}
				
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
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Target",
									sku: g[i],
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/target`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: g[i],
							site: 'NEWEGG',
							stop: false,
                          name: ""
						})

						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "NewEgg",
									sku: g[i],
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
									skuName: await getSku(g[i], await getProxies())
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
						await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
							try {
							rp.post({
							url : `http://localhost:7243/newEgg`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new newEggMonitor(g[i].toString())
						// monitor.task()
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
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Amd",
									sku: g[i],
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
								console.log(err)
							})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/amd`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
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
						console.log(pricerange)
							let currentBody = {
								  	site: "Walmart",
									sku: SKU,
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1])
							}
					if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
					console.log(err)
				})
							console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/walmart`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new walmartMonitor(g[i].toString())
						// monitor.task()
						await delay(30000)
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: g[i],
							site: 'BESTBUY',
							stop: false,
                          name: ""
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Best Buy",
									sku: g[i],
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
						if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
							console.log(err)
						})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/bestBuy`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
	
					} else if (site.toUpperCase() == 'ACADEMY') {
						skuBank.push({
							sku: g[i],
							site: 'ACADEMY',
							stop: false,
                          name: ""
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
								  	site: "Academy",
									sku: g[i],
									priceRangeMin: parseInt(pricerange.split(',')[0]),
									priceRangeMax: parseInt(pricerange.split(',')[1]),
							}
							if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
							console.log("No Max Price Range Detected")
							currentBody.priceRangeMax = 100000

						} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
							console.log("No Min Price Range Detected")
							currentBody.priceRangeMin = 1

						}
							await fs.writeFile('./GoMonitor/GoMonitors.json', JSON.stringify(skuBank), err => {
								console.log(err)
							})
						console.log(currentBody)
							try {
							rp.post({
							url : `http://localhost:7243/academy`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						})
							} catch (error) {
								console.log(error)
							}
						// let monitor = new targetMonitor(SKU.toString())
						// monitor.task()
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
async function getProxies () {
								try {
						// read contents of the file
						let proxyList = []
						const data = await fs.readFile('proxies.txt', 'utf-8');

						const lines = data.split(/\r?\n/);

						// print all lines
						lines.forEach((line) => {
							const lineSplit = line.split(':');
							const item1 = {
								ip : lineSplit[0],
								port: lineSplit[1],
								userAuth: lineSplit[2],
								userPass: lineSplit[3],
							};
							proxyList.push(item1);

							// console.log(line);
							// console.log(item1);
							// console.log('\n\n\n\n\n');
						});
						return proxyList;
					} catch (err) {
						console.error(err);
						fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
							console.log(err)
						}))
					}
}
async function getSku (skuName, proxyList) {
        try {
                   console.log(skuName)
                    let proxy1 = proxyList[Math.floor(Math.random() * proxyList.length)]
                    console.log(proxy1)
                    let fetchProductPage = await rp.get({
                    url : `https://www.newegg.com/prada/p/${skuName}`,
                    headers : {
                    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
                    "accept-language": "en-US,en;q=0.9",
                    "cache-control": "no-cache",
                    "pragma": "no-cache",
                    "sec-fetch-dest": "document",
                    "sec-fetch-mode": "navigate",
                    "sec-fetch-site": "none",
                    "sec-fetch-user": "?1",
                    "upgrade-insecure-requests": "1"
                    },
                    proxy : `http://${proxy1.userAuth}:${proxy1.userPass}@${proxy1.ip}:${proxy1.port}`
                })
                console.log(fetchProductPage?.statusCode)
                sku = fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[0] + '-'
                console.log(sku)
                sku = sku + fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[1] + '-'
                console.log(sku)
                sku = sku + fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[2] 
                console.log(sku)
                return sku
                } catch (error) {
                    console.log(error)
                    // skuBank[this.index].name = 'Restart'
                    // skuBank[this.index]["error"] = error.message
                    // skuBank[this.index].stop = true
                  //  console.log(skuBank[this.index])
                    await getSku()
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