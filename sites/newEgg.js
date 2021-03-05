const delay = require('delay');
const cheerio = require('cheerio');
const axios = require('axios').default;
const fs = require('fs');
const { gzip } = require('zlib');
const { response } = require('express');
const { resolve } = require('path');
const { json } = require('body-parser');
const fetch = require('node-fetch');
const Discord = require('discord.js');
// require ('newrelic');
const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
});
// https://discordapp.com/api/webhooks/816740348222767155/2APr1EdhzNO4hRWznexhMRlO0g7qOiCkI7HFtmuU7_r48PCWnGYmSTGJmRVX0LPCNN_t
 const webhookClient1 = new Discord.WebhookClient('816740348222767155', '2APr1EdhzNO4hRWznexhMRlO0g7qOiCkI7HFtmuU7_r48PCWnGYmSTGJmRVX0LPCNN_t');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');
// const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');

const webhook = require("webhook-discord")
 // https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
const Hook = new webhook.Webhook("https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM")

class newEggMonitor {
    constructor(sku) {
        this.trueSku = sku
        this.sku = sku.split(':')[0];
        this.skuName = sku.split(':')[1]
        this.delay = 175000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false
        this.imageUrl = ''

    }

    async task () {
        try {
            console.log('Start')
            console.log(this.sku)
            console.log(this.skuName)
            await this.getProxies()
         //   console.log(this.proxyList)
            await this.monitor()
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
            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                console.log(err)
            }))
        }
    }

    async monitor () {
       try {
        console.log('Starting Monitoring')
        var testing = ''
        return new Promise( async ( resolve, reject) => {
            for (let i = 0; i < this.proxyList.length; i++){
                let proxy = this.proxyList[i]
                console.log(`${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`)
                let monitorInterval = setInterval(async () => {
                    var { skuBank } = require('../dms')
                    let index = skuBank.findIndex(e => e.sku == this.trueSku)
                    if(skuBank[index]?.stop){
                        console.log('stopped!')
                        clearInterval(monitorInterval)
                        resolve('Stopped')
                        return;
                    }
                    try {
                        let fetchSite = await rp.get({
                            url : `https://www.newegg.com/product/api/ProductRealtime?ItemNumber=${this.sku}`,
                            proxy : `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                        })   
                       // console.log(fetchSite.body)
                        console.log(`Task ${i} : ${fetchSite.statusCode}`)
                        testing = fetchSite.body
                        let parsedBod = JSON.parse(fetchSite.body)
                        let image = parsedBod?.MainItem?.Image?.ItemCellImageName
                        
                        let productName = parsedBod?.MainItem?.Description?.Title
                        let originalPrice = parsedBod?.MainItem?.OriginalUnitPrice
                        let currentPrice = parsedBod?.MainItem?.FinalPrice
                        this.availability = parsedBod?.MainItem?.Instock
                        this.stockNumber = parsedBod?.MainItem?.Stock 
                        console.log(this.availability, this.stockNumber, productName)


                        if(!this.isStock && this.availability) {
                            // Send in stock webhook
                            this.isStock = true
                            let embed1 = new Discord.MessageEmbed()
                            .setColor('#07bf6e')
                            .setTitle('New Egg Monitor')
                            .setURL(`https://www.newegg.com/Prada/p/${this.skuName}`)
                            .addFields(
                                { name : 'Product Name', value : `${productName}`},
                                { name : 'Product Availability', value : `Product In Stock`, inline : true},
                                { name : 'Stock Number', value : `${this.stockNumber}`, inline : true}, 
                                { name : 'Current Price', value : currentPrice, inline : true} 
                                )
                            .addField("Links", `[ATC](https://secure.newegg.com/shopping/addtocart.aspx?submit=add&itemList=${this.sku}) | [Cart](https://secure.newegg.com/shop/cart)`)
                            // .addField('Original Price', originalPrice, true)
                            // .addField('Current Price', currentPrice, true)
                            .setThumbnail(`https://c1.neweggimages.com/ProductImageOriginal/${image}`)
                            .setTimestamp()
                            .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                            webhookClient1.send('Restock!', {
                                username: 'New Egg',
                                avatarURL: 'https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png',
                                embeds: [embed1],
                            })
                        } else if (!this.availability && this.isStock) {
                            this.isStock = false
                        } 
    
                    
                        // console.log(originalPrice, currentPrice)
                       // console.log(parsedBod)
        
        
                    } catch (error) {
                        console.log(error)
                        fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                            console.log(err)
                        }))
                          if(error.message.includes('Unexpected token')){
                            console.log(testing)
                            clearInterval(monitorInterval)
                            resolve('g')
                        }
                    }
                }, this.delay)
                await delay(3000)
            }
          

        })
       } catch (error) {
        fs.appendFileSync('bigError.txt', error.toString() + '\n', (err =>{
            console.log(err)
        }))
       }
  
    }

}

// const monitoring = new newEggMonitor(`19-113-569`);

// (async ()=>{
//     await monitoring.task()
// })()

module.exports = {
    newEggMonitor
}