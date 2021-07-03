const {MessageAttachment,  MessageEmbed} = require('discord.js');
const delay = require('delay');
const fs = require('fs').promises;
const config = require('./config.json')
const os = require('os')
//  var skuBank = []
let pushEndpoint = "https://monitors-9ad2c-default-rtdb.firebaseio.com/monitor"
let discordIds = "https://monitors-9ad2c-default-rtdb.firebaseio.com/validatedUsers"
const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip: true,
});
const port = 7243
// const secondServer = `http://ec2-3-236-148-149.compute-1.amazonaws.com`
const secondServer = `http://143.244.164.115`
const thirdServer = `http://64.227.28.51`

function SKUADD(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				let pricerange = ''
				let kw = []
				const content = message.content;
				const {group, isValidated} = await checkIfUserValidated(message)
				console.log(group)
				if(isValidated){
				const site = content.split(' ')[1];
				let original = content.split(`${site} `)[1]
				const SKU = content.split(' ')[2];
				if(content.includes('[') && site.toUpperCase() !== "TARGETNEW"){
					pricerange = content.split('[')[1].split(']')[0]
				}
				if(site.toUpperCase() == "TARGETNEW"){
					let kwArray = content.split('[')[1].split(']')[0]
					console.log(kwArray)
					let eachItem = kwArray.split(',')
					eachItem.forEach(e =>{
						kw.push(e)
					})
				}
				console.log(site)
				console.log(SKU)
				console.log(content)
				console.log(original)
				console.log(pricerange)
				
				if (SKU.length > 1 && site.length > 1) {
					let isContinue = true
					let skuBank = await getSkuBank()
					let caseSite = site.toUpperCase()
					let currentCompany = config[group]
					if(skuBank[caseSite]){
						if(skuBank[caseSite][SKU]){
							let skuWebhookArray = skuBank[caseSite][SKU]?.companies
							let isPresent = false
							skuWebhookArray?.forEach(e => {
								console.log(e.webhook, currentCompany[caseSite])
								if(e.webhook == currentCompany[caseSite]) isPresent = true
							})

							if(isPresent){
								message.channel.send(`${SKU} is already present in monitor`)
							} else {
								message.channel.send(`${SKU} is being added to monitor`)
								let arr = []
								skuWebhookArray.forEach(e =>{
									arr.push(e)
								})
								arr.push({
									company : group,
									webhook : currentCompany[caseSite],
									color : currentCompany?.companyColor,
									companyImage : currentCompany?.companyImage
								})
								console.log(arr)
								await updateSku(site, SKU, {companies : arr}, `${pushEndpoint}/${caseSite.toUpperCase()}/${SKU}/.json`)
							}
							console.log('Duplicate Found', isPresent)
							isContinue = false
					}
					}
					
					if(isContinue){
					if (site.toUpperCase() == 'TARGET') {
						await pushSku({
							sku: SKU,
							site: 'TARGET',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
						})
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
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
						
					} else if (site.toUpperCase() == 'NEWEGG') {
						await pushSku({
							sku: SKU,
							site: 'NEWEGG',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'GAMESTOP') {
							await pushSku({
								sku: SKU,
								site: 'GAMESTOP',
								stop: false,
								name: "",
								original : original,
								companies : [
									{
									company : group,
									webhook : currentCompany[caseSite],
									color : currentCompany?.companyColor,
									companyImage : currentCompany?.companyImage
								}]
							})
								let currentBody = {
										  site: "Game Stop",
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
							console.log(currentBody)
							startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'WALMART') {
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
							await pushSku({
							sku: SKU,
							site: 'WALMART',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
							})
							console.log(currentBody)
							startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'BESTBUY') {
						await pushSku({
							sku: SKU,
							site: 'BESTBUY',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]})
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
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'BIGLOTS') {
						await pushSku({
							sku: SKU,
							site: 'BIGLOTS',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]})
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
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'TARGETNEW') {
						await pushSku({
							sku: SKU,
							site: 'TARGETNEW',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]})
						console.log(kw)
							let currentBody = {
									  endpoint : SKU,
									  keywords : kw
							}
							console.log(currentBody)
							startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'ACADEMY') {
						await pushSku({
							sku: SKU,
							site: 'ACADEMY',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]})
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'AMD') {
						await pushSku({
							sku: SKU,
							site: 'AMD',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]})
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'SLICKDEALS' || site.toUpperCase() == 'SLICK' || site.toUpperCase() == 'SLICKDEAL') {
						await pushSku({
							sku: SKU,
							site: 'SLICKDEALS',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
						})
							let currentBody = {
								  	site: "Slick Deals",
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
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'WALMART NEW' || site.toUpperCase() == 'WALMARTNEW'){
						console.log(pricerange)
							let currentBody = {
								  	site: "Walmart New",
									skuName: "prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new",
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
							await pushSku({
							sku: SKU,
							site: 'WALMARTNEW',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
							})
							console.log(currentBody)
							startGoMonitor(currentBody, site.toUpperCase())
					}
					message.channel.send(`${SKU} Added to ${site}`)
					}
				}
				} else {
					message.channel.send(`${message.author} is not a validated user`)
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
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				console.log(site)
				console.log(`SKU - ${SKU}`)
				console.log(content)
				let skuBank = await getSkuBank()
				let caseSite = site.toUpperCase()
				let currentBody = skuBank[caseSite][SKU]
				const {group, isValidated} = await checkIfUserValidated(message)

				if(isValidated) {
					if(currentBody === undefined){
						console.log("undefined body")
						message.channel.send(`${SKU} Not Present \nCannot Delete ${SKU} from ${replaceWithTheCapitalLetter(site)}`)
					} else {
						console.log("Current length")
						console.log(currentBody.companies.length)
						message.channel.send(`Deleting ${SKU} from ${replaceWithTheCapitalLetter(site)}...`)
						if (currentBody.companies.length > 1){
							for(let i = 0; i < currentBody.companies.length; i++){
								if(currentBody.companies[i].company = group){
									console.log(currentBody.companies)
									currentBody.companies.splice(i, 1)
									console.log(currentBody.companies)
									await updateSku(caseSite, SKU, currentBody,  `${pushEndpoint}/${caseSite.toUpperCase()}/${SKU}.json`)}
									message.channel.send(`${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`)
								}
						} else {
							currentBody.stop = true
							console.log(currentBody)
							if (group == currentBody.companies[0].company){
								await updateSku(caseSite, SKU, currentBody,  `${pushEndpoint}/${caseSite.toUpperCase()}/${SKU}.json`)
								await delay(10000)
								await deleteSkuEnd(site, SKU)
								message.channel.send(`${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`)
							} else {
								console.log(`${group} is not present for this sku`)
								message.channel.send(`${SKU} Not Present \nCannot Delete ${SKU} from ${replaceWithTheCapitalLetter(site)}`)
								message.channel.send(`If this is an error please contact developer`)
							}
							
						}
						
					
					}
				} else {
					message.channel.send(`${message.author} is not a validated user`)
				}
				return;

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
				let validatedIds = await getValidatedIds()
				let parsed = JSON.parse(validatedIds)
				let isValidated = false
				let group
				parsed.ids.forEach(e =>{
					let id = e?.split('-')[0]
					console.log(id, message.author.id)
					if (id == message.author.id) {
						isValidated = true 
						group = e?.split('-')[1]
					}
				})
				if(isValidated){
				let skuBank = await getSkuBank()
				if(skuBank.length != 0){
					let bankArr = []
					let sites = Object.keys(skuBank)
					console.log(group)
					sites?.forEach(e =>{
						if (e != "-M_bpveXSTSxZkahEQkQ"){
							let currentSkus = Object.keys(skuBank[e])
							currentSkus.forEach(sku => {
								skuBank[e][sku]?.companies?.forEach(company =>{
									if (company.company == group){
										bankArr.push(`${e}-${sku}-${group}`)
									}
								})
							})
							}
					})
					console.log(bankArr)
					
				await fs.appendFile(`monitorBank-${message.author.username}.txt`, JSON.stringify(bankArr, null, 2) , err => {
						if(err) message.content.send('Error While Creating Text Document')
								else console.log("File Sent")
				})
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
				} else {
					message.channel.send(`${message.author} is not a validated user`)
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
				mass(string, content, message)
			}
		});
	} catch (error) {
		console.log(error);
	}
}
async function getSku (skuName, proxyList) {
        try {
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
					await delay(10000)
                    await getSku()
                }
}
async function startGoMonitor(currentBody, site){
	try {
		switch(site){
			// case "2341":
			case "TARGET":
			case "GAMESTOP":
				rp.post({
					url : `${secondServer}:${port}/${site}`,
					body : JSON.stringify(currentBody),
					headers : {
						"Content-Type": "application/json"
					}
				}, (response) => console.log(response?.statusCode))
				break
			case "WALMART":
			case "WALMARTNEW":
						rp.post({
							url : `${thirdServer}:${port}/${site}`,
							body : JSON.stringify(currentBody),
							headers : {
								"Content-Type": "application/json"
							}
						}, (response) => console.log(response?.statusCode))
						break
			default:
				rp.post({
					url : `http://localhost:${port}/${site}`,
					body : JSON.stringify(currentBody),
					headers : {
						"Content-Type": "application/json"
					}
				}, (response) => console.log(response?.statusCode))
		}
		
	} catch (error) {
		console.log(`Error Starting Go Monitor ${error}`)
	}
}
async function mass (string , content, message, groupName){
	//	const SKU = content.split(' ')[2];
			//	console.log(site)
			const site = content?.split(' ')[1]?.split('\n')[0].trim()
			console.log(content?.split(' ')[1]?.split('\n')[0].length)
			console.log(site.toUpperCase().length)
			console.log(site)
			let validatedIds = await getValidatedIds()
			let parsed = JSON.parse(validatedIds)
			let isValidated = false
			let group
			if(message){
				parsed.ids.forEach(e =>{
				let id = e?.split('-')[0]
				console.log(id, message.author.id)
				if (id == message.author.id) {
					isValidated = true 
					group = e?.split('-')[1]
				}
			})
			} else {
				isValidated = true
				group = groupName
			}
			if(isValidated){
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
					let original = g[i]
					if(g[i].includes('[')){
					pricerange = g[i].split('[')[1].split(']')[0]
					SKU = g[i].split(' ')[0]
					}
				
					console.log(g[i])
					console.log(site.toUpperCase())
					let skuBank = await getSkuBank()
					let caseSite = site.toUpperCase()
					let currentCompany = config[group]
					if(skuBank[caseSite]){
						if(skuBank[caseSite][SKU]){
							let skuWebhookArray = skuBank[caseSite][SKU]?.companies
							let isPresent = false
							skuWebhookArray.forEach(e => {
								console.log(e.webhook, currentCompany[caseSite])
								if(e.webhook == currentCompany[caseSite]) isPresent = true
							})
	
							if(isPresent){
								message.channel.send(`${SKU} is already present in monitor`)
							} else {
								message.channel.send(`${SKU} is being added to monitor`)
								let arr = []
								skuWebhookArray.forEach(e =>{
									arr.push(e)
								})
								arr.push({
									company : group,
									webhook : currentCompany[caseSite],
									color : currentCompany?.companyColor,
									companyImage : currentCompany?.companyImage
								})
								console.log(arr)
								await updateSku(site, SKU, {companies : arr}, `${pushEndpoint}/${caseSite.toUpperCase()}/${SKU}/.json`)
							}
							console.log('Duplicate Found', isPresent)
							isContinue = false
					}
					}
					if(isContinue){
					if (site.toUpperCase() == 'TARGET') {
						await pushSku({
							sku: SKU,
							site: 'TARGET',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'NEWEGG') {
						await pushSku({
							sku: SKU,
							site: 'NEWEGG' ,
        					stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
						startGoMonitor(currentBody, site.toUpperCase())

					} else if (site.toUpperCase() == 'GAMESTOP') {
						await pushSku({
							sku: SKU,
							site: 'GAMESTOP',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
						})
							let currentBody = {
									  site: "Game Stop",
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
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())

					} else if (site.toUpperCase() == 'AMD') {
							await pushSku({
							sku: SKU,
							site: 'AMD' ,
        					stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())

					} else if (site.toUpperCase() == 'WALMART') {
						await pushSku({
							sku: SKU,
							site: 'WALMART' ,
        					stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
							startGoMonitor(currentBody, site.toUpperCase())

						await delay(30000)
					} else if (site.toUpperCase() == 'BESTBUY') {
						await pushSku({
							sku: SKU,
							site: 'BESTBUY',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'ACADEMY') {
						await pushSku({
							sku: SKU,
							site: 'ACADEMY' ,
        					stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())

					} else if (site.toUpperCase() == 'SLICKDEALS' || site.toUpperCase() == 'SLICK' || site.toUpperCase() == 'SLICKDEAL') {
						await pushSku({
							sku: SKU,
							site: 'SLICKDEALS',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
						})
						// let monitor = new newEggMonitor(SKU.toString())
						// monitor.task()
							let currentBody = {
									  site: "Slick Deals",
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
							
						console.log(currentBody)
						startGoMonitor(currentBody, site.toUpperCase())
					} else if (site.toUpperCase() == 'WALMART NEW' || site.toUpperCase() == 'WALMART NEW'){
						console.log(pricerange)
							let currentBody = {
								  	site: "Walmart New",
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
							await pushSku({
							sku: SKU,
							site: 'WALMARTNEW',
							stop: false,
							name: "",
							original : original,
							companies : [
								{
								company : group,
								webhook : currentCompany[caseSite],
								color : currentCompany?.companyColor,
								companyImage : currentCompany?.companyImage
							}]
							})
							console.log(currentBody)
							startGoMonitor(currentBody, site.toUpperCase())
					}
					}
				await delay(30000)
				}
				
			}
				message.channel.send("SKUS Added")
			} else {
				message.channel.send(`${message.author} is not a validated user`)
			}
		
}
// Fire Base Sku Bank ----------------------------------------------
async function checkPresentSkus(){
	let skuBank = await rp.get({
		url : `${pushEndpoint}.json`
	})
	skuBank = JSON.parse(skuBank?.body)
	let sites = Object.keys(skuBank)
	let numberOfSkus = 0
	sites.forEach(site =>{
		let skus = Object.keys(skuBank[site])
		if(site != "-M_iJkLwZh3hW5Pjys5Z"){
			skus.forEach(sku =>{
			numberOfSkus++
		})
	}}
	)
	console.log(`Number Of Items In Monitor : ${numberOfSkus}`)
	await delay(3000)
	sites.forEach(async (e) => {
		if(e != "-M_iJkLwZh3hW5Pjys5Z"){
			let site = e
			let skus = Object.keys(skuBank[e])
			for(let i = 0; i < skus.length; i++) {
				let s = skus[i]
				let currentSku = skuBank[site][s].original
				let pricerange = ''
				if(currentSku?.includes('[') && site?.toUpperCase() !== "TARGETNEW"){
					pricerange = currentSku?.split('[')[1]?.split(']')[0]
					currentSku = currentSku?.split('[')[0]
				}
				let currentBody = {
				site: site,
				  sku: currentSku?.trim(),
				  priceRangeMin: parseInt(pricerange?.split(',')[0]),
				  priceRangeMax: parseInt(pricerange?.split(',')[1]),
				  skuName : site == 'NEWEGG' ? await getSku(currentSku?.trim(), await getProxies()) : site == "WALMARTNEW" ? "prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new" : ""
				  }
				  if(currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax){
					console.log("No Max Price Range Detected")
					currentBody.priceRangeMax = 100000
	
				} if(currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin){
					console.log("No Min Price Range Detected")
					currentBody.priceRangeMin = 1
	
				}
				console.log(currentBody)
				startGoMonitor(currentBody, site)
				await delay(15000)
			}
		}
		
		
	})

	
}
if(os.platform() == "win32" || os.platform() == "darwin"){
	console.log("Development Environment")
} else {
	checkPresentSkus()
}
async function getSkuBank(){
	let getbank = await rp.get({
		url : `${pushEndpoint}.json`
	})
//	console.log(JSON.parse(getbank?.body))
	return JSON.parse(getbank?.body)
}
async function pushSku(body){
	try {
		console.log("PUSHING ", body)
		let pushSku = await rp.patch({
		url : `${pushEndpoint}/${body.site}/${body.sku}.json`,
		body : JSON.stringify(body)
		})
		console.log(pushSku?.statusCode)
	} catch (error) {
		console.log(error.message)
	}
	
}
async function deleteSkuEnd(site, sku, group){
	try {
		console.log(`Deleting ${sku}/${site}`)
		let deleteSku = await rp.delete({
			url : `${pushEndpoint}/${site.toUpperCase()}/${encodeURIComponent(sku)}.json`
		})
		console.log(deleteSku?.statusCode)
	} catch (error) {
		console.log(error)
	}
}
async function updateSku(site, sku, newBody, url){
	try {
		// No need for site and sku // Only reason I kept it here is for console logs
		console.log(`Updating Sku ${sku}/${site}`)
		let updateSku = await rp.patch({
			url : url,
			body: JSON.stringify(newBody)
		})
		console.log(updateSku?.statusCode)
	} catch (error) {
		console.log(error)
	}
}
async function getValidatedIds (){
	try {
		let getIds = await rp.get({
		url : `${discordIds}/.json`
		})
		console.log(getIds?.statusCode)
		return getIds?.body
	} catch (error) {
		console.log(error)
	}
}
async function updateDiscordIdsDB (author, discordIdsArr, name){
try {
	let currentIds = []
	let parsed = JSON.parse(discordIdsArr)
	var isPresent = false
	console.log(parsed)
	parsed?.ids?.forEach(e => {
		currentIds.push(e)
		let id = e?.split('-')[0]
		if(id == author) isPresent = true
	})
	if(!isPresent) currentIds.push(`${author}-${name}`)
	else return false
	let updateIds = await rp.patch({
		url : `${discordIds}/.json`,
		body : JSON.stringify({ids : currentIds})
	})
	console.log(updateIds?.statusCode)
	return true
} catch (error) {
	console.log(error?.message)
	return false
}
}
async function validateUser(clients, triggerText, replyText){
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				let content = message.content.split('!validate')[1]
				let apikey = content?.split("apkey-")[1].split("%}")[0]
				if(apikey == undefined){
					message.reply("Please Submit Valid Api Key!")
				} else {
					console.log(apikey)
					let discordIdsDB = await getValidatedIds()
					switch(apikey){
						case "devAPIKEKg":
							message.reply("Valid Api Key")
							var update = await updateDiscordIdsDB(message.author.id, discordIdsDB, "DevTest")
							if (update) message.reply(`${message.author} Validated`)
							else message.reply(`${message.author} Could not be Validated \n Contact Dev if this continues to be an issue`)
							break
						case "j1ggedKRaFD#7d5e508f5e40":
							message.reply("Valid Api Key")
							var update = await updateDiscordIdsDB(message.author.id, discordIdsDB, "Jigged")
							if (update) message.reply(`${message.author} Validated`)
							else message.reply(`${message.author} Could not be Validated \n Contact Dev if this continues to be an issue`)
							break
						case "arialab91f37e6f36d2b2":
							message.reply("Valid Api Key")
							var update = await updateDiscordIdsDB(message.author.id, discordIdsDB, "Arial")
							if (update) message.reply(`${message.author} Validated`)
							else message.reply(`${message.author} Could not be Validated \n Contact Dev if this continues to be an issue`)
							break
						default:
							message.reply("Unknown Api Key")
					}
				}
				console.log(content)

		}});
	} catch (error) {
		console.log(error);
	}
}
async function checkIfUserValidated (message){
	let returnObj = {
		isValidated : false,
		group : undefined
	}
	let validatedIds = await getValidatedIds()
	let parsed = JSON.parse(validatedIds)
	parsed.ids.forEach(e => {
		let id = e?.split('-')[0]
		console.log(id, message.author.id)
		if (id == message.author.id) {
			returnObj.isValidated = true
			returnObj.group = e?.split('-')[1]
		}
	})
	return returnObj

}
//-----------------------------------------------------------------
function replaceWithTheCapitalLetter(values){
				return values.charAt(0).toUpperCase() + values.slice(1);
}
async function getProxies () {
	try {
// read contents of the file
let proxyList = []
const data = await fs.readFile('./GoMonitor/cloud.txt', 'utf-8');

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
});
console.log(`Proxy list Length : ${proxyList.length}`)
return proxyList;
} catch (err) {
console.error(err);
fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
console.log(err)
}))
}
}


// Scraper 

async function walmartScraper (clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				let validatedIds = await getValidatedIds()
				let parsed = JSON.parse(validatedIds)
				let isValidated = false
				let group
				parsed.ids.forEach(e =>{
					let id = e?.split('-')[0]
					console.log(id, message.author.id)
					if (id == message.author.id) {
						isValidated = true 
						group = e?.split('-')[1]
					}
				})
				let currentCompany = config[group]
				if(isValidated){
				console.log("Walmart Scraper Initiated")
				const SKU = message.content.split('/')[5]
				let currentMessage
				message.channel.send(`${message.author} Walmart Scraper for ${SKU} Initiated`).then(msg => {
					currentMessage = msg
				})
				let data = await getTerraSku(SKU)
				if(data == "Error") {
					currentMessage.edit(`${message.author} Error Grabbing Data for ${SKU}`)
					return
				}
				try {
				let ImageOffer = data.payload.selected.defaultImage
				let selectedProduct = data.payload.selected.product
				let productName = data.payload.products[selectedProduct].productAttributes.productName
				let ImageUrl = data.payload.images[ImageOffer].assetSizeUrls.DEFAULT
				let walmartSellerOffer

				Object.keys(data.payload.sellers).map(e =>{
					let currentSeller = data?.payload?.sellers[e]
					if (currentSeller?.sellerName == "Walmart.com") walmartSellerOffer = currentSeller?.sellerId
				})
				
				console.log(ImageUrl, walmartSellerOffer)
				let offerObj
				if(!walmartSellerOffer){
					currentMessage.edit(`${message.author} ${SKU} is not sold by Walmart`)
					offerObj = {
						offerId : "N/A",
						availabilityStatus : "Not Sold By Walmart",
						price : "N/A",
						image : ImageUrl,
						sku : SKU,
						productName : productName
					}
				} else {
					Object.keys(data.payload.offers).forEach(e =>{
					let currentOffer = data.payload.offers[e]
					if (currentOffer.sellerId == walmartSellerOffer)  offerObj = {
						offerId : currentOffer.id,
						availabilityStatus : currentOffer.productAvailability.availabilityStatus,
						price : currentOffer.pricesInfo.priceMap.CURRENT.price,
						image : ImageUrl,
						sku : SKU,
						productName : productName
					}
				})
				}
				
				var dt = new Date();
				// let body = {
				// 	"content": null,
				// 	"embeds": [
				// 	  {
				// 		"title": "Walmart Scraper",
				// 		"color": currentCompany.companyColor,
				// 		"fields": [
				// 		  {
				// 			"name": "Offer Id",
				// 			"value": offerObj.offerId
				// 		  },
				// 		  {
				// 			"name": "Availability",
				// 			"value": offerObj.availabilityStatus,
				// 			"inline": true
				// 		  },
				// 		  {
				// 			"name": "Walmart's Price",
				// 			"value": offerObj.price,
				// 			"inline": true
				// 		  },
				// 		  {
				// 			"name": "Links",
				// 			"value": `[Product](https://www.walmart.com/ip/prada/${offerObj.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${offerObj.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`
				// 		  }
				// 		],
				// 		"footer": {
				// 		  "text": "Prada#4873"
				// 		},
				// 		"timestamp": dt.toISOString(),
				// 		"thumbnail": {
				// 		  "url": offerObj.image
				// 		}
				// 	  }
				// 	],
				// 	"avatar_url": currentCompany.companyImage
				//   }
				const currentEmbed = new MessageEmbed()
				.setColor(currentCompany.companyColorV2)
				.setTitle(offerObj.productName)
				.setURL(`https://www.walmart.com/ip/prada/${offerObj.sku}`)
				.setThumbnail(offerObj.image)
				.addFields(
					{
						name: "Offer Id",
						value: offerObj.offerId
					  },
					  {
						name: "PID",
						value: offerObj.sku
					  },
					  {
						name: "Availability",
						value: offerObj.availabilityStatus,
						inline: true
					  },
					  {
						name: "Walmart's Price",
						value: offerObj.price,
						inline: true
					  },
					  {
						name: "Links",
						value: `[Product](https://www.walmart.com/ip/prada/${offerObj.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${offerObj.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`
					  }
					)
				.setFooter('Prada#4873')
				.setTimestamp()
				message.channel.send(currentEmbed)
				currentMessage.delete()
				} catch (error) {
					console.log(error)
					currentMessage.edit(`${message.author} Error Grabbing Data for ${SKU} BOJ`)

				}
				
				} else {
					message.channel.send(`${message.author} is not a validated user`)
				}
			
			}
		})
	} catch (error) {
		console.log(error)
	}
}

