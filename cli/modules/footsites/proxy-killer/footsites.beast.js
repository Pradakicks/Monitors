/** @format */
const chalk = require('chalk')
const {
    v4
} = require("uuid");
const fs = require("fs");
var https = require("https");
// const adyenEncrypt = require("node-adyen-encrypt")(24);
const adyenEncrypt = require("../node-adyen-encrypt/index")(18)


// const DB = require('../../file-crud/db');


const sendNotification = require("../../../lib/discord")
const { getProxy, proxyIsNotInUse } = require('../../../lib/proxies');
const harvest2Captcha = console.log('testtin') //require("../../../utils/captcha/recaptcha");
const generateDatadome = require("../utils/datadome");

require("colors");
require("console-stamp")(console, {
    colors: {
        stamp: "blue",
        label: "magenta",
    },
});

const util = require('util');
require('util.promisify').shim();

class Footsites {
    constructor(task, mainWindow) {
        this.mainWindow = mainWindow;
        this.task = task;
        this.site = task.site;
        this.pid = task.pid.replace(/\s/g, "");
        this.size = task.size;
        this.monitorDelay = task.monitorDelay;
        this.navigationDelay = task.navigationDelay;
        this.restockDelay = task.restockDelay;
        this.csrf = "";
        this.productInfo = {};
        this.cartId = null;
        this.ddCookieRes = "";
        this.aydenEncryptedCard = {};
        this.atcRetry = 5;
        this.emailRetry = 5;
        this.shippingRetry = 5;
        this.billingRetry = 5;
        this.paymentRetry = 3;
        this.stop = false;
        this.proxy = null;
        this.unFormattedproxy = null;
        this.profile = task.profile;
        this.body = ''
        this.bod1y = ''
        this.jar = require("request-promise").jar();
        this.request = require("request-promise").defaults({
            jar: this.jar,
            followAllRedirects: true,
            agent: new https.Agent({
                host: `www.${this.site}.com`,
                port: "443",
                path: "/",
                rejectUnauthorized: false,
                timeout : 200000
            }),
        });
    }
    async start() {



        try {
            console.log(`Starting ${this.site} task`);
            await this.setProxy();
            console.log(`Getting Session`);
            await this.initializeSession();
       //     console.log(`Initiallized: ${this.csrf}`);
            console.log(`Initiallized`);
            console.log(`Getting product with PID: ${this.pid}`);
            await this.delay(this.navigationDelay);
            console.log('Getting Details')
            await this.getItemDetails();
            if (!this.productInfo) return;
      //      console.log(`Found product:`, this.productInfo);
            await this.delay(this.navigationDelay);
            console.log(`Carting product with code: ${this.productInfo.variant}`);
            await this.addTocart()
            await this.delay(this.navigationDelay);
            console.log(`Carted product: ` + this.cartId);
            await this.delay(this.navigationDelay);
            await this.submitEmail();
            console.log(`Submitted Email: ${this.profile.email}`);
            await this.delay(this.navigationDelay);
            await this.submitShipping();
            console.log(`Submitted Shipping`);
            await this.delay(this.navigationDelay);
            await this.submitBilling();
            console.log(`Submitting payment`);
            await this.submitOrder();
        } catch (error) {
            console.log(error);
            this.atcRetry = 5;
            this.emailRetry = 5;
            this.shippingRetry = 5;
            this.billingRetry = 5;
            this.paymentRetry = 3;
            await proxyIsNotInUse(this.unFormattedproxy);
            console.log(`Restarting task`);
            this.jar = require("request-promise").jar();
            await this.delay(5000)
            await this.start()
        }
    }

    async setProxy() {
        try {
            if (this.stop) {
                return
            }
            if (true) {
                console.log('setting proxy')
                let proxies = await getProxy();
              //  console.log(proxies)
                this.proxy = proxies.formatted;
                this.unFormattedproxy = proxies.unformatted
            }
        } catch (error) {
            console.log(error)
            console.log('error setting proxy')
        }
    };

