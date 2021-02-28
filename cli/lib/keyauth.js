const { machineId, machineIdSync } = require('node-machine-id');
const { LICENSE_API_HOST, BOT_NAME, BOT_VERSION, AUTHORIZATION } = require('./constants');
const axios = require("axios")
const client = require('discord-rich-presence')('810325118550409286');


setTimeout(function () {
    client.updatePresence({
      state: 'Redefining Automation',
      details: '1.0.3 CLI', // Version
      startTimestamp: Date.now(),
      // endTimestamp: Date.now() + 1337,
      largeImageKey: 'rich2',
      // smallImageKey: 'rich2',
      instance: true,
    });
  }, 1000);


consoleIt = function(t) {
    process.stdout.write(t + '\n')
}

async function updateLicense(key, hwid) {
    await axios.patch(`https://api.metalabs.io/v4/licenses/${key}`, {
        metadata: { hwid }
    }, {
        headers: {
            'Authorization': `${AUTHORIZATION}`
        },
    }).then(res => res.data).then(async data => {
    //    consoleIt(data)
    })
}
async function retrieveLicenseCLI(keyInput = null) {
    try {
        let key;
        let id = await machineId();
        if (keyInput) {
            key = keyInput;
        }
        const license = await axios(`https://api.metalabs.io/v4/licenses/${key}`, {
            headers: {
                'Authorization': `${AUTHORIZATION}`
            }
        });
        const res = await license.data
     //   consoleIt(res)
            if (res.metadata.hwid) {
     
                if (res.metadata.hwid == id) {
                    consoleIt('Instance on Machine')
                  return true
                }
          
                consoleIt('License is already in use on another machine');
                return false
        } else {
            consoleIt('License is good to go!');
            await updateLicense(key, id);
            return true
        }

        // return res
    } catch (e) {
       // consoleIt(e)
        return false;
    }
}

module.exports = {
    updateLicense,
    retrieveLicenseCLI
}