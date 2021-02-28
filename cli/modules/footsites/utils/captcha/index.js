
const recaptcha = require("./recaptcha")

const {autoSolveBank, requestTokenFromAYCD} = require('./aycd');
const Configstore = require('configstore');
const config = new Configstore('VibrisCLI')


let _2captcha = false;

// const getCaptchaToken = async (id) => {
//     return new Promise((res, rej) => {
//         let myIntervall = setInterval(async () => {
//             console.log(`waiting for captcha: ${id}`)
//             try {
//                 console.log(manual_solve.captchaBank)
//                 let foundIndex = manual_solve.captchaBank.findIndex((i) => i.id = id)
//                 if (foundIndex>-1) {
//                     console.log("index",foundIndex)
//                     let found = manual_solve.captchaBank[foundIndex]
//                     manual_solve.captchaBank.splice(foundIndex, 1)
//                     clearInterval(myIntervall);
//                     console.log("found",found)
//                     res(found.token)
//                 }
//             } catch (error) {
//                 console.log(error)
//             }
//         }, 200);
//     })
// }

module.exports = async ({ initialURL, sitekey, id }) => {
    let cm = config.get('captcha')
    let captchaMethod = cm.split('/')[0]
    let captchaKey = cm.split('/')[1]
    
    if (captchaMethod == 'AYCD') {
        try {
            await requestTokenFromAYCD(initialURL, sitekey, id)
            return new Promise(async (resolve, reject) => {
                let int = setInterval(() => {
                    let foundIndex = autoSolveBank.findIndex(item => {
                        let g = (JSON.parse(item))
                        //   console.log(g)
                        return g.taskId === id
                    });
                    if (foundIndex > -1) {
                        let parsed = JSON.parse(autoSolveBank[foundIndex])
                        if (parsed?.taskId == id) {
                            let t = parsed?.token
                            autoSolveBank.splice(foundIndex, 1)
                            clearInterval(int)
                            resolve(t)
                        }
                    }

                }, 200)
            })
            
        } catch (error) {
            console.log(error)
            return null
        }
    } else if (captchaMethod == '2Cap') {
        try {
            const res = await recaptcha({ initialURL, sitekey, captchakey: captchaKey })
            return res
        } catch (error) {
            return null
        }
    } 
}