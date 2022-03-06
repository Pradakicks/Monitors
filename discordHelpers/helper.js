
const rp = require('request-promise').defaults({
    followAllRedirects: true,
    resolveWithFullResponse: true,
    gzip: true,
  });

const firstServer = `http://localhost`; // 12 Core z 24GB
const port = 7243;

module.exports = {
    rp,
    firstServer,
    port
}