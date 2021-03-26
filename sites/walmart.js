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
const { chromium } = require('playwright')
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

class walmartMonitor {
    constructor(sku) {
        this.trueSku = sku
        this.sku = sku.split(':')[0];
        this.skuName = sku.split(':')[1]
        this.delay = 850000; // this.delay = 390000
        this.startDelay = 100; //  this.startDelay = 6000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false
        this.imageUrl = ''
        this.maxPrice = ''
        this.minPrice = ''

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
            const data = fs.readFileSync('proxies.txt', 'utf-8');
			// read contents of the file
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

                let i = 0
                var { skuBank } = require('../dms')
                let index = skuBank.findIndex(e => e.sku == this.trueSku)
                    while(!skuBank[index]?.stop){
                        if(i+1 == this.proxyList.length){
                            i = 0
                        }
                        let proxy = this.proxyList[i]
                        i++
                        console.log(`${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`)
                        try {
                            
                            let monitoring = await rp.get({
                                url : `https://www.walmart.com/terra-firma/item/${this.sku}`,
                                proxy : `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                            })
                            let availableOffer = {offer, price, availability}
                         //   console.log()
                            let parsed = JSON.parse(monitoring?.body)
                            let offerList = parsed.payload.offers
                            let keys = Object.keys(offerList)
                            let currentImages = Object.keys(parsed.payload.images)
                            this.imageUrl = (parsed.payload.images[currentImages].assetSizeUrls.DEFAULT)
                            this.availability = false
                            this.productName = parsed.payload.products[Object.keys(parsed.payload.products)[0]].productAttributes.productName
                            keys.map(e => {
                                 let d = offerList[e]
                                 let status = d?.productAvailability.availabilityStatus
                                 let price = d?.pricesInfo.priceMap.CURRENT.price

                                 if(status == "IN_STOCK" && this.minPrice < price && price < this.maxPrice){
                                    this.availability == true
                                    console.log(this.availability)
                                    availableOffer.offer = e
                                    availableOffer.price = price
                                    availableOffer.availability = status
                                }
                                })
                                if (!this.isStock && this.availability) {
                                    // Send in stock webhook
                                    this.isStock = true
                                    let embed1 = new Discord.MessageEmbed()
                                        .setColor('#07bf6e')
                                        .setTitle('Walmart Monitor')
                                        .setThumbnail(`${this.imageUrl}`)
                                        .setURL(`https://www.walmart.com/ip/prada/${this.sku}`)
                                        .addField('Product Name', `${this.productName}`)
                                        .addField('Product Availability', 'In Stock!', true)
                                //        .addField('Stock Number', `${this.stockNumber}`, true)
                                        .addField("Links", `[Product](https://www.walmart.com/ip/prada/${this.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`)
                                        .addField('Price', availableOffer.price)
                                         .addField('OfferId', availableOffer.offer , true)
                                        //  .setImage(`${this.itemPicUrl}`)
                                        .setTimestamp()
                                        .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                                    webhookClient1.send('Restock!', {
                                        username: 'Walmart',
                                        avatarURL: 'https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png',
                                        embeds: [embed1],
                                    })
                                } else if (!this.availability && this.isStock) {
                                    this.isStock = false
                                } else if (!this.availability && !this.isStock){
                                    this.isStock = false
                                    this.availability = false
                                    // Not Important
                                  //  console.log('False false')
                                } else if(this.availability && this.isStock){
                                    this.availability = true
                                    this.isStock = true
                                    // Not Important
                                    // console.log(true, true)
                                } else {
                                    fs.appendFileSync('what.txt', this.availability + this.isStock + '\n', (err =>{
                                        console.log(err)
                                    }))
                                }
                        } catch (error) {
                            console.log(error)
                            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                                console.log(err)
                            }))
                              if(error.message.includes('Unexpected token')){
                                console.log(testing)
                                resolve('g')
                            } else if (error.message.includes('403')){
                                await delay(400000)
                                console.log('403 Access Denied')
                            }
                        }
                        await delay(this.startDelay)
                    }
                    console.log('stopped!')
                    resolve('Stopped')
                    return
              
        })
       } catch (error) {
        fs.appendFileSync('bigError.txt', error.toString() + '\n', (err =>{
            console.log(err)
        }))
       }
  
    }

    

}


module.exports = {
    walmartMonitor
}