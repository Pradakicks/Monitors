const request = require("request-promise");


const getMyIp = async (proxy) => {
    return new Promise((res, rej) => {
        request('https://api.ipify.org/?format=json', {
            proxy:proxy
        })
            .then(function (htmlString) {
                const jsoNParse = JSON.parse(htmlString)
                console.log(jsoNParse.ip)
                res(jsoNParse.ip)
            })
            .catch(function (err) {
                console.log(err)
                rej(err)
            });
    })
}

const randomRange = (min, max) => ~~(Math.random() * (max - min + 1)) + min;

function parseUrlBody(url) {
    let dd;
    try {
        dd = url;
        dd = {
            cid: dd.split('Cid=')[1].split('&')[0],
            hsh: dd.split('hash=')[1].split('&')[0],
            t: dd.split('t=')[1].split('&')[0],
            s: dd.split('s=')[1],
        };
    } catch (ex) {
        console.log(ex)
    }
    return dd;
}

function createQuery(
    {domain,
        cid,
        icid,
        hash,
        ip,
        s,
        captchaResponse }
) {
    let queryString = `?cid=${encodeURIComponent(cid)}`;
    queryString += `&icid=${encodeURIComponent(icid)}`;
    queryString += '&ccid=' + 'null';
    if (typeof captchaResponse === 'string') {
        queryString += `&g-recaptcha-response=${captchaResponse}`;
    } else {
        for (const i in captchaResponse) {
            queryString += `&${i.replace('_', '-response-')}=${encodeURIComponent(captchaResponse[i])}`;
        }
    }
    queryString += `&hash=${hash}`;
    queryString += `&ua=${encodeURIComponent(
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36',
    )}`;
    queryString += `&referer=${encodeURIComponent(domain)}`;
    queryString += `&parent_url=${encodeURIComponent(domain)}`;
    queryString += `&x-forwarded-for=${ip}`;
    queryString += `&captchaChallenge=${randomRange(100000000, 199999999)}`;
    queryString += `&s=${s}`;
    return queryString;
}


const getDataDomeCookie = async ({initialURL,recaptcha_response = "", proxy = null }) => {
    try {
        let ip = await getMyIp(proxy);
        let dd = await parseUrlBody(initialURL);
        if (!dd)return null;
        let checkURL = initialURL.replace("/captcha/", "/captcha/check").replace("initialCid", "icid")
            + "&s=17434"
            + "&captchaChallenge=180410550"
            + "&ccid=null"
            + '&g-recaptcha-response=' + encodeURIComponent(recaptcha_response)
            + '&x-forwarded-for' + encodeURIComponent(ip);
        let res = await request.get(checkURL, {
            headers: {
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
                "Accept": "*/*",
                "Accept-Language": "en-US,en;q=0.5",
                "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"
            },
            proxy: proxy ? proxy : null,
            transform: function (body, response, resolveWithFullResponse) {
                return {
                    data: JSON.parse(body),
                    status: response.statusCode
                };
            },
        });
        if (res.status == 200) {
            let datadome = res.data.cookie.split(";")[0];  // "datadome=6vIjy19aV
            console.log(" Success  [getDataDomeCookie] | code :", datadome);
            return datadome
        }

    } catch (error) {
        console.log(error)
        return null;
    }
}

module.exports = getDataDomeCookie;