async function getTerraSku (SKU) {
	try {

	let terra = await rp.get({
		url : `https://www.walmart.com/terra-firma/item/${SKU}`,
		headers : {
			"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"accept-language": "en-US,en;q=0.9",
			"cache-control": "no-cache",
			"pragma": "no-cache",
			"sec-ch-ua": "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"",
			"sec-ch-ua-mobile": "?0",
			"sec-fetch-dest": "document",
			"sec-fetch-mode": "navigate",
			"sec-fetch-site": "none",
			"sec-fetch-user": "?1",
			"service-worker-navigation-preload": "true",
			"upgrade-insecure-requests": "1",
			"user-agent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
		},
		proxy : "http://CloudsNYadrianzsGUqR7R:VwfvIFQR@154.28.145.201:2400"
	})
	console.log(terra.statusCode)
	if(terra.statusCode != 200 && terra.statusCode != 201){
		return "Error"
	} else {
		let parsed = JSON.parse(terra?.body)
		return parsed
	}
	} catch (error) {
		console.log(error)
		return "Error"
	}
}

async function sendWebhook(body, webhook) {
	try {
		let sent = await rp.post({
			url : webhook,
			body : body
		})
		return sent.statusCode
	} catch (error) {
		console.log(error)
		return error?.statusCode
	}
}
module.exports = {
	SKUADD,
	findCommand,
	deleteSku,
	checkBank,
	massAdd,
	validateUser,
	walmartScraper
} 