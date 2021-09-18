const got = require('got').default;
const { promisify } = require('util');
// const jsdom = require('jsdom');
// const { JSDOM } = jsdom;
const { CookieJar } = require('tough-cookie');
const cheerio = require('cheerio');
const cookieJar = new CookieJar();
const setCookie = promisify(cookieJar.setCookie.bind(cookieJar));
const ProxyAgent = require('proxy-agent');

class amazon {
  constructor(task) {
    this.userSettings = task;
    this.taskInfo = task.taskInfo;
    this.local = task.localUrl;
    this.proxy = task?.proxy;
    // this.proxy = task.proxy?.split(":").length > 2 ? task.proxy.trim().split(":") : null;
    const httpsAgent = new ProxyAgent(this.proxy);
    const httpAgent = new ProxyAgent(this.proxy);
    this.client = got.extend({
      cookieJar: cookieJar,
      headers: {
        authority: 'www.amazon.com',
        'sec-ch-ua':
          '"Chromium";v="88", "Google Chrome";v="89", ";Not A Brand";v="99"',
        'sec-ch-ua-mobile': '?0',
        dnt: '1',
        'upgrade-insecure-requests': '1',
        'user-agent':
          'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36',
        accept:
          'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'sec-fetch-site': 'none',
        'sec-fetch-mode': 'navigate',
        'sec-fetch-user': '?1',
        'sec-fetch-dest': 'document',
        'accept-language': 'en-US,en;q=0.9',
        cookie: `session-id=${this.sid};`,
      },
      https: {
        rejectUnauthorized: false,
      },

      agent: {
        http: httpAgent,
        https: httpsAgent,
      },
    });
  }

  async init() {
    await this.run();
  }

  async run() {
    // var monitorStatus = false;
    await this.delay(500);
    // console.log(this.taskInfo, this)
    var asin = this.taskInfo.asin;
    const monitorStatus = async () => {
      console.log('Searching For Item');
      const status = true;
      if (status) {
        this.taskInfo.productUrl = `https://www.amazon.com/gp/product/${asin}`;
        console.log(
          `https://www.amazon.com/gp/aod/ajax/ref=auto_load_aod?asin=${this.taskInfo.asin}&pc=dp`
        );
        await this.getProduct();
        await this.firstMonitor();
        await this.sendCookies();

        // await this.MainMonitor();
      } else {
        setTimeout(function () {
          monitorStatus();
        }, this.taskInfo.monitorDelay);
        // monitorStatus()
      }
    };
    monitorStatus();
  }

  async getProduct() {
    try {
      const response = await this.client.get(this.taskInfo.productUrl);
      // * Parse the product response for amazon formdata
      console.log('Generating Session');
      let $ = cheerio.load(response.body);
      this.image = $('#landingImage').attr('src');
      this.name = $('#productTitle').text();
      this.smid = $('input[id="merchantID"]').attr('value');
      this.sid = $('input[id="session-id"]').attr('value');
      console.log(this.sid);
    } catch (error) {
      console.log(error);
      await this.delay(this.taskInfo.monitorDelay);
      await this.getProduct();
    }
  }

  async firstMonitor() {
    try {
      console.log('Monitoring Product');
      const response = await this.client.get(
        `https://www.amazon.com/gp/aod/ajax/ref=auto_load_aod?asin=${this.taskInfo.asin}&pc=dp`,
        {
          cookieJar: cookieJar,
        }
      );
      let $ = cheerio.load(response.body);
      this.csrftoke = $('input[id="aod-atc-csrf-token"]').attr('value');
      console.log(this.csrftoke);
    } catch (error) {
      console.log(error);
      await this.delay(this.taskInfo.monitorDelay);
      await this.firstMonitor();
    }
  }

  async MainMonitor() {
    try {
      console.log('Monitoring');
      const response = await this.client.post(
        'https://data.amazon.com/api/marketplaces/ATVPDKIKX0DER/cart/carts/retail/items?ref=aod_dpdsk_used_1',
        {
          headers: {
            Connection: 'keep-alive',
            Pragma: 'no-cache',
            'Cache-Control': 'no-cache',
            'sec-ch-ua':
              '" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"',
            DNT: '1',
            'x-api-csrf-token': this.csrftoke,
            'Accept-Language': 'en-US',
            'sec-ch-ua-mobile': '?0',
            'User-Agent':
              'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
            'Content-Type':
              'application/vnd.com.amazon.api+json; type="cart.add-items.request/v1"',
            Accept:
              'application/vnd.com.amazon.api+json; type="cart.add-items/v1"',
            Origin: 'https://www.amazon.com',
            'Sec-Fetch-Site': 'same-site',
            'Sec-Fetch-Mode': 'cors',
            'Sec-Fetch-Dest': 'empty',
            Referer: 'https://www.amazon.com/',
            cookie: `session-id=${this.sid};`,
          },
          body: `{"items":[{"asin":"${this.taskInfo.asin}","offerListingId":"${this.taskInfo.offerId}","quantity":1}]}`,
          cookieJar: cookieJar,
        }
      );
      // * Parse the product response for amazon formdata
      console.log(response.statusCode);
      console.log(response.body);
    } catch (error) {
      console.log('Out of Stock, Retrying');
      console.log(error.response?.statusCode);
      console.log(error.response?.headers);
      await this.delay(this.taskInfo.monitorDelay);
      await this.MainMonitor();
    }
  }

  async sendCookies() {
    if (this.csrftoke && this.sid) {
        console.log({
            url: `${this.local}/AMAZONSESSION`,
            json: {
              sid: this.sid,
              csrf: this.csrftoke,
            },
          })
      let res = await got.post({
        url: `${this.local}/AMAZONSESSION`,
        json: {
          sid: this.sid,
          csrf: this.csrftoke,
        },
      });
      console.log(res.body)
    } else {
        await this.delay(5000)
        await this.getProduct();
        await this.firstMonitor();
        await this.sendCookies(); 
    }
  }
  delay(time) {
    return new Promise(function (resolve) {
      setTimeout(resolve, time);
    });
  }
}

const task = {
  password: '',
  proxy: null,
  userAgent:
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36',
  taskInfo: {
    asin: 'B09979GS5W',
    offerId:
      'f9FL%2FHvb8mR%2Fj3J361zTxy2J87piw0thjhoprF8qOayYYFTZsg1cIfDAEcG5D8tXKxNQiQsTBT%2F53D4%2FhW1Kgt1ruMYpP6XrgF5pie3MZfHXZ0xLMVgGaTUrNS9DOrabFGjD3TIzGMTMlDoz%2FKmCzA%3D%3D',
    // productUrl: 'https://www.amazon.com/gp/product/B09979GS5W',
    // priceRange: 1729.99,
    // monitorDelay: 500,
    // dom: ''
  },
};
// const d = new amazon(task);

// for (let index = 0; index < 1; index++) {
//   d.init();
// }
module.exports = amazon;