    async initializeSession() {

        if (this.stop) return;
        try {
            let url = `https://www.${this.site
                }.com/api/v3/session?timestamp=${Date.now()}`;
         //   console.log(url);
            const options = {
                method: "GET",
                uri: url,
                headers: {
                    Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
                    "Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    Connection: "keep-alive",
                    Host: `${this.site}.com`,
                    "Sec-Fetch-Dest": "document",
                    "Sec-Fetch-Mode": "navigate",
                    "Sec-Fetch-Site": "none",
                    "Sec-Fetch-User": "?1",
                    "Upgrade-Insecure-Requests": "1",
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
                },
                jar: this.jar,
                proxy: this.proxy,
                rejectUnauthorized: false,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: body,
                        status: response.statusCode
                    };
                },
            };
            const res = await this.request(options);
         //   console.log(res)
            if (res.status == 503) {
                throw new Error({
                    status: 503,
                    error: "queue"
                });
            }
            const response = JSON.parse(res.data);
            if (response.data) {
                this.csrf = response.data.csrfToken;
            } else {
              //  console.log(response);
                throw new Error("unable to initial session");
            }
        } catch (error) {
         //   console.log(error)
        
            if (error.response?.status == 503) {
                console.log("in queue");
                await this.delay(30000);
                await this.initializeSession();
            } else if (error.response?.status == 429){
                console.log('Too MAny Requests')
                await this.delay(30000)
                await this.initializeSession();
            } else if (error.message.includes('tunneling socket could not be established')){
                console.log('Handled Socket Error')
                console.log(error.message)
                await this.delay(50000)
                await this.initializeSession();
            } else {
                console.log(error.message)
                console.log("Proxy Error Retrying");
            //    console.log(error);
                await this.setProxy()
                await this.delay(1000);
                await this.initializeSession();
                // throw new Error("Blocked: try new proxy");
            }
        }
    }
    
    async getItemDetails() {
        try {
            if (this.stop) return;
            const url = `https://www.${this.site}.com/api/products;;/pdp/${this.pid
                }?timestamp=${Date.now()}`;
        //    console.log(url);
            let reqoptions = {
                method: "GET",
                uri: url,
                headers: {
                    Accept: "application/json",
                    "Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    Connection: "keep-alive",
                    Host: `${this.site}.com`,
                    "Sec-Fetch-Dest": "empty",
                    "Sec-Fetch-Mode": "cors",
                    "Sec-Fetch-Site": "none",
                    "Sec-Fetch-User": "?1",
                    "Upgrade-Insecure-Requests": "1",
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
                },
                jar: this.jar,
                proxy: this.proxy,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: JSON.parse(body),
                        body: body,
                        headers: response.headers,
                        status: response.statusCode,
                    };
                },
            };
            const res = await this.request(reqoptions)
            
            // console.log(res)
            this.csrf = res.headers["x-csrf-token"] ?
                res.headers["x-csrf-token"] :
                this.csrf;
            switch (res.status) {
                case 200:
                    const name = res.data.name;
                    const image = res.data.images[0].variations[2].url;
                    const foundBySize =
                        this.size.toUpperCase() == "RANDOM" ?
                            res.data.sellableUnits[
                            Math.floor(Math.random() * res.data.sellableUnits.length)
                            ] :
                            res.data.sellableUnits.find((i) =>
                                i.attributes.find((a) => a.value == this.size.split(",")[Math.floor(Math.random() * this.size.split(",").length)])
                            );
                    if (foundBySize) {
                        if (foundBySize.stockLevelStatus == "inStock") {
                            this.productInfo = {
                                name,
                                image,
                                price: foundBySize.price.value,
                                variant: foundBySize.code,
                            };
                        } else {
                            console.log("Item out of Stock");
                            await this.delay(this.monitorDelay);
                            await this.getItemDetails();
                        }
                    } else {
                        console.log("Unavailable size");
                    }
                    break;
                case 503:
                    console.log("in queue");
                    await this.delay(15000);
                    await this.getItemDetails();
                    break;
                default:
                    throw new Error("Unknow Error");
            }
        } catch (error) {
            if (error.statusCode == 400) {
                console.log("Product unavailable");
                await this.delay(this.restockDelay);
                await this.getItemDetails();
            } else if (error.response?.status == 429){
                console.log('Too MAny Requests')
                await this.delay(30000)
                await this.getItemDetails();
            } else if (error.message.includes('tunneling socket could not be established')){
                console.log('Proxy Error')
                console.log(error.message)
                await this.delay(50000)
                await this.getItemDetails();
            } else if (error.message.includes('Unexpected token < in')){
                console.log('Too MAny Requests')
                await this.delay(30000)
                await this.getItemDetails();
            } else {
                console.log(error.message)
                console.log("unable to Initiallized sessions");
                console.log("Proxy Error Retrying");
                await this.setProxy()
            //    console.log(error)
                await this.delay(this.restockDelay);
                await this.getItemDetails();
            }
        }
    }

    async addTocart() {
        if (this.stop) return;
        if (this.atcRetry == 1) {
            await this.setProxy()
        }
        try {
            const options = {
                method: "POST",
                body: `{"productQuantity":1,"productId":"${this.productInfo.variant}"}`,
                uri: `https://www.${this.site
                    }.com/api/users/carts/current/entries?timestamp=${Date.now()}`,
                headers: {
                    accept: "application/json",
                    "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    "cache-control": "no-cache",
                    "content-type": "application/json",
                    origin: "https://www.footlocker.com",
                    pragma: "no-cache",
                    referer: "https://www.footlocker.com/product/nike-air-force-1-low-mens/24300657.html",
                    "sec-fetch-dest": "empty",
                    "sec-fetch-mode": "cors",
                    "sec-fetch-site": "same-origin",
                    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
                    "x-csrf-token": this.csrf,
                    "x-fl-productid": this.productInfo.variant,
                    "x-fl-request-id": v4(),
                },
                jar: this.jar,
                proxy: this.proxy,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: JSON.parse(body),
                        headers: response.headers,
                        status: response.statusCode,
                    };
                },
            };
         //   console.log(this.jar)
            const res = await this.request(options);
            // console.log(res)
         //   console.log(this.jar)
            switch (res.status) {
                case 531:
                    console.log("Out of Stock");
                    throw new Error("Out of Stock");
                case 200:
                    if (res.data.guid) {
                        this.csrf = res.headers["x-csrf-token"];
                        this.cartId = res.data.guid;
                        return;
                    } else {
                        throw new Error("Unable to cart");
                    }
                default:
                    throw new Error("Unable to cart");
            }
        } catch (e) {
            if (e.statusCode == 403) {
                this.atcRetry = this.atcRetry - 1;
                if (this.atcRetry <= 0) {
                  //  console.log(e)
                    throw new Error("Error carting")
                }
                console.log("Getting Datadome cookie");
                console.log(
                    `Carting product with code: ${this.productInfo.variant}`
                );
                const foundCookie = e.response.headers["set-cookie"]
                    .find((i) => i.split("=")[0] == "datadome")?.split(";")[0]
                    .split("datadome=")[1];
                if (foundCookie) {
                    await this.setProxy()
                    console.log(`Getting Session`);
                    await this.initializeSession();
                    await this.delay(3000);
                    await this.addTocart();
                }
            } else {
                console.log("Retrying to cart");
                if (this.size.toUpperCase() == "RANDOM" || this.size.split(",").length > 1) {
                    await this.delay(3000);
                    await this.getItemDetails()
                }
                await this.delay(this.restockDelay);
                await this.addTocart();
            }
        }
    }

    async submitEmail() {

        if (this.stop) return;
        if (this.emailRetry == 1) {
            await this.setProxy()
        }
        try {
            const options = {
                method: "PUT",
                uri: `https://www.${this.site}.com/api/users/carts/current/email/${this.profile.email
                    }?timestamp=${Date.now()}`,
                headers: {
                    // Connection: "keep-alive", ////
                    accept: "application/json",
                    "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    "cache-control": "no-cache",
                    "content-type": "application/json",
                    origin: `https://www.${this.site}.com`,
                    referer: `https://www.${this.site}.com/checkout`,
                    pragma: "no-cache",
                    "sec-fetch-dest": "empty",
                    "sec-fetch-mode": "cors",
                    "sec-fetch-site": "same-origin",
                    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
                    "x-csrf-token": this.csrf,
                    "x-fl-request-id": v4(),
                },
                jar: this.jar,
                proxy: this.proxy,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: body,
                        status: response.statusCode,
                        headers: response.headers,
                    };
                },
            };
            const res = await this.request(options);
            if (res.status != 200) {
                console.log("unable to send email");
            } else {
              //  console.log(res);
            }
        } catch (e) {
         //   console.log(e)
            this.emailRetry = this.emailRetry - 1;
            if (this.emailRetry <= 0) {
             //   console.log(e)
                throw new Error("Error submitting Email" + e)
            }
            if (e.statusCode == 550) {
                await this.setProxy()
            }
            if (e.statusCode == 403) {
                await this.delay(3000);
              //  console.log(e)
                console.log("Getting Datadome cookie");
                const foundCookie = e.response.headers["set-cookie"]?.find((i) => i.split("=")[0] == "datadome")?.split(";")[0]
                    .split("datadome=")[1];
                if (foundCookie) {
                    await this.setProxy()
                    console.log(`Getting Session`);
                    await this.initializeSession();
                    await this.delay(3000);
                    await this.submitEmail();
                }
            } else {
                console.log(`Getting Session`);
                await this.initializeSession();
                await this.delay(3000);
                await this.submitEmail()
            }
        }
    }

    async submitShipping() {

        if (this.stop) return;
        if (this.shippingRetry == 1) {
            await this.setProxy()
        }
        try {
            const options = {
                method: "POST",
                uri: `https://www.${this.site
                    }.com/api/users/carts/current/addresses/shipping?timestamp=${Date.now()}`,
                headers: {
                    accept: "application/json",
                    "Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    "Cache-Control": "no-cache",
                    Connection: "keep-alive",
                    "content-type": "application/json",
                    Host: `www.${this.site}.com`,
                    Origin: `https://www.${this.site}.com`,
                    Pragma: "no-cache",
                    Referer: `https://www.${this.site}.com/product/-/${this.pid}.html`,
                    "Sec-Fetch-Dest": "empty",
                    "Sec-Fetch-Mode": "cors",
                    "Sec-Fetch-Site": "same-origin",
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
                    "x-csrf-token": this.csrf,
                    "x-fl-request-id": v4(),
                },
                jar: this.jar,
                proxy: this.proxy,
                body: `{"shippingAddress":{"setAsDefaultBilling":false,"setAsDefaultShipping":false,"firstName":"${this.profile.shippingAddress.firstName}","lastName":"${this.profile.shippingAddress.lastName}","email":"${this.email}","phone":"${this.profile.shippingAddress.phone}","billingAddress":false,"country":{"isocode":"US","name":"United States"},"defaultAddress":false,"id":null,"region":{"countryIso":"US","isocode":"${this.profile.shippingAddress.region.isocode}","isocodeShort":"${this.profile.shippingAddress.region.isocodeShort}","name":"${this.profile.shippingAddress.region.name}"},"setAsBilling":true,"shippingAddress":true,"visibleInAddressBook":false,"type":"default","LoqateSearch":"","postalCode":"${this.profile.shippingAddress.postalCode}","town":"${this.profile.shippingAddress.town}","regionFPO":null,"line1":"${this.profile.shippingAddress.line1}","recordType":" "}}`,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: body,
                        status: response.statusCode,
                        headers: response.headers,
                    };
                },
            };
            const res = await this.request(options);
            // console.log(res)
            if (res.status != 201) {
                console.log("unable to send shipping");
            //    console.log(res)
                // throw new Error(res)
            } else {
              //  console.log(res);
            }
        } catch (e) {
            // console.log(e.response)
            this.shippingRetry = this.shippingRetry - 1;
            if (this.shippingRetry <= 0) {
               // console.log(e)
                throw new Error("Error submitting shipping:" + e)
            }
            if (e.statusCode == 403) {
                await this.delay(3000);
                console.log("Getting Datadome cookie");
                const foundCookie = e.response.headers["set-cookie"]
                    .find((i) => i.split("=")[0] == "datadome")?.split(";")[0]
                    .split("datadome=")[1];
                if (foundCookie) {
                    await this.setProxy()
                    console.log(`Getting Session`);
                    await this.initializeSession();
                    await this.delay(3000);
                    await this.submitShipping();
                }
            } else {
                console.log(`Getting Session`);
                await this.initializeSession();
                await this.delay(3000);
                await this.submitShipping()
            }
        }
    }

    async submitBilling() {

        if (this.stop) return;
        if (this.billingRetry == 1) {
            await this.setProxy()
        }
        try {
            const options = {
                method: "POST",
                uri: `https://www.${this.site
                    }.com/api/users/carts/current/set-billing?timestamp=${Date.now()}`,
                headers: {
                    Connection: "keep-alive", ////
                    accept: "application/json",
                    "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                    "accept-encoding": "gzip, deflate, br",
                    "content-type": "application/json",
                    origin: `https://www.${this.site}.com`,
                    referer: `https://www.${this.site}.com/checkout`,
                    pragma: "no-cache",
                    "sec-fetch-dest": "empty",
                    "sec-fetch-mode": "cors",
                    "sec-fetch-site": "same-origin",
                    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
                    "x-csrf-token": this.csrf,
                    "x-fl-request-id": v4(),
                },
                jar: this.jar,
                proxy: this.proxy,
                body: `{"setAsDefaultBilling":false,"setAsDefaultShipping":false,"firstName":"${this.profile.billingAddress.firstName}","lastName":"${this.profile.billingAddress.lastName}","email":"${this.profile.billingAddress.email}","phone":"${this.profile.billingAddress.phone}","billingAddress":false,"country":{"isocode":"US","name":"United States"},"defaultAddress":false,"id":null,"region":{"countryIso":"US","isocode":"${this.profile.billingAddress.region.isocode}","isocodeShort":"${this.profile.billingAddress.region.isocodeShort}","name":"${this.profile.billingAddress.region.name}"},"setAsBilling":false,"shippingAddress":true,"visibleInAddressBook":false,"type":"default","LoqateSearch":"","postalCode":"${this.profile.billingAddress.postalCode}","town":"${this.profile.billingAddress.town}","regionFPO":null,"line1":"${this.profile.billingAddress.line1}","recordType":" "}`,
                transform: function (body, response, resolveWithFullResponse) {
                    return {
                        data: body,
                        status: response.statusCode,
                        headers: response.headers,
                    };
                },
            };
            const res = await this.request(options);
            if (res.status === 400) {
                console.log("Failed while submitting billing");
                await this.delay(5000);
                await this.submitBilling();
            }
            if (res.status !== 200) {
                console.log(" Billing error");
            }
            // console.log(res);
        } catch (e) {
            this.billingRetry = this.billingRetry - 1;
            if (this.billingRetry <= 0) {
              //  console.log(e)
                throw new Error("Error submitting Billing:" + e)
            }
            if (e.statusCode == 403) {
                await this.delay(3000);
             //   console.log(e)
                console.log("Getting Datadome cookie");
                const foundCookie = e.response.headers["set-cookie"]
                    .find((i) => i.split("=")[0] == "datadome")?.split(";")[0]
                    .split("datadome=")[1];
                if (foundCookie) {
                    await this.setProxy()
                    console.log(`Getting Session`);
                    await this.initializeSession();
                    await this.delay(3000);
                    await this.submitBilling();
                }
            } else {
                console.log(`Getting Session`);
                await this.initializeSession();
                await this.delay(3000);
                await this.submitBilling()
            }
        }
    }

    EncryptCard(key, card) {
        const cseInstance = adyenEncrypt.createEncryption(key, {});
        return {
            encryptedCardNumber: cseInstance.encrypt({
                number: card.number,
                generationtime: new Date().toISOString(),
            }),
            encryptedExpiryMonth: cseInstance.encrypt({
                expiryMonth: card.month,
                generationtime: new Date().toISOString(),
            }),
            encryptedExpiryYear: cseInstance.encrypt({
                expiryYear: card.year,
                generationtime: new Date().toISOString(),
            }),
            encryptedSecurityCode: cseInstance.encrypt({
                cvc: card.ccv,
                generationtime: new Date().toISOString(),
            }),
        };
    }
    async submitOrder() {
        if (this.stop) return;
        const aydenKey =
            "10001|A237060180D24CDEF3E4E27D828BDB6A13E12C6959820770D7F2C1671DD0AEF4729670C20C6C5967C664D18955058B69549FBE8BF3609EF64832D7C033008A818700A9B0458641C5824F5FCBB9FF83D5A83EBDF079E73B81ACA9CA52FDBCAD7CD9D6A337A4511759FA21E34CD166B9BABD512DB7B2293C0FE48B97CAB3DE8F6F1A8E49C08D23A98E986B8A995A8F382220F06338622631435736FA064AEAC5BD223BAF42AF2B66F1FEA34EF3C297F09C10B364B994EA287A5602ACF153D0B4B09A604B987397684D19DBC5E6FE7E4FFE72390D28D6E21CA3391FA3CAADAD80A729FEF4823F6BE9711D4D51BF4DFCB6A3607686B34ACCE18329D415350FD0654D";
        const cardInfo = this.EncryptCard(aydenKey, {
            number: this.profile.paymentInfo.cardNumber,
            month: this.profile.paymentInfo.cardMonth,
            year: this.profile.paymentInfo.cardYear,
            ccv: this.profile.paymentInfo.cardCvv,
        });
    //    console.log(cardInfo)
        const payload = {
            preferredLanguage: "en",
            termsAndCondition: false,
            deviceId: "",
            cartId: this.cartId,
            encryptedCardNumber: cardInfo.encryptedCardNumber,
            encryptedExpiryMonth: cardInfo.encryptedExpiryMonth,
            encryptedExpiryYear: cardInfo.encryptedExpiryYear,
            encryptedSecurityCode: cardInfo.encryptedSecurityCode,
            paymentMethod: "CREDITCARD",
            returnUrl: `https://www.${this.site}.com/adyen/checkout`,
            browserInfo: {
                screenWidth: 1920,
                screenHeight: 1080,
                colorDepth: 24,
                userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
                timeZoneOffset: -120,
                language: "en-US",
                javaEnabled: false,
            },
        };
        const options = {
            method: "POST",
            uri: `https://www.${this.site
                }.com/api/users/orders?timestamp=${Date.now()}`,
            body: JSON.stringify(payload),
            headers: {
                accept: "application/json",
                "accept-encoding": "gzip, deflate, br",
                "accept-language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7,de;q=0.6",
                "content-type": "application/json",
                origin: `https://www.${this.site}.com`,
                referer: `https://www.${this.site}.com/checkout`,
                "sec-fetch-dest": "empty",
                "sec-fetch-mode": "cors",
                "sec-fetch-site": "same-origin",
                "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
                "x-csrf-token": this.csrf,
                "x-fl-request-id": v4(),
            },
            jar: this.jar,
            proxy: this.proxy,
            transform: function (body, response, resolveWithFullResponse) {
                return {
                    data: JSON.parse(body),
                    status: response.statusCode,
                };
            },
        };
        try {
            const res = await this.request(options);
           // console.log(res);
            if (res.status == 200) {
                console.log(chalk.green("Success"));
            } else {
                console.log(chalk.magenta("Status Unknown"));
            }
            console.log("Success");
            await sendNotification({
                store: this.site,
                image: this.productInfo.image,
                itemName: this.productInfo.name,
                size: this.task.size,
                price: this.productInfo.price,
                pid: this.task.pid,
                status: "Success",
                storeName: this.task.site,
                module: this.task.mode,
                kw: this.task.keywords,
                id: this.task.id,
            })
            // await fs.writeFileSync('./checkout.json', JSON.stringify(res.data))
        } catch (error) {
        //    console.log(error)
        //    console.log(error?.response?.data)
            console.log(chalk.red("Failed"));
            await sendNotification({
                store: this.site,
                image: this.productInfo.image,
                itemName: this.productInfo.name,
                size: this.task.size,
                price: this.productInfo.price,
                pid: this.task.pid,
                status: "Failed",
                storeName: this.task.site,
                module: this.task.mode,
                kw: this.task.keywords,
                id: this.task.id,
            })
            // if (error.statusCode == 400) {
            //   let dreason = JSON?.parse(
            //     error?.response?.data
            //   )?.errors[0]?.message?.includes("not available") ?
            //     "out of stock" :
            //     "payment failed";
            //   console.log(error.response);
            //   console.log(dreason);
            //   return;
            // }
            // console.log("Failed");
            // console.log(error.response);
        }
    }

    delay(time) {
        return new Promise(function (resolve) {
            setTimeout(resolve, time);
        });
    }

    async allCookies(jar) {
        const store = jar._jar.store;
        return (
            await Promise.all(
                Object.keys(store.idx).map((d) =>
                    util.promisify(store.findCookies).call(store, d, null)
                )
            )
        ).flat();
    }


}


module.exports = {
    Footsites
}
