const delay = require('delay');
const cheerio = require('cheerio');
const axios = require('axios').default;
const fs = require('fs');
const { gzip } = require('zlib');
const { response } = require('express');
const { resolve } = require('path');
const fetch = require('node-fetch');
const Discord = require('discord.js');
const request = require('request')
// const {chromium} = require('playwright')
const { chromium } = require("playwright-chromium");
var tough = require('tough-cookie');
var Cookie = tough.Cookie;
// require ('newrelic');
var cookieJar = request.jar()
const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
    jar : cookieJar
});
const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');

const webhook = require("webhook-discord");
 // https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
const Hook = new webhook.Webhook("https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM")
function getChromiumExecPath() {
    let platform = process.platform;
    if (platform === "win32") {
        return "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
    } else if (platform === "darwin") {
        return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome";
    } else {
        return chromium.executablePath().replace("app.asar", "app.asar.unpacked");
    }
}
class gameStopMonitor {
    constructor(sku) {
        this.sku = sku;
        this.delay = 10000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false
        this.browser = ''
        this.stop = ''

    }

    async task () {
        try {
            console.log('Start')
            await this.getProxies()
         await this.getCookies()
         
         //   await this.monitor()
        } catch (error) {
            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                console.log(err)
            }))
        }
    }
    async getProxies() {
		try {
			// read contents of the file
			const data = fs.readFileSync('proxies.txt', 'utf-8');
			// split the contents by new line
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
				this.proxyList.push(item1);

				// console.log(line);
				// console.log(item1);
				// console.log('\n\n\n\n\n');
			});
		} catch (err) {
			console.error(err);
		}
    }
    async getCookies (){
        try {
            // ::
            console.log('Fetching Page')
            this.browser = await chromium.launch({
                headless: false,
                chromiumSandbox: false,
                proxy : {server : `134.195.155.61:8080`, username : 'bfe6yp', password : 'gsuep4'},
                timeout : 0
            });
            console.log('First Check')
            const context = await this.browser.newContext();
            console.log('Second Check')
            const page = await context.newPage();
            console.log('Third Check')
            await page.goto(`https://www.gamestop.com/pradaIsTalented/${this.sku}.html`);
            console.log('Fourth Check')
            await page.waitForSelector('h1')
            console.log('Fetching Cookies')
            const cookies = await context.cookies();
            console.log('Cookies Fetched')
         //   console.log(cookies)
            fs.writeFileSync('./test.txt', cookies.toString(), (err => console.log(err)))
            cookies.forEach(json => {
                                const { name, domain } = json
                                json.key = name
                                json.expires = json.expires > 0 ? new Date(json.expires * 1000) : 'Infinity'
                                const cookie = Cookie.fromJSON(json)
                                let url = 'https://' + domain
                                cookieJar.setCookie(cookie.toString(), url)
                            })

            console.log('Cookies Stored')
                // console.log(cookieJar)
            let body = ''
            let reload = setInterval(async () => {
                var { skuBank } = require('../dms')
                let index = skuBank.findIndex(e => e.sku == this.sku)
                if(skuBank[index].stop){
                    console.log('stoppped!!!!!')
                    clearInterval(reload)
                    resolve('Stopped')
                    return;
                }
                await page.reload()
            }, 600000)
            for (let i = 0 ; i < this.proxyList.length; i++){
                let proxy = this.proxyList[i]
                console.log(`${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`)
                let monitorInterval = setInterval(async () => {
                    var { skuBank } = require('../dms')
                    let index = skuBank.findIndex(e => e.sku == this.sku)
                    if(skuBank[index].stop){
                        console.log('stoppped!!!!!')
                        clearInterval(monitorInterval)
                        resolve('Stopped')
                        return;
                    }
                    try {
                        
                        let stock = await page.evaluate(async (sku) => {
                            try {
                                let test  = await fetch(`https://www.gamestop.com/on/demandware.store/Sites-gamestop-us-Site/default/Product-Variation?pid=${sku}&redesignFlag=true&rt=productDetailsRedesign`).then(res => res.text().then(res => {
                                    body = res
                               //     console.log(res)
                                }))
                                
                            } catch (error) {
                                console.log(error)
                                fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                                    console.log(err)
                                }))
                            }
                           return body
                           
                    }, this.sku)
                   // console.log(stock.statusCode)
                    let parsedBody = JSON.parse(stock)
                    let productName = parsedBody?.gtmData?.productInfo?.name
                    let productSku = parsedBody?.gtmData?.productInfo?.sku
                    let originalPrice = parsedBody?.gtmData?.price?.basePrice
                    let currentPrice = parsedBody?.gtmData?.price?.sellingPrice
                    let image = parsedBody?.product?.images?.large[0]?.url
                    this.availability = parsedBody?.gtmData?.productInfo?.availability
                    console.log({
                        task : i,
                        productName : productName,
                        productSku : productSku,
                        originalPrice : originalPrice,
                        currentPrice :  currentPrice,
                        image : image,
                        availability : this.availability
                    })
                    if(this.availability == 'Available'){
                    this.availability = true
                    } else if (this.availability == 'Not Available'){
                        this.availability = false
                    }
        
        
                    if(!this.isStock && this.availability){
                         // Send in stock webhook
                         this.isStock = true
                         let embed1 = new Discord.MessageEmbed()
                         .setColor('#00FF00')
                         .setTitle('Game Stop Monitor')
                         .setURL(`https://www.gamestop.com/Prada/${this.sku}.html`)
                         .addField('Product Name', `${productName}`)
                         .addField('Product Availability', 'Product In Stock', true)
                         .addField('Product Pid', productSku , true)
                         .addField('Original Price', originalPrice)
                         .addField('Current Price', currentPrice)
                         .setImage(`${image}`)
                         .setTimestamp()
                         .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                         webhookClient1.send('Restock!', {
                             username: 'Game Stop',
                             avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
                             embeds: [embed1],
                         })
                    } else if (!this.availability && this.isStock) {
                        this.isStock = false
                    }
        
        
                    } catch (error) {
                        fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                            console.log(err)
                        }))
                        console.log(error)
                    }
                }, 1000)
                await delay(100)
            }
           
        } catch (error) {
            console.log(error)
        }
       
        
       
        

    }


}

// const monitoring = new gameStopMonitor(`B158467A`);

// (async ()=>{
//     await monitoring.task()
// }) // ()


module.exports = {
    gameStopMonitor
}