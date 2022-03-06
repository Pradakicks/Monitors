async function getSku(skuName, proxyList) {
    try {
      let proxy1 = proxyList[Math.floor(Math.random() * proxyList.length)];
      console.log(proxy1);
      let fetchProductPage = await rp.get({
        url: `https://www.newegg.com/prada/p/${skuName}`,
        headers: {
          accept:
            'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
          'accept-language': 'en-US,en;q=0.9',
          'cache-control': 'no-cache',
          pragma: 'no-cache',
          'sec-fetch-dest': 'document',
          'sec-fetch-mode': 'navigate',
          'sec-fetch-site': 'none',
          'sec-fetch-user': '?1',
          'upgrade-insecure-requests': '1',
        },
        proxy: `http://${proxy1.userAuth}:${proxy1.userPass}@${proxy1.ip}:${proxy1.port}`,
      });
      console.log(fetchProductPage?.statusCode);
      sku =
        fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[0] + '-';
      console.log(sku);
      sku =
        sku +
        fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[1] +
        '-';
      console.log(sku);
      sku =
        sku + fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[2];
      console.log(sku);
      return sku;
    } catch (error) {
      console.log(error);
      await delay(10000);
      await getSku(skuName, proxyList);
    }
  }