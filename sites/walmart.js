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
//const { chromium } = require('playwright')
const puppeteer = require('puppeteer')
const request = require('request')
var tough = require('tough-cookie');
var Cookie = tough.Cookie
// require ('newrelic');
var cookieJar = request.jar()
const rp = require('request-promise').defaults({
    followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
    timeout: 5000 
});
// https://discordapp.com/api/webhooks/826289643455643658/tRuYU2WQGSoyD5gH2QL8dKecI59F8IyH_wds5_pio7pOst79cBWs6wEe0jdkGI1qeYMC 
  const webhookClient1 = new Discord.WebhookClient('826289643455643658', 'tRuYU2WQGSoyD5gH2QL8dKecI59F8IyH_wds5_pio7pOst79cBWs6wEe0jdkGI1qeYMC');



 //Test
 // const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');


class walmartMonitor {
    constructor(sku, price) {
        this.trueSku = sku
        this.price = price
        this.sku = sku
        this.delay = 850000; // this.delay = 390000
        this.startDelay = 3000; //  this.startDelay = 6000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false
        this.imageUrl = ''
        this.maxPrice = parseInt(price?.split(' ')[1])
        this.minPrice =  parseInt(price?.split(' ')[0])
    }

    async task () {
        try {
            console.log('Start')
            console.log(this.sku)
            if(this.price?.length > 2){
            console.log('Price Range Detected')
            } else {
                this.maxPrice = 1000000000
                this.minPrice = 1
            }
            console.log(this.maxPrice, this.minPrice)
            await this.getProxies()
            await this.monitor()
        } catch (error) {
            console.log(error)
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
                            skuBank[index].name = this.productName
                        // console.log(`${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`)

                        try {
                            
                            let monitoring = await rp.get({
                                headers : {
                                                "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
                                                "accept-language": "en-US,en;q=0.9",
                                                "cache-control": "max-age=0",
                                                "sec-fetch-dest": "document",
                                                "sec-fetch-mode": "navigate",
                                                "sec-fetch-site": "none",
                                                "sec-fetch-user": "?1",
                                                "upgrade-insecure-requests": "1",
                                                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
                                   
                                },
                                url : `https://www.walmart.com/terra-firma/item/${this.sku}`,
                                proxy : `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                            })

                            let availableOffer = 
                            {
                              offer : "",
                              price: "",
                              availability : ""
                            }
                         //   console.log()
                            testing = monitoring?.body
                            let parsed = JSON.parse(monitoring?.body)
                            let offerList = parsed?.payload?.offers
                            let keys = Object.keys(offerList)
                            let currentImages = Object.keys(parsed?.payload?.images)

                            this.imageUrl = (parsed?.payload?.images[currentImages[0]]?.assetSizeUrls?.DEFAULT)
                            this.availability = false
                            this.productName = parsed?.payload?.products[Object.keys(parsed?.payload?.products)[0]]?.productAttributes?.productName
                                keys.map(e => {
                                 let d = offerList[e]
                                 let status = d?.productAvailability?.availabilityStatus
                                 let price = d?.pricesInfo?.priceMap?.CURRENT?.price

                                if(status == "IN_STOCK" && this.minPrice < price && price < this.maxPrice){
                                    this.availability = true
                                  //  console.log(this.availability)
                                    availableOffer.offer = e
                                    availableOffer.price = price
                                    availableOffer.availability = status
                                }
                                })
                                console.log(`Task ${i} | ${this.sku} | ${monitoring?.statusCode}: ${JSON.stringify(availableOffer)} ${this.availability} & ${this.isStock}`)
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
                                        .addField('Price', availableOffer.price, true)
                                        .addField('Sku', this.sku, true)
                                        .addField('OfferId', availableOffer.offer , true)
                                //      .addField('Stock Number', `${this.stockNumber}`, true)
                                        .addField("Links", `[Product](https://www.walmart.com/ip/prada/${this.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${this.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`)
                                        //  .setImage(`${this.itemPicUrl}`)
                                        .setTimestamp()
                                        .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                                        webhookClient1.send({
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
                               // console.log(testing)
                              //  resolve('g')
                            } else if (error.message.includes('403')){
                                // await delay(400000)
                                console.log('403 Access Denied')
                            } else if(error?.message?.includes('Cannot convert undefined or null to object')){
                                console.log('Walmart Undefined')
                            } else {
                                console.log(testing)
                            }
                            
                        }

                        await delay(this.startDelay)
                    }
                    console.log('stopped!')
                    resolve('Stopped')
                    return;
              
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