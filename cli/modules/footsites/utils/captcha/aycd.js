const AutoSolve = require('autosolve-client');
const fs = require('fs')


//Then, we instantiate the connections
var autoSolveBank = []


let autoSolveInit = async (accessToken, apiKey) => {
    try {
        const autoSolve = AutoSolve.getInstance({
            "accessToken": accessToken,
            "apiKey": apiKey,
            "clientKey": "VibrisBots-96ceed26-20a7-4636-bdfa-10ffad5c5ec7",
            "shouldAlertOnCancel": true,
            "debug": true
        });
    
        await autoSolve.init(accessToken, apiKey).then(() => {
    
            //Register our handlers for responses from AutoSolve
            autoSolve.ee.on(`AutoSolveResponse`, (data) => {
                console.log(`Auto Solve Response ${data}`)
                autoSolveBank.push(data)
                //do stuff here with token response object for task
            })
    
            autoSolve.ee.on(`AutoSolveResponse_Cancel`, (data) => {
                console.log(`Auto Solve Response_Cancel ${data}`)
                //do stuff here to handle an AutoSolve request that was cancelled
            })
    
            autoSolve.ee.on(`AutoSolveError`, (data) => {
                console.log(`Auto Solve Error ${data}`)
                //emits errors related to connection events
            })
        })
        return true
    } catch (error) {
        fs.appendFileSync('errorCli.txt', error.toString() + '\n', (err => {
            console.log(err)
        }))
    }
   
    
};

// console.log(autoSolveInit('132358-0fa77464-ac49-4151-a986-91d305510414', 'd5a6543a-5ae3-49a3-bdb8-deb20defcdc2'))

let requestTokenFromAYCD = async (url, siteKey, id) => {
    try {
           // let autoSolveObject = AutoSolve.getInstance();
    
    const autoSolve = AutoSolve.getInstance();

    let requestToken = await autoSolve.sendTokenRequest({
        taskId: id,
        url: url,
        siteKey: siteKey,
        version: 0,
    })

    } catch (error) {
        fs.appendFileSync('errorCli.txt', error.toString() + '\n', (err => {
            console.log(err)
        }))
    }
 



};

module.exports = {
    autoSolveBank,
    autoSolveInit,
    requestTokenFromAYCD
}

    







// console.log(`autoSolveObject ${(JSON.stringify(autoSolveObject))} `)


