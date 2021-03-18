const delay = require('delay');
const cheerio = require('cheerio');
const axios = require('axios').default;
const fs = require('fs');
const {
    gzip
} = require('zlib');
const {
    response
} = require('express');
const {
    resolve
} = require('path');
const {
    json
} = require('body-parser');
const fetch = require('node-fetch');
const Discord = require('discord.js');
// require ('newrelic');
const rp = require('request-promise').defaults({
    followAllRedirects: true,
    resolveWithFullResponse: true,
    gzip: true,
});
const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');

const webhook = require("webhook-discord")
// https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
const Hook = new webhook.Webhook("https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM")



class targetMonitor {
    constructor(sku) {
        this.sku = sku;
        this.delay = 1000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false;
        this.productName = '';
        this.itemPicUrl = ''

    }

    async task() {
        try {
            console.log('Start')
            await this.getProxies()
            //   console.log(this.proxyList)
            await this.monitor()
        } catch (error) {
            fs.appendFileSync('errors.txt', error.toString() + '\n', (err =>{
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
                    ip: lineSplit[0],
                    port: lineSplit[1],
                    userAuth: lineSplit[2],
                    userPass: lineSplit[3],
                };
                this.proxyList.push(item1);

                // console.log(line);

                // console.log('\n\n\n\n\n');
            });
        } catch (err) {
            console.error(err);
            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                console.log(err)
            }))
        }
    }
    async monitor() {
        console.log('Starting Monitoring')
        var testing = ''
        return new Promise(async (resolve, reject) => {
            let getItemPic = await rp.get({
                url: `https://www.amd.com/en/direct-buy/${this.sku}/us`
            }, ((error, response, body) => {
                // console.log(body)
                const $ = cheerio.load(response.body)
                let itemPic = $('#product-details-info > div.container > div > div.product-page-image.col-flex-lg-7.col-flex-sm-12 > img')[0]                                                                                                
                this.productName = itemPic.alt
                this.itemPicUrl = itemPic.href
                console.log(this.productName, this.itemPicUrl)
            }))
            for (let i = 0; i < this.proxyList.length; i++) {
                let proxy = this.proxyList[i]
                let monitorInterval = setInterval(async () => {
                    var {
                        skuBank
                    } = require('../dms')
                    let index = skuBank.findIndex(e => e.sku == this.sku)
                    if (skuBank[index].stop) {
                        console.log('stoppped!!!!!')
                        clearInterval(monitorInterval)
                        resolve('Stopped')
                        return;
                    }
                    try {
                        let check = await rp.get({
                            url : `https://www.amd.com/en/direct-buy/${this.sku}/us`,
                            proxy: `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                        })
                        let test = await rp.get({
                            url = 'https://www.amd.com/en/direct-buy/add-to-cart/5450881600',
                            proxy: `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                        }, ((error, response, body)=>{
                            console.log(body)
                        }))
                        const $ = cheerio.load(check.body)
                        let atcCheck = $('#product-details-info > div.container > div > div.product-page-description.col-flex-lg-5.col-flex-sm-12 > button')
                        let length  = atcCheck.length
                        console.log(length)  
                        if(length == 1){
                            console.log(atcCheck[0])
                            this.availability = true
                        } else if (length == 0){
                            this.availability = false
                        } else {
                            fs.appendFileSync('./errors.txt', this.availability + this.availability + '\n', (err =>{
                                console.log(err)
                            }))
                        }
                        console.log(check.statusCode)


                        if (!this.isStock && this.availability) {
                            // Send in stock webhook
                            this.isStock = true
                            let embed1 = new Discord.MessageEmbed()
                                .setColor('#07bf6e')
                                .setTitle('Target Monitor')
                                .setThumbnail(`${this.itemPicUrl}`)
                                .setURL(`https://www.target.com/prada/-/A-${this.sku}`)
                                .addField('Product Name', `${this.productName}`)
                                .addField('Product Availability', 'In Stock!', true)
                              //  .addField('Stock Number', `${this.stockNumber}`, true)
                                .addField("Links", `[Product](https://www.amd.com/en/direct-buy/${this.sku}/us) | [Add To Cart](https://www.amd.com/en/direct-buy/add-to-cart/${this.sku})`)
                                // .addField('Original Price', originalPrice)
                                // .addField('Current Price', currentPrice , true)
                                //  .setImage(`${this.itemPicUrl}`)
                                .setTimestamp()
                                .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                            webhookClient1.send('Restock!', {
                                username: 'Target',
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

                        console.log(`Task ${i} : ${this.availability}, ${this.productName}}`)
                        // console.log(originalPrice, currentPrice)
                        // console.log(parsedBod)


                    } catch (error) {
                        console.log(error)
                        fs.appendFileSync('errors.txt', error.toString() + '\n', (err =>{
                            console.log(err)
                        }))
                        // Add fs errors
                        if (error.message.includes('Unexpected token')) {
                            console.log(testing)
                            clearInterval(monitorInterval)
                            resolve('g')
                        }
                    }
                }, this.delay)
                await delay(100)
            }


        })
    }


}

// const monitoring = new targetMonitor(81409304);

// (async ()=>{
//     await monitoring.task()
// }) ()

module.exports = {
    targetMonitor
}