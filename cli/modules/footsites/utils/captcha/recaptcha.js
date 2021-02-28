const Client = require('@infosimples/node_two_captcha');


const solveCaptcha = async ({ initialURL, sitekey, captchakey }) => {
    // 2captcha configuration
    let client = new Client(captchakey, {
        timeout: 100000,
        polling: 5000,
        throwErrors: false
    });
    let recaptcha_response;
    try {
        let response = await client.decodeRecaptchaV2({
            googlekey: sitekey,
            pageurl: initialURL
        })
        let recaptcha_response = response.text
        return recaptcha_response
    } catch (error) {
        return recaptcha_response
    }
}

module.exports = solveCaptcha;