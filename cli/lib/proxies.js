const Configstore = require('configstore');
const { get } = require('request-promise');
const config = new Configstore('VibrisCLI')

const getProxy = async () => {
    let proxy = config.get('proxies')
  //  console.log(proxy)
    let group = []
    let type = typeof proxy[0]
    let typ1e = Array.isArray(proxy[0])
    // console.log(proxy[0])
    // console.log(typ1e)
    // console.log(type)
    //  console.log(proxy.split('\n').length)
    if(type == 'string'){
    //    console.log('test')
        await proxy.forEach(e => {
            group.push({ proxyUrl: e, inUse: false })
        })
     //   console.log(group)
    } else {
        group = proxy
    }
  //  console.log(group)
    // console.log(proxy)
    let groupList = group;
    let foundProxy = groupList.find((i) => i.inUse == false) ? groupList.find((i) => i.inUse == false) : groupList[Math.floor(Math.random() * groupList.length)]
 //   console.log(foundProxy)
    let elementsIndex = groupList.findIndex(element => element.proxyUrl == foundProxy.proxyUrl);
    let newArray = [...groupList]
    newArray[elementsIndex] = { ...newArray[elementsIndex], inUse: true }
   // console.log(newArray)
    await config.set('proxies', newArray)

    let goodProxy;
    var arraysplitted = foundProxy.proxyUrl.split(':');
    var IP = arraysplitted.shift();
    var PORT = arraysplitted.shift();
    var USERNAME = arraysplitted.shift();
    var PASSWORD = arraysplitted.shift();
    if (USERNAME) {
        goodProxy = `http://${USERNAME}:${PASSWORD}@${IP}:${PORT}`;
    } else {
        goodProxy = `http://${IP}:${PORT}`;
    }
 //   console.log(foundProxy.proxyUrl)
    return {
        formatted: goodProxy,
        unformatted: foundProxy.proxyUrl
    };
} 


const proxyIsNotInUse = async (proxy) => {
    if (!proxy) return;
    let proxys = config.get('proxies')
    let groupList = proxys;
    let elementsIndex = groupList.findIndex(element => element.proxyUrl == proxy);
    let newArray = [...groupList]
    newArray[elementsIndex] = { ...newArray[elementsIndex], inUse: false }
    await config.set('proxies', newArray)
}

module.exports = {
    getProxy,
    proxyIsNotInUse
}