const {rp } = require('../discordHelpers/helper')


async function walmartScraper(clients, triggerText, replyText) {
    try {
      clients.on('message', async (message) => {
        if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
          const { group, isValidated } = await checkIfUserValidated(message);
          let currentCompany = config.groups[group];
          if (isValidated) {
            console.log('Walmart Scraper Initiated');
            const SKU = message.content.split('/')[5];
            let currentMessage = await message.channel.send(
              `${message.author} Walmart Scraper for ${SKU} Initiated`
            );
            let data = await getTerraSku(SKU);
            if (data == 'Error') {
              currentMessage.edit(
                `${message.author} Error Grabbing Data for ${SKU}`
              );
              return;
            }
            try {
              let ImageOffer = data.payload.selected.defaultImage;
              let selectedProduct = data.payload.selected.product;
              let productName =
                data.payload.products[selectedProduct].productAttributes
                  .productName;
              let ImageUrl =
                data.payload.images[ImageOffer].assetSizeUrls.DEFAULT;
              let walmartSellerOffer;
  
              Object.keys(data.payload.sellers).map((e) => {
                let currentSeller = data?.payload?.sellers[e];
                if (currentSeller?.sellerName == 'Walmart.com')
                  walmartSellerOffer = currentSeller?.sellerId;
              });
  
              console.log(ImageUrl, walmartSellerOffer);
              let offerObj;
              if (!walmartSellerOffer) {
                currentMessage.edit(
                  `${message.author} ${SKU} is not sold by Walmart`
                );
                offerObj = {
                  offerId: 'N/A',
                  availabilityStatus: 'Not Sold By Walmart',
                  price: 'N/A',
                  image: ImageUrl,
                  sku: SKU,
                  productName: productName,
                };
              } else {
                Object.keys(data.payload.offers).forEach((e) => {
                  let currentOffer = data.payload.offers[e];
                  console.log(currentOffer?.pricesInfo);
                  if (currentOffer.sellerId == walmartSellerOffer)
                    offerObj = {
                      offerId: currentOffer.id ? currentOffer.id : 'N/A',
                      availabilityStatus: currentOffer?.productAvailability
                        ?.availabilityStatus
                        ? currentOffer.productAvailability.availabilityStatus
                        : 'N/A',
                      price: currentOffer?.pricesInfo?.priceMap?.CURRENT?.price
                        ? currentOffer.pricesInfo.priceMap.CURRENT.price
                        : 'No Price Found',
                      image: ImageUrl,
                      sku: SKU,
                      productName: productName,
                    };
                });
              }
  
              const currentEmbed = new MessageEmbed()
                .setColor(currentCompany.companyColorV2)
                .setTitle(offerObj.productName)
                .setURL(`https://www.walmart.com/ip/prada/${offerObj.sku}`)
                .setThumbnail(offerObj.image)
                .addFields(
                  {
                    name: 'Offer Id',
                    value: offerObj.offerId,
                  },
                  {
                    name: 'PID',
                    value: offerObj.sku,
                  },
                  {
                    name: 'Availability',
                    value: offerObj.availabilityStatus,
                    inline: true,
                  },
                  {
                    name: "Walmart's Price",
                    value: offerObj.price,
                    inline: true,
                  },
                  {
                    name: 'Links',
                    value: `[Product](https://www.walmart.com/ip/prada/${offerObj.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${offerObj.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`,
                  }
                )
                .setFooter('Prada#4873')
                .setTimestamp();
              message.channel.send(currentEmbed);
              currentMessage.delete();
            } catch (error) {
              console.log(error);
              currentMessage.edit(
                `${message.author} Error Grabbing Data for ${SKU} BOJ`
              );
            }
          } else {
            message.channel.send(`${message.author} is not a validated user`);
          }
        }
      });
    } catch (error) {
      console.log(error);
    }
  }

async function getTerraSku(SKU) {
    try {
      // let proxies = await getProxies()
      const proxy1 = proxyList[Math.floor(Math.random() * proxyList.length)];
      console.log(proxy1);
      let terra = await rp.get({
        url: `https://www.walmart.com/terra-firma/item/${SKU}`,
        headers: {
          accept:
            'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
          'accept-language': 'en-US,en;q=0.9',
          'cache-control': 'no-cache',
          pragma: 'no-cache',
          'sec-ch-ua':
            '" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"',
          'sec-ch-ua-mobile': '?0',
          'sec-fetch-dest': 'document',
          'sec-fetch-mode': 'navigate',
          'sec-fetch-site': 'none',
          'sec-fetch-user': '?1',
          'service-worker-navigation-preload': 'true',
          'upgrade-insecure-requests': '1',
          'user-agent':
            'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36',
        },
        proxy: `http://${proxy1.userAuth}:${proxy1.userPass}@${proxy1.ip}:${proxy1.port}`,
      });
      console.log(`Terra Status Code : ${terra.statusCode}`);
      if (terra.statusCode != 200 && terra.statusCode != 201) {
        return 'Error';
      } else {
        let parsed = JSON.parse(terra?.body);
        return parsed;
      }
    } catch (error) {
      console.log(error);
      return 'Error';
    }
  }

module.exports = {
    walmartScraper
}