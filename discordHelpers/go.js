const { rp,firstServer, port } = require('./helper')

async function startGoMonitor(currentBody, site) {
    try {
      switch (site) {
        case 'BESTBUY': {
          rp.post(
            {
              url: `${firstServer}:${port}/${site}`,
              body: JSON.stringify(currentBody),
              headers: {
                'Content-Type': 'application/json',
              },
            },
            (response) => console.log(response?.statusCode)
          );
          break;
        }
        default: {
          rp.post(
            {
              url: `${firstServer}:${port}/${site}`,
              body: JSON.stringify(currentBody),
              headers: {
                'Content-Type': 'application/json',
              },
            },
            (response) => console.log(response?.statusCode)
          );
        }
      }
    } catch (error) {
      console.log(`Error Starting Go Monitor ${error}`);
    }
  }

  module.exports = {
      startGoMonitor
  }